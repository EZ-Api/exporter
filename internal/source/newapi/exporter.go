// Package newapi provides the export logic for New API data.
package newapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EZ-Api/exporter/internal/schema"
)

// ExporterConfig holds configuration for the exporter.
type ExporterConfig struct {
	IncludeTokens    bool // Whether to include tokens in export
	IncludeAbilities bool // Whether to include abilities (bindings)
	Verbose          bool // Enable verbose logging
}

// DefaultExporterConfig returns default configuration.
func DefaultExporterConfig() ExporterConfig {
	return ExporterConfig{
		IncludeTokens:    true,
		IncludeAbilities: false,
		Verbose:          false,
	}
}

// Exporter handles the export process from New API.
type Exporter struct {
	connector *Connector
	config    ExporterConfig
	result    *schema.ExportResult
}

// NewExporter creates a new exporter instance.
func NewExporter(connector *Connector, config ExporterConfig) *Exporter {
	return &Exporter{
		connector: connector,
		config:    config,
		result:    schema.NewExportResult(),
	}
}

// Export performs the full export process.
func (e *Exporter) Export() (*schema.ExportResult, error) {
	// Export channels -> providers
	if err := e.exportChannels(); err != nil {
		return nil, fmt.Errorf("failed to export channels: %w", err)
	}

	// Export users and tokens -> masters and keys
	if e.config.IncludeTokens {
		if err := e.exportUsersAndTokens(); err != nil {
			return nil, fmt.Errorf("failed to export users/tokens: %w", err)
		}
	}

	// Export abilities -> bindings (optional)
	if e.config.IncludeAbilities {
		if err := e.exportAbilities(); err != nil {
			return nil, fmt.Errorf("failed to export abilities: %w", err)
		}
	}

	return e.result, nil
}

// exportChannels exports all channels as providers.
func (e *Exporter) exportChannels() error {
	channels, err := e.connector.GetAllChannels()
	if err != nil {
		return err
	}

	for _, ch := range channels {
		providers := e.channelToProviders(ch)
		for _, p := range providers {
			e.result.AddProvider(p)
		}
	}

	return nil
}

// channelToProviders converts a New API channel to one or more EZ-API providers.
// Multi-key channels are split into multiple providers.
func (e *Exporter) channelToProviders(ch Channel) []schema.Provider {
	// Parse keys (newline separated for multi-key)
	keys := parseKeys(ch.Key)
	isMultiKey := len(keys) > 1

	// Parse groups (comma separated for multi-group)
	groups := parseGroups(ch.Group)
	primaryGroup := groups[0]

	// Map channel type to provider type
	providerType, typeOK := MapChannelType(ch.Type)
	if !typeOK {
		e.result.AddWarning(fmt.Sprintf(
			"Channel '%s' (ID=%d) has unknown type %d, mapped to 'custom'",
			ch.Name, ch.ID, ch.Type,
		))
	}

	// Map status
	status := MapChannelStatus(ch.Status)

	// Base URL
	baseURL := ""
	if ch.BaseURL != nil {
		baseURL = *ch.BaseURL
	}

	// Weight
	weight := 1
	priority := 0
	if ch.Weight != nil && *ch.Weight > 0 {
		weight = int(*ch.Weight)
	} else if ch.Priority != nil && *ch.Priority > 0 {
		weight = int(*ch.Priority)
	}
	if ch.Priority != nil && *ch.Priority > 0 {
		priority = int(*ch.Priority)
	}

	// Auto ban
	autoBan := true
	if ch.AutoBan != nil {
		autoBan = *ch.AutoBan == 1
	}

	// Parse models
	models := parseModels(ch.Models)

	// Create original backup
	original := e.createOriginalBackup(ch)

	// Check for unmappable fields and add warnings
	e.checkUnmappableFields(ch)

	// Create providers for each key
	var providers []schema.Provider
	for i, key := range keys {
		name := ch.Name
		if isMultiKey && i > 0 {
			name = fmt.Sprintf("%s-%d", ch.Name, i+1)
		}

		p := schema.Provider{
			OriginalID:   ch.ID,
			Name:         name,
			Type:         providerType,
			BaseURL:      baseURL,
			APIKey:       key,
			Models:       models,
			PrimaryGroup: primaryGroup,
			AllGroups:    groups,
			Weight:       weight,
			Priority:     priority,
			Status:       status,
			AutoBan:      autoBan,
			IsMultiKey:   isMultiKey,
			Original:     original,
		}

		if isMultiKey {
			p.MultiKeyIndex = i + 1
			p.OriginalName = ch.Name
		}

		providers = append(providers, p)
	}

	return providers
}

// checkUnmappableFields checks for fields that cannot be mapped and adds warnings.
func (e *Exporter) checkUnmappableFields(ch Channel) {
	if ch.Priority != nil && *ch.Priority != 0 {
		e.result.AddWarning(fmt.Sprintf(
			"Channel '%s' (ID=%d) has priority=%d which is not supported in EZ-API",
			ch.Name, ch.ID, *ch.Priority,
		))
	}

	if ch.ModelMapping != nil && *ch.ModelMapping != "" {
		e.result.AddWarning(fmt.Sprintf(
			"Channel '%s' (ID=%d) has model_mapping which is not migrated. Use EZ-API Binding instead.",
			ch.Name, ch.ID,
		))
	}

	if ch.StatusCodeMapping != nil && *ch.StatusCodeMapping != "" {
		e.result.AddWarning(fmt.Sprintf(
			"Channel '%s' (ID=%d) has status_code_mapping which is not supported in EZ-API",
			ch.Name, ch.ID,
		))
	}

	if ch.Setting != nil && *ch.Setting != "" {
		e.result.AddWarning(fmt.Sprintf(
			"Channel '%s' (ID=%d) has custom settings which are not migrated",
			ch.Name, ch.ID,
		))
	}

	if ch.ParamOverride != nil && *ch.ParamOverride != "" {
		e.result.AddWarning(fmt.Sprintf(
			"Channel '%s' (ID=%d) has param_override which is not supported in EZ-API",
			ch.Name, ch.ID,
		))
	}

	if ch.HeaderOverride != nil && *ch.HeaderOverride != "" {
		e.result.AddWarning(fmt.Sprintf(
			"Channel '%s' (ID=%d) has header_override which is not supported in EZ-API",
			ch.Name, ch.ID,
		))
	}

	// Multi-group warning
	groups := parseGroups(ch.Group)
	if len(groups) > 1 {
		e.result.AddWarning(fmt.Sprintf(
			"Channel '%s' (ID=%d) belongs to multiple groups %v. Only '%s' is used as primary group. Consider creating Bindings for other groups.",
			ch.Name, ch.ID, groups, groups[0],
		))
	}
}

// createOriginalBackup creates a JSON backup of original channel data.
func (e *Exporter) createOriginalBackup(ch Channel) json.RawMessage {
	data, err := json.Marshal(ch)
	if err != nil {
		return nil
	}
	return data
}

// exportUsersAndTokens exports users and tokens as masters and keys.
func (e *Exporter) exportUsersAndTokens() error {
	// Get all users with tokens
	users, err := e.connector.GetUsersWithTokens()
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	// Create a map to track masters
	masterMap := make(map[int]string) // user_id -> master_name

	for _, user := range users {
		// Create master from user
		master := e.userToMaster(user)
		e.result.AddMaster(master)
		masterMap[user.ID] = master.Name

		// Get tokens for this user
		tokens, err := e.connector.GetTokensByUserID(user.ID)
		if err != nil {
			e.result.AddWarning(fmt.Sprintf(
				"Failed to get tokens for user '%s' (ID=%d): %v",
				user.Username, user.ID, err,
			))
			continue
		}

		// Convert tokens to keys
		for _, token := range tokens {
			key := e.tokenToKey(token, master.Name)
			e.result.AddKey(key)
		}
	}

	return nil
}

// userToMaster converts a New API user to an EZ-API master.
func (e *Exporter) userToMaster(user User) schema.Master {
	return schema.Master{
		Name:             user.Username,
		Group:            user.Group,
		Namespaces:       []string{user.Group},
		DefaultNamespace: user.Group,
		MaxChildKeys:     10, // Default value
		GlobalQPS:        3,  // Default value
		Status:           MapUserStatus(user.Status),
		SourceUserID:     user.ID,
		SourceEmail:      user.Email,
	}
}

// tokenToKey converts a New API token to an EZ-API key.
func (e *Exporter) tokenToKey(token Token, masterRef string) schema.Key {
	key := schema.Key{
		MasterRef:               masterRef,
		OriginalToken:           token.Key,
		Group:                   token.Group,
		Status:                  MapTokenStatus(token.Status),
		Scopes:                  []string{"chat:*", "completions:*"},
		Namespaces:              []string{token.Group},
		ModelLimitsEnabled:      token.ModelLimitsEnabled,
		UnlimitedQuota:          token.UnlimitedQuota,
		OriginalID:              token.ID,
		TokenPlaintextAvailable: true,
	}

	// Parse model limits
	if token.ModelLimitsEnabled && token.ModelLimits != "" {
		key.ModelLimits = parseModels(token.ModelLimits)
	}

	// Parse expiration time
	if token.ExpiredTime > 0 && token.ExpiredTime != -1 {
		key.ExpiresAt = TimestampToTime(token.ExpiredTime)
	}

	// Parse IP whitelist
	if token.AllowIPs != nil && *token.AllowIPs != "" {
		key.AllowIPs = strings.Split(*token.AllowIPs, ",")
		for i := range key.AllowIPs {
			key.AllowIPs[i] = strings.TrimSpace(key.AllowIPs[i])
		}
	}

	// Set quota
	if !token.UnlimitedQuota {
		quota := int64(token.RemainQuota)
		key.QuotaLimit = &quota
		used := int64(token.UsedQuota)
		key.QuotaUsed = &used
	}

	return key
}

// exportAbilities exports abilities as bindings.
func (e *Exporter) exportAbilities() error {
	abilities, err := e.connector.GetAllAbilities()
	if err != nil {
		return err
	}

	for _, ab := range abilities {
		binding := e.abilityToBinding(ab)
		e.result.AddBinding(binding)
	}

	return nil
}

// abilityToBinding converts a New API ability to an EZ-API binding.
func (e *Exporter) abilityToBinding(ab Ability) schema.Binding {
	status := "active"
	if !ab.Enabled {
		status = "disabled"
	}

	return schema.Binding{
		Namespace:  ab.Group,
		RouteGroup: ab.Group, // Use group as route_group
		Model:      ab.Model,
		Status:     status,
	}
}

// ============================================
// Helper Functions
// ============================================

// parseKeys parses the key field which may contain multiple keys separated by newlines.
func parseKeys(key string) []string {
	if key == "" {
		return nil
	}

	lines := strings.Split(key, "\n")
	var keys []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			keys = append(keys, line)
		}
	}

	if len(keys) == 0 {
		return []string{key}
	}

	return keys
}

// parseGroups parses the group field which may contain multiple groups separated by commas.
func parseGroups(group string) []string {
	if group == "" {
		return []string{"default"}
	}

	parts := strings.Split(group, ",")
	var groups []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			groups = append(groups, part)
		}
	}

	if len(groups) == 0 {
		return []string{"default"}
	}

	return groups
}

// parseModels parses the models field which contains comma-separated model names.
func parseModels(models string) []string {
	if models == "" {
		return nil
	}

	parts := strings.Split(models, ",")
	var result []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

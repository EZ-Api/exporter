// Package schema defines the intermediate JSON format for migration.
// This format is used as a bridge between New API and EZ-API.
// Reference: SPEC_newapi_migration_tool.md Section 7
package schema

import (
	"encoding/json"
	"time"
)

// ExportResult represents the complete export output.
type ExportResult struct {
	Version  string   `json:"version"`            // Schema version, e.g., "1.0.0"
	Source   Source   `json:"source"`             // Source system information
	Data     Data     `json:"data"`               // Exported data
	Warnings []string `json:"warnings,omitempty"` // Export warnings
}

// Source represents the source system information.
type Source struct {
	Type       string    `json:"type"`        // Source type, e.g., "newapi"
	Version    string    `json:"version"`     // Source version if detectable
	ExportedAt time.Time `json:"exported_at"` // Export timestamp
}

// Data contains all exported entities.
type Data struct {
	Providers []Provider `json:"providers,omitempty"`
	Masters   []Master   `json:"masters,omitempty"`
	Keys      []Key      `json:"keys,omitempty"`
	Bindings  []Binding  `json:"bindings,omitempty"`
}

// Provider represents an EZ-API provider (mapped from New API channel).
type Provider struct {
	OriginalID   int      `json:"original_id"`          // Original channel ID
	Name         string   `json:"name"`                 // Provider name
	Type         string   `json:"type"`                 // Provider type (mapped from int)
	BaseURL      string   `json:"base_url,omitempty"`   // Custom base URL
	APIKey       string   `json:"api_key"`              // API key (single key)
	Models       []string `json:"models,omitempty"`     // Supported models
	PrimaryGroup string   `json:"primary_group"`        // Primary group (first from multi-group)
	AllGroups    []string `json:"all_groups,omitempty"` // All groups (for multi-group channels)
	Weight       int      `json:"weight"`               // Load balancing weight
	Priority     int      `json:"priority,omitempty"`   // Channel priority (optional, fallback for weight)
	Status       string   `json:"status"`               // active/disabled
	AutoBan      bool     `json:"auto_ban"`             // Auto ban on failure

	// Multi-key tracking
	IsMultiKey    bool   `json:"is_multi_key,omitempty"`    // Was this from a multi-key channel
	MultiKeyIndex int    `json:"multi_key_index,omitempty"` // Index in multi-key split (1-based)
	OriginalName  string `json:"original_name,omitempty"`   // Original channel name before split

	// Original data backup (for fields that cannot be mapped)
	Original json.RawMessage `json:"_original,omitempty"`
}

// Master represents an EZ-API master (inferred from New API user).
type Master struct {
	Name             string   `json:"name"`                        // Master name (from username)
	Group            string   `json:"group"`                       // User group
	Namespaces       []string `json:"namespaces,omitempty"`        // Accessible namespaces
	DefaultNamespace string   `json:"default_namespace,omitempty"` // Default namespace
	MaxChildKeys     int      `json:"max_child_keys,omitempty"`    // Max child keys allowed
	GlobalQPS        int      `json:"global_qps,omitempty"`        // Global QPS limit
	Status           string   `json:"status"`                      // active/suspended

	// Source tracking
	SourceUserID int    `json:"_source_user_id"` // Original user ID
	SourceEmail  string `json:"_source_email,omitempty"`
}

// Key represents an EZ-API key (mapped from New API token).
type Key struct {
	MasterRef     string `json:"master_ref"`      // Reference to master name
	OriginalToken string `json:"original_token"`  // Original token (plaintext)
	Group         string `json:"group,omitempty"` // Token group
	Status        string `json:"status"`          // active/disabled/expired/exhausted

	// Access control
	Scopes     []string `json:"scopes,omitempty"`     // Permission scopes
	Namespaces []string `json:"namespaces,omitempty"` // Accessible namespaces

	// Model limits
	ModelLimitsEnabled bool     `json:"model_limits_enabled,omitempty"`
	ModelLimits        []string `json:"model_limits,omitempty"` // Allowed models

	// Expiration
	ExpiresAt *time.Time `json:"expires_at,omitempty"` // Expiration time

	// IP restrictions
	AllowIPs []string `json:"allow_ips,omitempty"` // IP whitelist

	// Quota (if supported)
	QuotaLimit     *int64 `json:"quota_limit,omitempty"`     // -1 = unlimited
	QuotaUsed      *int64 `json:"quota_used,omitempty"`      // Used quota
	UnlimitedQuota bool   `json:"unlimited_quota,omitempty"` // Unlimited flag

	// Source tracking
	OriginalID              int  `json:"_original_id"`                         // Original token ID
	TokenPlaintextAvailable bool `json:"_token_plaintext_available,omitempty"` // Was plaintext available
}

// Binding represents an EZ-API binding (optional, from abilities).
type Binding struct {
	Namespace  string `json:"namespace"`   // Namespace (from group)
	RouteGroup string `json:"route_group"` // Route group
	Model      string `json:"model"`       // Model name
	Status     string `json:"status"`      // active/disabled
}

// NewExportResult creates a new export result with default values.
func NewExportResult() *ExportResult {
	return &ExportResult{
		Version: "1.0.0",
		Source: Source{
			Type:       "newapi",
			Version:    "unknown",
			ExportedAt: time.Now(),
		},
		Data:     Data{},
		Warnings: []string{},
	}
}

// AddWarning adds a warning message to the export result.
func (r *ExportResult) AddWarning(msg string) {
	r.Warnings = append(r.Warnings, msg)
}

// AddProvider adds a provider to the export result.
func (r *ExportResult) AddProvider(p Provider) {
	r.Data.Providers = append(r.Data.Providers, p)
}

// AddMaster adds a master to the export result.
func (r *ExportResult) AddMaster(m Master) {
	r.Data.Masters = append(r.Data.Masters, m)
}

// AddKey adds a key to the export result.
func (r *ExportResult) AddKey(k Key) {
	r.Data.Keys = append(r.Data.Keys, k)
}

// AddBinding adds a binding to the export result.
func (r *ExportResult) AddBinding(b Binding) {
	r.Data.Bindings = append(r.Data.Bindings, b)
}

// ToJSON serializes the export result to JSON.
func (r *ExportResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// Summary returns a summary of exported entities.
type Summary struct {
	Providers int `json:"providers"`
	Masters   int `json:"masters"`
	Keys      int `json:"keys"`
	Bindings  int `json:"bindings"`
	Warnings  int `json:"warnings"`
}

// GetSummary returns a summary of the export result.
func (r *ExportResult) GetSummary() Summary {
	return Summary{
		Providers: len(r.Data.Providers),
		Masters:   len(r.Data.Masters),
		Keys:      len(r.Data.Keys),
		Bindings:  len(r.Data.Bindings),
		Warnings:  len(r.Warnings),
	}
}

// Package newapi defines data models for New API database tables.
// These structures are based on the actual New API source code.
// Reference: SPEC_newapi_migration_tool.md Appendix A
package newapi

import (
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Channel represents the channels table in New API.
// Source: model/channel.go
type Channel struct {
	ID                 int          `json:"id" gorm:"primaryKey"`
	Type               int          `json:"type" gorm:"default:0"`                           // Channel type enum, see channel_type.go
	Key                string       `json:"key" gorm:"not null"`                             // API key, supports multi-key (newline separated)
	OpenAIOrganization *string      `json:"openai_organization"`                             // OpenAI organization ID
	TestModel          *string      `json:"test_model"`                                      // Test model
	Status             int          `json:"status" gorm:"default:1"`                         // 1=enabled, 2=manually disabled, 3=auto disabled
	Name               string       `json:"name" gorm:"index"`                               // Channel name
	Weight             *uint        `json:"weight" gorm:"default:0"`                         // Weight for load balancing
	CreatedTime        int64        `json:"created_time" gorm:"bigint"`                      // Creation timestamp
	TestTime           int64        `json:"test_time" gorm:"bigint"`                         // Last test timestamp
	ResponseTime       int          `json:"response_time"`                                   // Response time in ms
	BaseURL            *string      `json:"base_url" gorm:"column:base_url;default:''"`      // Custom base URL
	Other              string       `json:"other"`                                           // Other config (legacy)
	Balance            float64      `json:"balance"`                                         // Balance in USD
	BalanceUpdatedTime int64        `json:"balance_updated_time" gorm:"bigint"`              // Balance update timestamp
	Models             string       `json:"models"`                                          // Supported models (comma separated)
	Group              string       `json:"group" gorm:"type:varchar(64);default:'default'"` // Groups (comma separated for multi-group)
	UsedQuota          int64        `json:"used_quota" gorm:"bigint;default:0"`              // Used quota
	ModelMapping       *string      `json:"model_mapping" gorm:"type:text"`                  // Model mapping JSON
	StatusCodeMapping  *string      `json:"status_code_mapping" gorm:"type:varchar(1024)"`   // Status code mapping
	Priority           *int64       `json:"priority" gorm:"bigint;default:0"`                // Priority
	AutoBan            *int         `json:"auto_ban" gorm:"default:1"`                       // Auto ban on failure
	OtherInfo          string       `json:"other_info"`                                      // Other info JSON
	Tag                *string      `json:"tag" gorm:"index"`                                // Tag
	Setting            *string      `json:"setting" gorm:"type:text"`                        // Extra settings JSON
	ParamOverride      *string      `json:"param_override" gorm:"type:text"`                 // Param override JSON
	HeaderOverride     *string      `json:"header_override" gorm:"type:text"`                // Header override JSON
	Remark             *string      `json:"remark" gorm:"type:varchar(255)"`                 // Remark
	ChannelInfo        *ChannelInfo `json:"channel_info" gorm:"type:json"`                   // Multi-key management info
	OtherSettings      string       `json:"settings" gorm:"column:settings"`                 // Other settings (e.g., Azure version)
}

// TableName returns the table name for Channel.
func (Channel) TableName() string {
	return "channels"
}

// ChannelInfo represents multi-key management structure.
type ChannelInfo struct {
	IsMultiKey             bool           `json:"is_multi_key"`
	MultiKeySize           int            `json:"multi_key_size"`
	MultiKeyStatusList     map[int]int    `json:"multi_key_status_list"`
	MultiKeyDisabledReason map[int]string `json:"multi_key_disabled_reason,omitempty"`
	MultiKeyDisabledTime   map[int]int64  `json:"multi_key_disabled_time,omitempty"`
	MultiKeyPollingIndex   int            `json:"multi_key_polling_index"`
	MultiKeyMode           string         `json:"multi_key_mode"` // "random" or "polling"
}

// Scan implements sql.Scanner for ChannelInfo.
func (c *ChannelInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return nil
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, c)
}

// Token represents the tokens table in New API.
// Source: model/token.go
type Token struct {
	ID                 int            `json:"id" gorm:"primaryKey"`
	UserID             int            `json:"user_id" gorm:"index"`                   // Associated user
	Key                string         `json:"key" gorm:"type:char(48);uniqueIndex"`   // Plaintext, 48 char random string
	Status             int            `json:"status" gorm:"default:1"`                // 1=enabled, 2=disabled, 3=expired, 4=exhausted
	Name               string         `json:"name" gorm:"index"`                      // Token name
	CreatedTime        int64          `json:"created_time" gorm:"bigint"`             // Creation timestamp
	AccessedTime       int64          `json:"accessed_time" gorm:"bigint"`            // Last access timestamp
	ExpiredTime        int64          `json:"expired_time" gorm:"bigint;default:-1"`  // Expiration timestamp, -1 = never
	RemainQuota        int            `json:"remain_quota" gorm:"default:0"`          // Remaining quota
	UnlimitedQuota     bool           `json:"unlimited_quota"`                        // Unlimited quota flag
	ModelLimitsEnabled bool           `json:"model_limits_enabled"`                   // Model limits enabled
	ModelLimits        string         `json:"model_limits" gorm:"type:varchar(1024)"` // Allowed models (comma separated)
	AllowIPs           *string        `json:"allow_ips" gorm:"default:''"`            // IP whitelist
	UsedQuota          int            `json:"used_quota" gorm:"default:0"`            // Used quota
	Group              string         `json:"group" gorm:"default:''"`                // Group
	CrossGroupRetry    bool           `json:"cross_group_retry" gorm:"default:false"` // Cross-group retry (only for auto group)
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`                         // Soft delete
}

// TableName returns the table name for Token.
func (Token) TableName() string {
	return "tokens"
}

// User represents the users table in New API.
// Source: model/user.go
type User struct {
	ID              int            `json:"id" gorm:"primaryKey"`
	Username        string         `json:"username" gorm:"unique;index"`                    // Username
	Password        string         `json:"password" gorm:"not null"`                        // Password (bcrypt hash)
	DisplayName     string         `json:"display_name" gorm:"index"`                       // Display name
	Role            int            `json:"role" gorm:"type:int;default:1"`                  // 0=guest, 1=common, 10=admin, 100=root
	Status          int            `json:"status" gorm:"type:int;default:1"`                // 1=enabled, 2=disabled
	Email           string         `json:"email" gorm:"index"`                              // Email
	GitHubID        string         `json:"github_id" gorm:"column:github_id;index"`         // GitHub OAuth ID
	DiscordID       string         `json:"discord_id" gorm:"column:discord_id;index"`       // Discord OAuth ID
	OidcID          string         `json:"oidc_id" gorm:"column:oidc_id;index"`             // OIDC ID
	WeChatID        string         `json:"wechat_id" gorm:"column:wechat_id;index"`         // WeChat ID
	TelegramID      string         `json:"telegram_id" gorm:"column:telegram_id;index"`     // Telegram ID
	AccessToken     *string        `json:"access_token" gorm:"type:char(32);uniqueIndex"`   // System admin token
	Quota           int            `json:"quota" gorm:"type:int;default:0"`                 // Remaining quota
	UsedQuota       int            `json:"used_quota" gorm:"type:int;default:0"`            // Used quota
	RequestCount    int            `json:"request_count" gorm:"type:int;default:0"`         // Request count
	Group           string         `json:"group" gorm:"type:varchar(64);default:'default'"` // User group
	AffCode         string         `json:"aff_code" gorm:"type:varchar(32);uniqueIndex"`    // Affiliate code
	AffCount        int            `json:"aff_count" gorm:"type:int;default:0"`             // Affiliate count
	AffQuota        int            `json:"aff_quota" gorm:"type:int;default:0"`             // Affiliate remaining quota
	AffHistoryQuota int            `json:"aff_history_quota" gorm:"type:int;default:0"`     // Affiliate history quota
	InviterID       int            `json:"inviter_id" gorm:"type:int;index"`                // Inviter ID
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`                                  // Soft delete
	LinuxDOID       string         `json:"linux_do_id" gorm:"column:linux_do_id;index"`     // LinuxDO ID
	Setting         string         `json:"setting" gorm:"type:text"`                        // User settings JSON
	Remark          string         `json:"remark,omitempty" gorm:"type:varchar(255)"`       // Remark
	StripeCustomer  string         `json:"stripe_customer" gorm:"type:varchar(64);index"`   // Stripe customer ID
}

// TableName returns the table name for User.
func (User) TableName() string {
	return "users"
}

// Ability represents the abilities table in New API.
// Source: model/ability.go
// Note: This table is auto-generated from Channel's models and group fields.
type Ability struct {
	Group     string  `json:"group" gorm:"type:varchar(64);primaryKey"`  // Group name (composite PK)
	Model     string  `json:"model" gorm:"type:varchar(255);primaryKey"` // Model name (composite PK)
	ChannelID int     `json:"channel_id" gorm:"primaryKey;index"`        // Channel ID (composite PK)
	Enabled   bool    `json:"enabled"`                                   // Enabled flag
	Priority  *int64  `json:"priority" gorm:"bigint;default:0;index"`    // Priority
	Weight    uint    `json:"weight" gorm:"default:0;index"`             // Weight
	Tag       *string `json:"tag" gorm:"index"`                          // Tag
}

// TableName returns the table name for Ability.
func (Ability) TableName() string {
	return "abilities"
}

// Redemption represents the redemptions table in New API.
// Source: model/redemption.go
type Redemption struct {
	ID           int            `json:"id" gorm:"primaryKey"`
	UserID       int            `json:"user_id"`                              // Creator user ID
	Key          string         `json:"key" gorm:"type:char(32);uniqueIndex"` // Redemption code (32 chars)
	Status       int            `json:"status" gorm:"default:1"`              // 1=enabled, 2=disabled, 3=used
	Name         string         `json:"name" gorm:"index"`                    // Redemption code name
	Quota        int            `json:"quota" gorm:"default:100"`             // Quota value
	CreatedTime  int64          `json:"created_time" gorm:"bigint"`           // Creation timestamp
	RedeemedTime int64          `json:"redeemed_time" gorm:"bigint"`          // Redemption timestamp
	UsedUserID   int            `json:"used_user_id"`                         // User ID who used the code
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`                       // Soft delete
	ExpiredTime  int64          `json:"expired_time" gorm:"bigint"`           // Expiration timestamp, 0=never
}

// TableName returns the table name for Redemption.
func (Redemption) TableName() string {
	return "redemptions"
}

// NullableInt64 is a helper for nullable int64 fields.
type NullableInt64 struct {
	sql.NullInt64
}

// NullableString is a helper for nullable string fields.
type NullableString struct {
	sql.NullString
}

// TimestampToTime converts a Unix timestamp to time.Time pointer.
// Returns nil if timestamp is -1 (meaning never expires).
func TimestampToTime(ts int64) *time.Time {
	if ts == -1 || ts == 0 {
		return nil
	}
	t := time.Unix(ts, 0)
	return &t
}

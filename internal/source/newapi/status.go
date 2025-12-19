// Package newapi provides status enum mapping.
// Reference: SPEC_newapi_migration_tool.md Appendix C
package newapi

// UserStatus represents user status enum in New API.
type UserStatus int

const (
	UserStatusEnabled  UserStatus = 1
	UserStatusDisabled UserStatus = 2
)

// ToEzAPIStatus converts New API user status to EZ-API status string.
func (s UserStatus) ToEzAPIStatus() string {
	switch s {
	case UserStatusEnabled:
		return "active"
	case UserStatusDisabled:
		return "suspended"
	default:
		return "suspended"
	}
}

// TokenStatus represents token status enum in New API.
type TokenStatus int

const (
	TokenStatusEnabled   TokenStatus = 1
	TokenStatusDisabled  TokenStatus = 2
	TokenStatusExpired   TokenStatus = 3
	TokenStatusExhausted TokenStatus = 4
)

// ToEzAPIStatus converts New API token status to EZ-API status string.
func (s TokenStatus) ToEzAPIStatus() string {
	switch s {
	case TokenStatusEnabled:
		return "active"
	case TokenStatusDisabled:
		return "disabled"
	case TokenStatusExpired:
		return "expired"
	case TokenStatusExhausted:
		return "exhausted"
	default:
		return "disabled"
	}
}

// IsActive returns true if token is in active state.
func (s TokenStatus) IsActive() bool {
	return s == TokenStatusEnabled
}

// ChannelStatus represents channel status enum in New API.
type ChannelStatus int

const (
	ChannelStatusUnknown          ChannelStatus = 0
	ChannelStatusEnabled          ChannelStatus = 1
	ChannelStatusManuallyDisabled ChannelStatus = 2
	ChannelStatusAutoDisabled     ChannelStatus = 3
)

// ToEzAPIStatus converts New API channel status to EZ-API status string.
func (s ChannelStatus) ToEzAPIStatus() string {
	switch s {
	case ChannelStatusEnabled:
		return "active"
	case ChannelStatusManuallyDisabled, ChannelStatusAutoDisabled, ChannelStatusUnknown:
		return "disabled"
	default:
		return "disabled"
	}
}

// IsActive returns true if channel is in active state.
func (s ChannelStatus) IsActive() bool {
	return s == ChannelStatusEnabled
}

// RedemptionStatus represents redemption code status enum in New API.
type RedemptionStatus int

const (
	RedemptionStatusEnabled  RedemptionStatus = 1
	RedemptionStatusDisabled RedemptionStatus = 2
	RedemptionStatusUsed     RedemptionStatus = 3
)

// UserRole represents user role enum in New API.
type UserRole int

const (
	RoleGuestUser  UserRole = 0
	RoleCommonUser UserRole = 1
	RoleAdminUser  UserRole = 10
	RoleRootUser   UserRole = 100
)

// IsAdmin returns true if user is admin or root.
func (r UserRole) IsAdmin() bool {
	return r >= RoleAdminUser
}

// MapUserStatus maps integer user status to EZ-API status string.
func MapUserStatus(status int) string {
	return UserStatus(status).ToEzAPIStatus()
}

// MapTokenStatus maps integer token status to EZ-API status string.
func MapTokenStatus(status int) string {
	return TokenStatus(status).ToEzAPIStatus()
}

// MapChannelStatus maps integer channel status to EZ-API status string.
func MapChannelStatus(status int) string {
	return ChannelStatus(status).ToEzAPIStatus()
}

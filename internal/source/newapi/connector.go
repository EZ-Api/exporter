// Package newapi provides database connectors for New API.
package newapi

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectorType represents the database type.
type ConnectorType string

const (
	ConnectorTypeMySQL  ConnectorType = "mysql"
	ConnectorTypeSQLite ConnectorType = "sqlite"
)

// ConnectorConfig holds database connection configuration.
type ConnectorConfig struct {
	Type     ConnectorType
	DSN      string // MySQL: "user:pass@tcp(host:port)/dbname", SQLite: file path
	LogLevel logger.LogLevel
}

// Connector provides database access for New API.
type Connector struct {
	db     *gorm.DB
	config ConnectorConfig
}

// NewConnector creates a new database connector.
func NewConnector(config ConnectorConfig) (*Connector, error) {
	var dialector gorm.Dialector

	switch config.Type {
	case ConnectorTypeMySQL:
		dialector = mysql.Open(config.DSN)
	case ConnectorTypeSQLite:
		dialector = sqlite.Open(config.DSN)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Connector{
		db:     db,
		config: config,
	}, nil
}

// NewMySQLConnector creates a MySQL connector with the given DSN.
func NewMySQLConnector(dsn string) (*Connector, error) {
	return NewConnector(ConnectorConfig{
		Type:     ConnectorTypeMySQL,
		DSN:      dsn,
		LogLevel: logger.Silent,
	})
}

// NewSQLiteConnector creates a SQLite connector with the given file path.
func NewSQLiteConnector(path string) (*Connector, error) {
	return NewConnector(ConnectorConfig{
		Type:     ConnectorTypeSQLite,
		DSN:      path,
		LogLevel: logger.Silent,
	})
}

// Close closes the database connection.
func (c *Connector) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the underlying GORM database instance.
func (c *Connector) GetDB() *gorm.DB {
	return c.db
}

// ============================================
// Channel Operations
// ============================================

// GetAllChannels retrieves all channels from the database.
func (c *Connector) GetAllChannels() ([]Channel, error) {
	var channels []Channel
	err := c.db.Find(&channels).Error
	return channels, err
}

// GetChannelByID retrieves a channel by ID.
func (c *Connector) GetChannelByID(id int) (*Channel, error) {
	var channel Channel
	err := c.db.First(&channel, id).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// GetActiveChannels retrieves all active (enabled) channels.
func (c *Connector) GetActiveChannels() ([]Channel, error) {
	var channels []Channel
	err := c.db.Where("status = ?", ChannelStatusEnabled).Find(&channels).Error
	return channels, err
}

// CountChannels returns the total number of channels.
func (c *Connector) CountChannels() (int64, error) {
	var count int64
	err := c.db.Model(&Channel{}).Count(&count).Error
	return count, err
}

// ============================================
// Token Operations
// ============================================

// GetAllTokens retrieves all tokens from the database.
func (c *Connector) GetAllTokens() ([]Token, error) {
	var tokens []Token
	err := c.db.Find(&tokens).Error
	return tokens, err
}

// GetTokenByID retrieves a token by ID.
func (c *Connector) GetTokenByID(id int) (*Token, error) {
	var token Token
	err := c.db.First(&token, id).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetTokensByUserID retrieves all tokens for a user.
func (c *Connector) GetTokensByUserID(userID int) ([]Token, error) {
	var tokens []Token
	err := c.db.Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

// GetActiveTokens retrieves all active (enabled) tokens.
func (c *Connector) GetActiveTokens() ([]Token, error) {
	var tokens []Token
	err := c.db.Where("status = ?", TokenStatusEnabled).Find(&tokens).Error
	return tokens, err
}

// CountTokens returns the total number of tokens.
func (c *Connector) CountTokens() (int64, error) {
	var count int64
	err := c.db.Model(&Token{}).Count(&count).Error
	return count, err
}

// ============================================
// User Operations
// ============================================

// GetAllUsers retrieves all users from the database.
func (c *Connector) GetAllUsers() ([]User, error) {
	var users []User
	err := c.db.Find(&users).Error
	return users, err
}

// GetUserByID retrieves a user by ID.
func (c *Connector) GetUserByID(id int) (*User, error) {
	var user User
	err := c.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetActiveUsers retrieves all active (enabled) users.
func (c *Connector) GetActiveUsers() ([]User, error) {
	var users []User
	err := c.db.Where("status = ?", UserStatusEnabled).Find(&users).Error
	return users, err
}

// GetUsersWithTokens retrieves all users who have at least one token.
func (c *Connector) GetUsersWithTokens() ([]User, error) {
	var users []User
	err := c.db.Where("id IN (SELECT DISTINCT user_id FROM tokens)").Find(&users).Error
	return users, err
}

// CountUsers returns the total number of users.
func (c *Connector) CountUsers() (int64, error) {
	var count int64
	err := c.db.Model(&User{}).Count(&count).Error
	return count, err
}

// ============================================
// Ability Operations
// ============================================

// GetAllAbilities retrieves all abilities from the database.
func (c *Connector) GetAllAbilities() ([]Ability, error) {
	var abilities []Ability
	err := c.db.Find(&abilities).Error
	return abilities, err
}

// GetAbilitiesByChannelID retrieves all abilities for a channel.
func (c *Connector) GetAbilitiesByChannelID(channelID int) ([]Ability, error) {
	var abilities []Ability
	err := c.db.Where("channel_id = ?", channelID).Find(&abilities).Error
	return abilities, err
}

// GetAbilitiesByGroup retrieves all abilities for a group.
func (c *Connector) GetAbilitiesByGroup(group string) ([]Ability, error) {
	var abilities []Ability
	err := c.db.Where("`group` = ?", group).Find(&abilities).Error
	return abilities, err
}

// CountAbilities returns the total number of abilities.
func (c *Connector) CountAbilities() (int64, error) {
	var count int64
	err := c.db.Model(&Ability{}).Count(&count).Error
	return count, err
}

// ============================================
// Utility Methods
// ============================================

// Ping tests the database connection.
func (c *Connector) Ping() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GetDatabaseStats returns database statistics.
type DatabaseStats struct {
	Channels  int64
	Tokens    int64
	Users     int64
	Abilities int64
	DBType    ConnectorType
}

// GetStats returns counts of all entities.
func (c *Connector) GetStats() (*DatabaseStats, error) {
	stats := &DatabaseStats{
		DBType: c.config.Type,
	}

	var err error
	stats.Channels, err = c.CountChannels()
	if err != nil {
		return nil, fmt.Errorf("failed to count channels: %w", err)
	}

	stats.Tokens, err = c.CountTokens()
	if err != nil {
		return nil, fmt.Errorf("failed to count tokens: %w", err)
	}

	stats.Users, err = c.CountUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	stats.Abilities, err = c.CountAbilities()
	if err != nil {
		return nil, fmt.Errorf("failed to count abilities: %w", err)
	}

	return stats, nil
}

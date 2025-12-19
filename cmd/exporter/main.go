// Package main provides the CLI entry point for the exporter.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zenfun/exporter/internal/source/newapi"
	"gorm.io/gorm/logger"
)

var (
	// Version information
	version = "0.1.0"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "exporter",
	Short: "New API to EZ-API data exporter",
	Long: `Exporter is a CLI tool for exporting data from New API database
to an intermediate JSON format that can be imported into EZ-API.

Supported data sources:
  - MySQL: Direct connection to New API MySQL database
  - SQLite: SQLite database file from New API

Example:
  # Export from MySQL
  exporter export --source-type mysql --source-dsn "user:pass@tcp(localhost:3306)/new_api" -o export.json

  # Export from SQLite
  exporter export --source-type sqlite --source-path /path/to/new_api.db -o export.json`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export data from New API database",
	Long: `Export channels, tokens, and users from New API database
to an intermediate JSON format for EZ-API import.`,
	RunE: runExport,
}

var (
	// Export command flags
	sourceType       string
	sourceDSN        string
	sourcePath       string
	outputFile       string
	includeTokens    bool
	includeAbilities bool
	dryRun           bool
	verbose          bool
)

func init() {
	// Add export command
	rootCmd.AddCommand(exportCmd)

	// Export command flags
	exportCmd.Flags().StringVar(&sourceType, "source-type", "mysql", "Database type (mysql or sqlite)")
	exportCmd.Flags().StringVar(&sourceDSN, "source-dsn", "", "MySQL DSN (user:pass@tcp(host:port)/dbname)")
	exportCmd.Flags().StringVar(&sourcePath, "source-path", "", "SQLite database file path")
	exportCmd.Flags().StringVarP(&outputFile, "output", "o", "export.json", "Output file path")
	exportCmd.Flags().BoolVar(&includeTokens, "include-tokens", true, "Include tokens in export")
	exportCmd.Flags().BoolVar(&includeAbilities, "include-abilities", false, "Include abilities (bindings) in export")
	exportCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Validate without writing output file")
	exportCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
}

func runExport(cmd *cobra.Command, args []string) error {
	// Validate flags
	if sourceType != "mysql" && sourceType != "sqlite" {
		return fmt.Errorf("invalid source-type: %s (must be 'mysql' or 'sqlite')", sourceType)
	}

	if sourceType == "mysql" && sourceDSN == "" {
		return fmt.Errorf("--source-dsn is required for MySQL")
	}

	if sourceType == "sqlite" && sourcePath == "" {
		return fmt.Errorf("--source-path is required for SQLite")
	}

	// Create connector
	var connector *newapi.Connector
	var err error

	logLevel := logger.Silent
	if verbose {
		logLevel = logger.Info
	}

	switch sourceType {
	case "mysql":
		connector, err = newapi.NewConnector(newapi.ConnectorConfig{
			Type:     newapi.ConnectorTypeMySQL,
			DSN:      sourceDSN,
			LogLevel: logLevel,
		})
	case "sqlite":
		connector, err = newapi.NewConnector(newapi.ConnectorConfig{
			Type:     newapi.ConnectorTypeSQLite,
			DSN:      sourcePath,
			LogLevel: logLevel,
		})
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer connector.Close()

	// Test connection
	if err := connector.Ping(); err != nil {
		return fmt.Errorf("database connection test failed: %w", err)
	}

	if verbose {
		fmt.Println("✓ Database connection successful")
	}

	// Get stats
	stats, err := connector.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get database stats: %w", err)
	}

	fmt.Printf("Database: %s\n", sourceType)
	fmt.Printf("  Channels:  %d\n", stats.Channels)
	fmt.Printf("  Tokens:    %d\n", stats.Tokens)
	fmt.Printf("  Users:     %d\n", stats.Users)
	fmt.Printf("  Abilities: %d\n", stats.Abilities)
	fmt.Println()

	// Create exporter
	exporter := newapi.NewExporter(connector, newapi.ExporterConfig{
		IncludeTokens:    includeTokens,
		IncludeAbilities: includeAbilities,
		Verbose:          verbose,
	})

	// Run export
	fmt.Println("Exporting data...")
	result, err := exporter.Export()
	if err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	// Print summary
	summary := result.GetSummary()
	fmt.Println()
	fmt.Println("Export Summary:")
	fmt.Printf("  Providers: %d\n", summary.Providers)
	fmt.Printf("  Masters:   %d\n", summary.Masters)
	fmt.Printf("  Keys:      %d\n", summary.Keys)
	fmt.Printf("  Bindings:  %d\n", summary.Bindings)
	fmt.Printf("  Warnings:  %d\n", summary.Warnings)

	// Print warnings
	if len(result.Warnings) > 0 {
		fmt.Println()
		fmt.Println("Warnings:")
		for i, w := range result.Warnings {
			if i >= 10 && !verbose {
				fmt.Printf("  ... and %d more (use --verbose to see all)\n", len(result.Warnings)-10)
				break
			}
			fmt.Printf("  - %s\n", w)
		}
	}

	// Write output
	if dryRun {
		fmt.Println()
		fmt.Println("Dry run complete. No file written.")
		return nil
	}

	data, err := result.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize result: %w", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Println()
	fmt.Printf("✓ Export saved to: %s\n", outputFile)

	// Print file size
	info, _ := os.Stat(outputFile)
	if info != nil {
		fmt.Printf("  File size: %s\n", formatBytes(info.Size()))
	}

	return nil
}

// formatBytes formats bytes to human readable string.
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Add stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show database statistics",
	Long:  "Connect to the New API database and show entity counts.",
	RunE:  runStats,
}

func init() {
	rootCmd.AddCommand(statsCmd)

	statsCmd.Flags().StringVar(&sourceType, "source-type", "mysql", "Database type (mysql or sqlite)")
	statsCmd.Flags().StringVar(&sourceDSN, "source-dsn", "", "MySQL DSN")
	statsCmd.Flags().StringVar(&sourcePath, "source-path", "", "SQLite database file path")
	statsCmd.Flags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
}

func runStats(cmd *cobra.Command, args []string) error {
	// Validate flags
	if sourceType != "mysql" && sourceType != "sqlite" {
		return fmt.Errorf("invalid source-type: %s", sourceType)
	}

	if sourceType == "mysql" && sourceDSN == "" {
		return fmt.Errorf("--source-dsn is required for MySQL")
	}

	if sourceType == "sqlite" && sourcePath == "" {
		return fmt.Errorf("--source-path is required for SQLite")
	}

	// Create connector
	var connector *newapi.Connector
	var err error

	switch sourceType {
	case "mysql":
		connector, err = newapi.NewMySQLConnector(sourceDSN)
	case "sqlite":
		connector, err = newapi.NewSQLiteConnector(sourcePath)
	}

	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer connector.Close()

	stats, err := connector.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	fmt.Printf("Database Statistics (%s)\n", sourceType)
	fmt.Println("========================")
	fmt.Printf("Channels:  %d\n", stats.Channels)
	fmt.Printf("Tokens:    %d\n", stats.Tokens)
	fmt.Printf("Users:     %d\n", stats.Users)
	fmt.Printf("Abilities: %d\n", stats.Abilities)

	if verbose {
		// Get more details
		fmt.Println()
		fmt.Println("Active entities:")

		activeChannels, _ := connector.GetActiveChannels()
		fmt.Printf("  Active channels: %d\n", len(activeChannels))

		activeTokens, _ := connector.GetActiveTokens()
		fmt.Printf("  Active tokens:   %d\n", len(activeTokens))

		activeUsers, _ := connector.GetActiveUsers()
		fmt.Printf("  Active users:    %d\n", len(activeUsers))
	}

	return nil
}

// Add validate command
var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate an export file",
	Long:  "Validate the structure of an export JSON file.",
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Check required fields
	requiredFields := []string{"version", "source", "data"}
	for _, field := range requiredFields {
		if _, ok := result[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Check version
	version, ok := result["version"].(string)
	if !ok {
		return fmt.Errorf("version must be a string")
	}
	fmt.Printf("Version: %s\n", version)

	// Check source
	source, ok := result["source"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("source must be an object")
	}
	fmt.Printf("Source: %s (exported at: %s)\n",
		source["type"],
		source["exported_at"],
	)

	// Check data
	data_obj, ok := result["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("data must be an object")
	}

	fmt.Println()
	fmt.Println("Data counts:")
	if providers, ok := data_obj["providers"].([]interface{}); ok {
		fmt.Printf("  Providers: %d\n", len(providers))
	}
	if masters, ok := data_obj["masters"].([]interface{}); ok {
		fmt.Printf("  Masters:   %d\n", len(masters))
	}
	if keys, ok := data_obj["keys"].([]interface{}); ok {
		fmt.Printf("  Keys:      %d\n", len(keys))
	}
	if bindings, ok := data_obj["bindings"].([]interface{}); ok {
		fmt.Printf("  Bindings:  %d\n", len(bindings))
	}

	// Check warnings
	if warnings, ok := result["warnings"].([]interface{}); ok && len(warnings) > 0 {
		fmt.Printf("\nWarnings: %d\n", len(warnings))
	}

	fmt.Println()
	fmt.Println("✓ File is valid")

	return nil
}

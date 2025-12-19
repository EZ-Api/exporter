# New API Exporter

> Export data from New API database to intermediate JSON format for EZ-API import.

## Overview

This tool exports channels, tokens, and users from a New API database to a JSON format that can be imported into EZ-API. It handles:

- **Channel → Provider** mapping with multi-key splitting
- **User/Token → Master/Key** mapping with proper relationships
- **Multi-group** channels (exports all groups as bindings)
- **Type/Status** enum conversion
- **Warning collection** for unmappable fields

## Installation

```bash
# Clone the repository
git clone https://github.com/EZ-Api/exporter.git
cd exporter

# Build
go build -o exporter ./cmd/exporter

# Or install globally
go install ./cmd/exporter
```

## Usage

### Export from MySQL

```bash
# Basic export
exporter export \
  --source-type mysql \
  --source-dsn "user:pass@tcp(localhost:3306)/new_api" \
  -o export.json

# With all options
exporter export \
  --source-type mysql \
  --source-dsn "user:pass@tcp(localhost:3306)/new_api" \
  --include-tokens=true \
  --include-abilities=false \
  --verbose \
  -o export.json
```

### Export from SQLite

```bash
# First, copy the SQLite file from Docker container if needed
docker cp newapi-container:/data/new_api.db ./new_api.db

# Then export
exporter export \
  --source-type sqlite \
  --source-path ./new_api.db \
  -o export.json
```

### Dry Run

Validate the export without writing a file:

```bash
exporter export \
  --source-type sqlite \
  --source-path ./new_api.db \
  --dry-run
```

### Show Database Statistics

```bash
exporter stats \
  --source-type mysql \
  --source-dsn "user:pass@tcp(localhost:3306)/new_api"
```

### Validate Export File

```bash
exporter validate export.json
```

## Command Reference

### `exporter export`

Export data from New API database.

| Flag | Default | Description |
|------|---------|-------------|
| `--source-type` | `mysql` | Database type (`mysql` or `sqlite`) |
| `--source-dsn` | - | MySQL DSN (required for MySQL) |
| `--source-path` | - | SQLite file path (required for SQLite) |
| `-o, --output` | `export.json` | Output file path |
| `--include-tokens` | `true` | Include tokens in export |
| `--include-abilities` | `false` | Include abilities (bindings) |
| `--dry-run` | `false` | Validate without writing |
| `--verbose` | `false` | Enable verbose output |

### `exporter stats`

Show database entity counts.

### `exporter validate [file]`

Validate an export JSON file structure.

## Output Format

The export produces a JSON file with this structure:

```json
{
  "version": "1.0.0",
  "source": {
    "type": "newapi",
    "version": "unknown",
    "exported_at": "2025-01-01T00:00:00Z"
  },
  "data": {
    "providers": [...],
    "masters": [...],
    "keys": [...],
    "bindings": [...]
  },
  "warnings": [...]
}
```

### Provider (from Channel)

```json
{
  "original_id": 1,
  "name": "openai-primary",
  "type": "openai",
  "base_url": "https://api.openai.com/v1",
  "api_key": "sk-xxx",
  "models": ["gpt-4", "gpt-3.5-turbo"],
  "primary_group": "default",
  "all_groups": ["default", "vip"],
  "weight": 1,
  "status": "active",
  "auto_ban": true,
  "is_multi_key": false,
  "_original": {...}
}
```

### Master (from User)

```json
{
  "name": "user123",
  "group": "default",
  "namespaces": ["default"],
  "default_namespace": "default",
  "max_child_keys": 10,
  "global_qps": 3,
  "status": "active",
  "_source_user_id": 123
}
```

### Key (from Token)

```json
{
  "master_ref": "user123",
  "original_token": "sk-xxxx...",
  "group": "default",
  "status": "active",
  "model_limits_enabled": true,
  "model_limits": ["gpt-4"],
  "expires_at": "2025-12-31T00:00:00Z",
  "allow_ips": ["192.168.1.0/24"],
  "_original_id": 456,
  "_token_plaintext_available": true
}
```

## Multi-Key Handling

When a New API channel has multiple keys (newline separated), the exporter splits them into multiple providers:

```
Original Channel:
  name: "openai-main"
  key: "sk-key1\nsk-key2\nsk-key3"

Exported Providers:
  1. name: "openai-main",   api_key: "sk-key1"
  2. name: "openai-main-2", api_key: "sk-key2"
  3. name: "openai-main-3", api_key: "sk-key3"
```

## Multi-Group Handling

Channels with multiple groups (comma separated) use the first group as primary:

```
Original Channel:
  group: "default,vip,enterprise"

Exported Provider:
  primary_group: "default"
  all_groups: ["default", "vip", "enterprise"]
```

A warning is generated suggesting to create Bindings for other groups.

## Warnings

The exporter generates warnings for:

- Unknown channel types (mapped to "custom")
- Multi-group channels (only first group used as primary)
- Unsupported fields (priority, model_mapping, status_code_mapping, etc.)

## Development

### Prerequisites

- Go 1.23+
- Access to New API database (MySQL or SQLite)

### Build

```bash
go build -o exporter ./cmd/exporter
```

### Test

```bash
go test ./...
```

### Project Structure

```
exporter/
├── cmd/exporter/main.go          # CLI entry point
├── internal/
│   ├── source/newapi/
│   │   ├── models.go             # New API table structures
│   │   ├── connector.go          # Database connection
│   │   ├── exporter.go           # Export logic
│   │   ├── channel_type.go       # Type enum mapping
│   │   └── status.go             # Status enum mapping
│   └── schema/
│       └── intermediate.go       # Output JSON format
├── go.mod
└── README.md
```

## Related Documentation

- [SPEC: New API → EZ-API Migration Tool](../devlog/spec/SPEC_newapi_migration_tool.md)
- [EZ-API Documentation](../ez-api/README.md)

## License

MIT
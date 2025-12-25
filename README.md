# New API 数据导出工具

> 从 New API 数据库导出数据到中间 JSON 格式，用于 EZ-API 导入。

## 概述

本工具从 New API 数据库导出 channels、tokens 和 users 数据到 JSON 格式，供 EZ-API 导入使用。主要功能：

- **Channel → Provider** 映射，支持多 Key 拆分
- **User/Token → Master/Key** 映射，自动推断关联关系
- **多分组**渠道导出（所有分组作为 bindings 导出）
- **类型/状态**枚举自动转换
- **警告收集**，标记无法映射的字段

## 安装

### 方式一：go install（推荐）

```bash
go install github.com/EZ-Api/exporter/cmd/exporter@latest
```

**卸载：**
```bash
rm -f $(go env GOPATH)/bin/exporter
# 或者如果设置了 GOBIN
rm -f $(go env GOBIN)/exporter
```

### 方式二：下载预编译二进制

从 [GitHub Releases](https://github.com/EZ-Api/exporter/releases) 下载对应平台的可执行文件：

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### 方式三：从源码编译

```bash
# 克隆仓库
git clone https://github.com/EZ-Api/exporter.git
cd exporter

# 编译
go build -o exporter ./cmd/exporter
```

## 使用方法

### 从 MySQL 导出

```bash
# 基础导出
exporter export \
  --source-type mysql \
  --source-dsn "user:pass@tcp(localhost:3306)/new_api" \
  -o export.json

# 完整参数
exporter export \
  --source-type mysql \
  --source-dsn "user:pass@tcp(localhost:3306)/new_api" \
  --include-tokens=true \
  --include-abilities=false \
  --verbose \
  -o export.json
```

### 从 SQLite 导出

```bash
# 如果是 Docker 部署，先复制 SQLite 文件
docker cp newapi-container:/data/new_api.db ./new_api.db

# 然后导出
exporter export \
  --source-type sqlite \
  --source-path ./new_api.db \
  -o export.json
```

### 空运行模式

验证导出但不写入文件：

```bash
exporter export \
  --source-type sqlite \
  --source-path ./new_api.db \
  --dry-run
```

### 查看数据库统计

```bash
exporter stats \
  --source-type mysql \
  --source-dsn "user:pass@tcp(localhost:3306)/new_api"
```

### 验证导出文件

```bash
exporter validate export.json
```

## 命令参考

### `exporter export`

从 New API 数据库导出数据。

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--source-type` | `mysql` | 数据库类型（`mysql` 或 `sqlite`） |
| `--source-dsn` | - | MySQL DSN（MySQL 必填） |
| `--source-path` | - | SQLite 文件路径（SQLite 必填） |
| `-o, --output` | `export.json` | 输出文件路径 |
| `--include-tokens` | `true` | 是否包含 tokens |
| `--include-abilities` | `false` | 是否包含 abilities（bindings） |
| `--dry-run` | `false` | 仅验证不写入 |
| `--verbose` | `false` | 详细输出 |

### `exporter stats`

显示数据库实体统计。

### `exporter validate [file]`

验证导出 JSON 文件结构。

## 输出格式

导出生成的 JSON 文件结构：

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

### Provider（来自 Channel）

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

### Master（来自 User）

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

### Key（来自 Token）

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

## 多 Key 处理

当 New API 的 channel 包含多个 key（换行分隔）时，导出工具会将它们拆分为多个 provider：

```
原始 Channel：
  name: "openai-main"
  key: "sk-key1\nsk-key2\nsk-key3"

导出后的 Providers：
  1. name: "openai-main",   api_key: "sk-key1"
  2. name: "openai-main-2", api_key: "sk-key2"
  3. name: "openai-main-3", api_key: "sk-key3"
```

## 多分组处理

包含多个分组（逗号分隔）的 channel，使用第一个分组作为主分组：

```
原始 Channel：
  group: "default,vip,enterprise"

导出后的 Provider：
  primary_group: "default"
  all_groups: ["default", "vip", "enterprise"]
```

系统会生成警告，建议为其他分组创建 Bindings。

## 警告信息

导出工具会生成以下类型的警告：

- 未知的渠道类型（映射为 "custom"）
- 多分组渠道（仅使用第一个分组作为主分组）
- 不支持的字段（priority、model_mapping、status_code_mapping 等）

## 开发

### 环境要求

- Go 1.24+
- New API 数据库访问权限（MySQL 或 SQLite）

### 编译

```bash
go build -o exporter ./cmd/exporter
```

### 测试

```bash
go test ./...
```

### 项目结构

```
exporter/
├── cmd/exporter/main.go          # CLI 入口
├── internal/
│   ├── source/newapi/
│   │   ├── models.go             # New API 表结构
│   │   ├── connector.go          # 数据库连接
│   │   ├── exporter.go           # 导出逻辑
│   │   ├── channel_type.go       # 类型枚举映射
│   │   └── status.go             # 状态枚举映射
│   └── schema/
│       └── intermediate.go       # 输出 JSON 格式定义
├── go.mod
└── README.md
```

## 相关文档

- [SPEC: New API → EZ-API 迁移工具](../devlog/spec/SPEC_newapi_migration_tool.md)
- [EZ-API 文档](../ez-api/README.md)

## 许可证

MIT
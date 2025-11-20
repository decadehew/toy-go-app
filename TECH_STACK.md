# 技術棧 (Tech Stack)

本項目採用以下技術棧進行開發：

## 後端 (Backend)

- **語言**: [Go](https://go.dev/) (v1.25.1)
- **Web 框架**: [Chi](https://github.com/go-chi/chi) (v5) - 輕量級、慣用的 Go HTTP 路由器
- **數據庫驅動**: [pgx](https://github.com/jackc/pgx) (v5) - 高性能 PostgreSQL 驅動
- **數據驗證**: [validator](https://github.com/go-playground/validator) (v10) - 結構體和字段驗證

## 數據庫 (Database)

- **PostgreSQL**: 主要關聯式數據庫 (由 `pgx` 驅動推斷)
- **遷移工具**: [golang-migrate](https://github.com/golang-migrate/migrate) - 用於數據庫版本控制 (見 `Makefile`)

## 開發工具 (Development Tools)

- **任務管理**: Makefile - 用於自動化常用命令 (測試、遷移、文檔生成)
- **熱重載 (Live Reload)**: [Air](https://github.com/air-verse/air) - 用於 Go 應用的實時重載 (見 `.air.toml`)
- **環境變量管理**: [direnv](https://direnv.net/) (推斷自 `.envrc`)
- **API 文檔**: [Swag](https://github.com/swaggo/swag) - 自動生成 Swagger 2.0 文檔 (見 `Makefile`)

## 項目結構

- `cmd/`: 應用程序入口點
- `internal/`: 私有應用程序和庫代碼
- `web/`: 前端資源 (目前為空)
- `scripts/`: 腳本文件
- `bin/`: 編譯後的二進制文件

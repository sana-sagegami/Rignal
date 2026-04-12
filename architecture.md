auto-zen-backend/
├── cmd/
│   └── api/
│       └── main.go          # 全体のエントリーポイント（起動担当）
├── internal/
│   ├── domain/              # データの設計図（モデル）
│   │   ├── user.go
│   │   └── zen_record.go
│   ├── infrastructure/      # 外部との接続（DB）
│   │   └── database.go
│   └── interfaces/          # 窓口（APIの受付担当）
│       └── handlers/
│           ├── user_handler.go
│           └── log_handler.go
├── go.mod
└── go.sum
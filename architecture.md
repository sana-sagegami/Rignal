auto-zen-backend/
├── cmd/
│   └── api/
│       └── main.go                  # アプリケーションのエントリーポイント
├── controllers/
│   └── http/
│       ├── user_controller.go       # ユーザー関連のHTTPハンドラ
│       └── log_controller.go        # ログ関連のHTTPハンドラ
├── models/
│   ├── user.go                      # ユーザーのドメインモデル
│   └── zen_record.go                # 瞑想記録のドメインモデル
├── repositories/
│   ├── user_repository.go           # ユーザーデータベースアクセス
│   └── log_repository.go            # ログデータベースアクセス
├── services/
│   ├── user_service.go              # ユーザービジネスロジック
│   └── log_service.go               # ログビジネスロジック
├── middlewares/
│   └── auth_middleware.go           # 認証ミドルウェア
├── dto/
│   └── http/                        # HTTPリクエスト/レスポンスの定義
├── infra/
│   └── database.go                  # データベース接続設定
├── migrations/                      # データベースマイグレーション
├── hub/                             # WebSocketなど共有ロジック
├── scripts/                         # 補助スクリプト
├── Dockerfile                       # Dockerコンテナ定義
├── docker-compose.yml               # Docker Compose設定
├── go.mod                           # Goモジュール定義
├── go.sum                           # Goモジュール依存関係ロック
└── Makefile                         # ビルドとデプロイメント
# ビルド用
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# 実行用
# ビルドツールが入っていない、超軽量なOS
FROM alpine:latest

# 本番環境でSSL通信が必要になった時のための証明書
RUN apk --no-code add ca-certificates

WORKDIR /app/

# ビルド用から作成したバイナリファイルだけをコピー
COPY --from=builder /app/main .
# .envファイル
COPY --from=builder /app/.env .

CMD ["./main"]
FROM golang:alpine
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/api/main.go
CMD ["./main"]
package main

import (
	"auto-zen-backend/controllers/http"
	"auto-zen-backend/infra"
	"auto-zen-backend/middlewares"
	"auto-zen-backend/repositories"
	"auto-zen-backend/services"
	"auto-zen-backend/worker"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// インフラ層
	db := infra.InitDB()

	// リポジトリ層
	logRepo := repositories.NewLogRepository(db)
	userRepo := repositories.NewUserRepository(db)

	// サービス層
	logService := services.NewLogService(logRepo)
	userService := services.NewUserService(userRepo)

	// コントローラー層
	logCtrl := http.NewLogController(logService)
	userCtrl := http.NewUserController(userService)

	// Ginルーター
	r := gin.Default()

authorized := r.Group("/")
authorized.Use(middlewares.AuthMiddleware())
{
		// ルーティング
	authorized.GET("/logs", logCtrl.GetLogs)
	authorized.POST("/save", logCtrl.SaveLog)
	authorized.DELETE("/delete", logCtrl.DeleteLog)
}
	// ユーザー
	r.POST("/signup", userCtrl.Signup)
	r.POST("/login", userCtrl.Login)

	fmt.Printf("/signup handler type: %T\n", userCtrl.Signup)

	// worker を goroutine で起動
    poller := worker.NewPoller(ouraClient, analyzerService)
    go poller.Start(ctx)

	// サーバー起動
	fmt.Println("Auto-Zen Backend is starting on :8081...")
	r.Run(":8081")
}

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	httpctrl "rignal/controllers/http"
	"rignal/infra"
	"rignal/infra/oura"
	"rignal/middlewares"
	"rignal/repositories"
	"rignal/services"
	"rignal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	// --- インフラ層 ---
	db := infra.InitDB()
	ouraClient := oura.NewClient(os.Getenv("OURA_ACCESS_TOKEN"))

	// --- リポジトリ層 ---
	logRepo := repositories.NewLogRepository(db)
	userRepo := repositories.NewUserRepository(db)
	readinessRepo := repositories.NewReadinessRepository(db)
	sleepRepo := repositories.NewSleepRepository(db)
	summaryRepo := repositories.NewSummaryRepository(db)
	ibiRepo := repositories.NewIBIRepository(db)

	// --- サービス層 ---
	logService := services.NewLogService(logRepo)
	userService := services.NewUserService(userRepo)
	analyzerService := services.NewAnalyzerService(ouraClient, readinessRepo, sleepRepo, summaryRepo, ibiRepo)

	// --- コントローラー層 ---
	logCtrl := httpctrl.NewLogController(logService)
	userCtrl := httpctrl.NewUserController(userService)
	webhookCtrl := httpctrl.NewWebhookController(analyzerService, os.Getenv("OURA_WEBHOOK_VERIFY_TOKEN"))
	summaryCtrl := httpctrl.NewSummaryController(summaryRepo)

	// --- Gin ルーター ---
	r := gin.Default()

	authorized := r.Group("/")
	authorized.Use(middlewares.AuthMiddleware())
	{
		authorized.GET("/logs", logCtrl.GetLogs)
		authorized.POST("/save", logCtrl.SaveLog)
		authorized.DELETE("/delete", logCtrl.DeleteLog)
	}

	// Swift アプリ向け（認証不要）
	r.GET("/summary", summaryCtrl.GetSummary)

	// 開発用: 手動で分析を実行して今日のサマリーを生成する
	r.POST("/admin/analyze", func(c *gin.Context) {
		yesterday := time.Now().AddDate(0, 0, -1)
		go func() {
			if err := analyzerService.RunDailyAnalysis(context.Background(), yesterday); err != nil {
				fmt.Printf("[admin/analyze] failed: %v\n", err)
			}
		}()
		c.JSON(200, gin.H{"message": "分析を開始しました（昨日のデータで今日のサマリーを生成）"})
	})

	r.POST("/signup", userCtrl.Signup)
	r.POST("/login", userCtrl.Login)
	r.POST("/webhook/oura", webhookCtrl.HandleOuraEvent)

	// --- バックグラウンドワーカー ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	poller := worker.NewPoller(ouraClient, analyzerService, ibiRepo)
	go poller.Start(ctx)

	fmt.Println("Rignal is starting on :8081...")
	r.Run(":8081")
}

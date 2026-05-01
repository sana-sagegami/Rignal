package http

import (
	dto "auto-zen-backend/dto/http"
	"auto-zen-backend/repositories"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SummaryController struct {
	summaryRepo repositories.SummaryRepository
}

func NewSummaryController(sr repositories.SummaryRepository) *SummaryController {
	return &SummaryController{summaryRepo: sr}
}

// GetSummary は指定日（省略時は今日）のコンディションサマリーを返す。
// Swift アプリが毎朝プッシュ通知後に詳細を取得するために使用する。
func (c *SummaryController) GetSummary(ctx *gin.Context) {
	date := time.Now()
	if d := ctx.Query("date"); d != "" {
		parsed, err := time.Parse("2006-01-02", d)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "date は YYYY-MM-DD 形式で指定してください"})
			return
		}
		date = parsed
	}

	summary, err := c.summaryRepo.FindByDate(date)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "指定日のサマリーが見つかりません"})
		return
	}

	ctx.JSON(http.StatusOK, dto.FromDailySummary(summary))
}


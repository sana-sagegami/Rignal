package http

import (
	"rignal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	Service services.LogService
}

func NewLogController(s services.LogService) *LogController {
	return &LogController{Service: s}
}

func (ctrl *LogController) GetLogs(c *gin.Context) {
	logs, err := ctrl.Service.GetAllLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データの取得に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func (ctrl *LogController) SaveLog(c *gin.Context) {
	var input struct {
		Task		 string `json:"task" binding:"required"`
		Duration int		`json:"duration" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "タスク名と集中時間は必須です"})
		return
	}

	if err := ctrl.Service.SaveLog(input.Task, input.Duration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ctrl *LogController) DeleteLog(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDが必要です"})
		return
	}

	if err := ctrl.Service.DeleteLog(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "削除に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

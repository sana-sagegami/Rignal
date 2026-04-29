package http

import (
	"auto-zen-backend/services"
	"encoding/json"
	"net/http"
)

type WebhookController struct {
    analyzerService services.AnalyzerService
    verifyToken     string
}

func (c *WebhookController) HandleOuraEvent(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("x-oura-verification-token") != c.verifyToken {
        http.Error(w, "unauthorized", 401)
        return
    }

    var payload dto.OuraWebhookPayload
    json.NewDecoder(r.Body).Decode(&payload)

    // readiness が来たら分析実行
    if payload.DataType == "daily_readiness" {
        go c.analyzerService.RunDailyAnalysis(payload.Day)
    }

    w.WriteHeader(204)
}
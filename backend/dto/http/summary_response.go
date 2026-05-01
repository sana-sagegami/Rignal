package http

import (
	"auto-zen-backend/models"
	"time"
)

// SummaryResponse は Swift アプリ向け /summary エンドポイントのレスポンス型。
type SummaryResponse struct {
	Date             string  `json:"date"`
	ConditionScore   int     `json:"condition_score"`
	FocusPeakStart   *string `json:"focus_peak_start"`
	FocusPeakEnd     *string `json:"focus_peak_end"`
	RecommendBedtime *string `json:"recommend_bedtime"`
	SleepDebtMinutes int     `json:"sleep_debt_minutes"`
}

// FromDailySummary は DailySummary モデルを SummaryResponse に変換する。
// 時刻は RFC3339 形式で返す（Swift の ISO8601DateFormatter に対応）。
func FromDailySummary(s *models.DailySummary) SummaryResponse {
	resp := SummaryResponse{
		Date:             s.Date.Format("2006-01-02"),
		ConditionScore:   s.ConditionScore,
		SleepDebtMinutes: s.SleepDebtMinutes,
	}
	if s.FocusPeakStart != nil {
		str := s.FocusPeakStart.Format(time.RFC3339)
		resp.FocusPeakStart = &str
	}
	if s.FocusPeakEnd != nil {
		str := s.FocusPeakEnd.Format(time.RFC3339)
		resp.FocusPeakEnd = &str
	}
	if s.RecommendBedtime != nil {
		str := s.RecommendBedtime.Format(time.RFC3339)
		resp.RecommendBedtime = &str
	}
	return resp
}

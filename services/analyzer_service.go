package services

import (
	"auto-zen-backend/infra/oura"
	"auto-zen-backend/models"
	"auto-zen-backend/repositories"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type AnalyzerService interface {
	RunDailyAnalysis(ctx context.Context, date time.Time) error
}

type analyzerService struct {
	ouraClient    *oura.Client
	readinessRepo repositories.ReadinessRepository
	sleepRepo     repositories.SleepRepository
	summaryRepo   repositories.SummaryRepository
	ibiRepo       repositories.IBIRepository
}

func NewAnalyzerService(
	c *oura.Client,
	rr repositories.ReadinessRepository,
	sr repositories.SleepRepository,
	smr repositories.SummaryRepository,
	ir repositories.IBIRepository,
) AnalyzerService {
	return &analyzerService{
		ouraClient:    c,
		readinessRepo: rr,
		sleepRepo:     sr,
		summaryRepo:   smr,
		ibiRepo:       ir,
	}
}

// RunDailyAnalysis は Oura データを取得・保存し、コンディションスコアと通知内容を算出して
// daily_summaries に保存する。Webhook と poller の両方から呼ばれる。
func (s *analyzerService) RunDailyAnalysis(ctx context.Context, date time.Time) error {
	dateStr := date.Format("2006-01-02")

	if err := s.fetchAndSaveReadiness(ctx, date, dateStr); err != nil {
		return fmt.Errorf("analyzer: %w", err)
	}
	if err := s.fetchAndSaveSleep(ctx, date, dateStr); err != nil {
		return fmt.Errorf("analyzer: %w", err)
	}

	readiness, err := s.readinessRepo.FindByDate(date)
	if err != nil {
		return fmt.Errorf("analyzer: load readiness: %w", err)
	}
	sleep, err := s.sleepRepo.FindByDate(date)
	if err != nil {
		return fmt.Errorf("analyzer: load sleep: %w", err)
	}

	conditionScore := s.calcConditionScore(date, readiness, sleep)
	summary := s.buildSummary(date, conditionScore, sleep)

	if err := s.summaryRepo.Save(summary); err != nil {
		return fmt.Errorf("analyzer: save summary: %w", err)
	}

	return nil
}

// --- Oura データ取得＆保存 ---

func (s *analyzerService) fetchAndSaveReadiness(ctx context.Context, date time.Time, dateStr string) error {
	data, err := s.ouraClient.GetDailyReadiness(ctx, dateStr)
	if err != nil {
		return fmt.Errorf("fetch readiness: %w", err)
	}
	raw, _ := json.Marshal(data)
	record := &models.ReadinessRecord{Date: date, RawJSON: string(raw)}
	if data.Score != nil {
		record.Score = *data.Score
	}
	if data.Contributors.HRVBalance != nil {
		record.HRVBalance = *data.Contributors.HRVBalance
	}
	return s.readinessRepo.Save(record)
}

func (s *analyzerService) fetchAndSaveSleep(ctx context.Context, date time.Time, dateStr string) error {
	// daily_sleep からスコアを取得
	daily, err := s.ouraClient.GetDailySleep(ctx, dateStr)
	if err != nil {
		return fmt.Errorf("fetch daily_sleep: %w", err)
	}

	raw, _ := json.Marshal(daily)
	record := &models.SleepRecord{Date: date, RawJSON: string(raw)}
	if daily.Score != nil {
		record.Score = *daily.Score
	}

	// /sleep から実際の効率・時間・起床時刻を取得して上書き
	detail, err := s.ouraClient.GetSleep(ctx, dateStr)
	if err == nil {
		if detail.Efficiency != nil {
			record.Efficiency = *detail.Efficiency
		}
		if detail.TotalSleepDuration != nil {
			record.TotalMinutes = *detail.TotalSleepDuration / 60
		}
		if wt, err := time.Parse(time.RFC3339, detail.BedtimeEnd); err == nil {
			record.WakeTime = &wt
		}
	}
	// GetSleep のエラーは非致命的（daily_sleep のデータだけでも続行）

	return s.sleepRepo.Save(record)
}

// --- コンディションスコア計算 ---

// calcConditionScore は CLAUDE.md の式に従ってコンディションスコア（0–100）を算出する。
func (s *analyzerService) calcConditionScore(
	date time.Time,
	readiness *models.ReadinessRecord,
	sleep *models.SleepRecord,
) int {
	hrvScore := s.calcHRVScore(date, readiness)
	sleepScore := calcSleepScore(sleep)

	score := float64(readiness.Score)*0.40 +
		hrvScore*0.40 +
		sleepScore*0.20

	return clamp(int(math.Round(score)), 0, 100)
}

// calcHRVScore は今日の RMSSD を過去7日平均で正規化して 0–100 に変換する。
// IBI データがない場合は readiness の HRVBalance を fallback として使う。
func (s *analyzerService) calcHRVScore(date time.Time, readiness *models.ReadinessRecord) float64 {
	todayRMSSD := s.calcRMSSD(date)
	if todayRMSSD == 0 {
		return float64(readiness.HRVBalance)
	}

	avg7 := s.calcAvgRMSSD7Days(date)
	if avg7 == 0 {
		return float64(readiness.HRVBalance)
	}

	// 過去7日平均を基準（=50点）として線形正規化
	score := todayRMSSD/avg7*50 + 50
	return math.Min(100, math.Max(0, score))
}

// calcRMSSD はその日の IBI レコードから RMSSD（単位: ms）を計算する。
// 公式: √( Σ(IBI[i]-IBI[i-1])² / N )
func (s *analyzerService) calcRMSSD(date time.Time) float64 {
	records, err := s.ibiRepo.FindByDate(date)
	if err != nil || len(records) < 2 {
		return 0
	}
	var sumSq float64
	for i := 1; i < len(records); i++ {
		diff := records[i].IntervalMs - records[i-1].IntervalMs
		sumSq += diff * diff
	}
	return math.Sqrt(sumSq / float64(len(records)-1))
}

// calcAvgRMSSD7Days は過去7日分の RMSSD を計算して平均を返す。
func (s *analyzerService) calcAvgRMSSD7Days(date time.Time) float64 {
	var sum float64
	var count int
	for i := range 7 {
		d := date.AddDate(0, 0, -(i + 1))
		rmssd := s.calcRMSSD(d)
		if rmssd > 0 {
			sum += rmssd
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// calcSleepScore は睡眠スコア（0–100）を算出する。
// 式: (efficiency×0.5 + duration_ratio×0.5) × 100
func calcSleepScore(sleep *models.SleepRecord) float64 {
	efficiencyRate := float64(sleep.Efficiency) / 100.0
	durationRate := math.Min(1.0, float64(sleep.TotalMinutes)/450.0)
	return (efficiencyRate*0.5 + durationRate*0.5) * 100
}

// --- DailySummary 組み立て ---

// buildSummary はコンディションスコアをもとに集中ピーク・睡眠負債・推奨就寝時刻を計算して
// DailySummary を返す。
func (s *analyzerService) buildSummary(date time.Time, conditionScore int, sleep *models.SleepRecord) *models.DailySummary {
	summary := &models.DailySummary{
		Date:           date,
		ConditionScore: conditionScore,
	}

	wakeTime := resolveWakeTime(date, sleep)
	peakStart, peakEnd := calcFocusPeak(conditionScore, wakeTime)
	summary.FocusPeakStart = &peakStart
	summary.FocusPeakEnd = &peakEnd

	debtMin := s.calcSleepDebt()
	summary.SleepDebtMinutes = debtMin
	bedtime := calcRecommendBedtime(date, debtMin)
	summary.RecommendBedtime = &bedtime

	return summary
}

// resolveWakeTime は SleepRecord.WakeTime を返す。
// nil（データ未取得）の場合は当日 07:00 を fallback として使う。
func resolveWakeTime(date time.Time, sleep *models.SleepRecord) time.Time {
	if sleep.WakeTime != nil {
		return *sleep.WakeTime
	}
	return time.Date(date.Year(), date.Month(), date.Day(), 7, 0, 0, 0, date.Location())
}

// calcFocusPeak はコンディションスコアに応じて集中ピーク時間帯を返す（CLAUDE.md 仕様）。
func calcFocusPeak(score int, wakeTime time.Time) (start, end time.Time) {
	switch {
	case score >= 80:
		start = wakeTime.Add(90 * time.Minute)  // +1.5h
		end = start.Add(4 * time.Hour)
	case score >= 60:
		start = wakeTime.Add(120 * time.Minute) // +2.0h
		end = start.Add(3 * time.Hour)
	default:
		start = wakeTime.Add(150 * time.Minute) // +2.5h
		end = start.Add(2 * time.Hour)
	}
	return
}

// calcSleepDebt は過去7日の睡眠記録から累積睡眠負債（分）を計算する。
// 目標睡眠時間 450分（7.5時間）を下回った分を合計する。
func (s *analyzerService) calcSleepDebt() int {
	records, err := s.sleepRepo.FindRecent(7)
	if err != nil {
		return 0
	}
	const target = 450
	var debt int
	for _, r := range records {
		if r.TotalMinutes > 0 && r.TotalMinutes < target {
			debt += target - r.TotalMinutes
		}
	}
	return debt
}

// calcRecommendBedtime は推奨就寝時刻を計算する（CLAUDE.md 仕様）。
//
//	必要睡眠時間 = 450分 + 睡眠負債(過去7日累積)/7
//	推奨就寝時刻 = 翌日07:00 - 必要睡眠時間 - 30min(バッファ)
//
// Phase 2 では翌日07:00を固定の目標起床時刻とする。
// カレンダー連携で「翌日最初の予定時刻」を使うのは将来対応。
func calcRecommendBedtime(date time.Time, debtMin int) time.Time {
	tomorrow7AM := time.Date(date.Year(), date.Month(), date.Day()+1, 7, 0, 0, 0, date.Location())
	requiredSleep := time.Duration(450+debtMin/7) * time.Minute
	return tomorrow7AM.Add(-(requiredSleep + 30*time.Minute))
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

package worker

import (
	"auto-zen-backend/infra/oura"
	"auto-zen-backend/models"
	"auto-zen-backend/repositories"
	"auto-zen-backend/services"
	"context"
	"log"
	"time"
)

type Poller struct {
	ouraClient      *oura.Client
	analyzerService services.AnalyzerService
	ibiRepo         repositories.IBIRepository
}

func NewPoller(c *oura.Client, a services.AnalyzerService, ir repositories.IBIRepository) *Poller {
	return &Poller{ouraClient: c, analyzerService: a, ibiRepo: ir}
}

// Start は毎朝 6:00 に HRV データをポーリングするスケジューラー。
// ctx がキャンセルされると停止する。
func (p *Poller) Start(ctx context.Context) {
	for {
		next := nextOccurrence(6, 0)
		select {
		case <-time.After(time.Until(next)):
			p.poll(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// nextOccurrence は今日または翌日の hour:min の時刻を返す
func nextOccurrence(hour, min int) time.Time {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}

// poll は前日の interbeat_interval を取得して DB に保存したあと、RunDailyAnalysis を呼ぶ。
// Oura の IBI データは睡眠中に記録されるため、毎朝6:00時点では前日分が確定している。
func (p *Poller) poll(ctx context.Context) {
	yesterday := time.Now().AddDate(0, 0, -1)
	dateStr := yesterday.Format("2006-01-02")

	ibi, err := p.ouraClient.GetInterbeatInterval(ctx, dateStr)
	if err != nil {
		log.Printf("[poller] GetInterbeatInterval failed (date=%s): %v", dateStr, err)
		// IBI が取れなくても他のデータは取得済みのため分析は続行
	} else {
		if err := p.saveIBIRecords(ibi, yesterday); err != nil {
			log.Printf("[poller] saveIBIRecords failed (date=%s): %v", dateStr, err)
		}
	}

	if err := p.analyzerService.RunDailyAnalysis(ctx, yesterday); err != nil {
		log.Printf("[poller] RunDailyAnalysis failed (date=%s): %v", dateStr, err)
	}
}

// saveIBIRecords は Oura API の InterbeatInterval レスポンスを IBIRecord スライスに変換して
// バッチ INSERT する。
//
// Oura API の Items は RR間隔（秒単位）の配列。
// Timestamp はその日のデータ取得開始時刻で、Interval 秒ごとに1サンプル記録される。
func (p *Poller) saveIBIRecords(ibi *oura.InterbeatInterval, date time.Time) error {
	if len(ibi.Items) == 0 {
		return nil
	}

	startTime, err := time.Parse(time.RFC3339, ibi.Timestamp)
	if err != nil {
		// Timestamp のパースに失敗した場合はその日の0時を起点として扱う
		startTime = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	}

	records := make([]models.IBIRecord, 0, len(ibi.Items))
	intervalDur := time.Duration(ibi.Interval * float64(time.Second))
	for i, item := range ibi.Items {
		records = append(records, models.IBIRecord{
			RecordedAt: startTime.Add(time.Duration(i) * intervalDur),
			IntervalMs: item * 1000, // 秒 → ミリ秒変換
		})
	}

	return p.ibiRepo.BatchInsert(records)
}

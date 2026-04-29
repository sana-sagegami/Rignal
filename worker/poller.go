package worker

import (
	"auto-zen-backend/infra/oura"
	"auto-zen-backend/services"
	"context"
	"time"
)

type Poller struct {
    ouraClient *oura.Client
    analyzerService services.AnalyzerService
}

func (p *Poller) Start(ctx context.Context) {
    // 毎朝 6:00 に HRV データを取得
    for {
        now := time.Now()
        next := nextOccurrence(6, 0) // 翌 6:00 まで待機
        select {
        case <-time.After(time.Until(next)):
            p.poll(ctx)
        case <-ctx.Done():
            return
        }
    }
}
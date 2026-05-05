package repositories

import (
	"rignal/models"
	"time"

	"gorm.io/gorm"
)

type IBIRepository interface {
	BatchInsert(records []models.IBIRecord) error
	FindByDate(date time.Time) ([]models.IBIRecord, error)
}

type ibiRepository struct {
	db *gorm.DB
}

func NewIBIRepository(db *gorm.DB) IBIRepository {
	return &ibiRepository{db: db}
}

// BatchInsert は1日分の IBI レコードを500件ずつに分けて INSERT する。
// 1日あたり数千件になるため、一括 INSERT でなくバッチに分割してメモリ使用量を抑える。
func (r *ibiRepository) BatchInsert(records []models.IBIRecord) error {
	if len(records) == 0 {
		return nil
	}
	return r.db.CreateInBatches(records, 500).Error
}

// FindByDate はその日（00:00:00〜翌日00:00:00）の IBI レコードをすべて返す。
// recorded_at の昇順で返すため、RMSSD 計算で順序依存の差分演算が正確になる。
func (r *ibiRepository) FindByDate(date time.Time) ([]models.IBIRecord, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)

	var records []models.IBIRecord
	err := r.db.
		Where("recorded_at >= ? AND recorded_at < ?", dayStart, dayEnd).
		Order("recorded_at ASC").
		Find(&records).Error
	return records, err
}

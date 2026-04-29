package models

import "time"

// IBIRecord は Oura の interbeat_interval（RR間隔）を1件ずつ保存するモデル。
// 1日数千件になるため、INSERT は BatchInsert を使うこと（CLAUDE.md 注意事項）。
type IBIRecord struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RecordedAt time.Time `gorm:"not null;index"`
	IntervalMs float64   `gorm:"not null"`
}

func (IBIRecord) TableName() string { return "ibi_records" }

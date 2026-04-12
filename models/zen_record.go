package models

import "time"

type ZenRecord struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Task      string    `json:"task"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"timestamp"`
}
package models

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Username     string `json:"username" gorm:"unique;not null"`
	PasswordHash string `json:"-"`
}
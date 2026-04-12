package models

type User struct {
	ID 					 uint   `json:"id" gorm:"primaryKey"`
	Username  	 string `json:"username" gorm:"unique;not null"`
	// json:"-" でJSONにパスワードを含めない（セキュリティ対策）
	PasswordHash string `json:"-"`
}
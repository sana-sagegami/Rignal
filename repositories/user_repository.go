package repositories

import (
	"auto-zen-backend/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// ユーザー名からユーザー情報を取得（ログイン）
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	// Firstは最初に見つかった1件を取得 見つからない場合はエラー
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
package repositories

import (
	"rignal/models"
	"gorm.io/gorm"
)

type LogRepository interface { // interfaceはテストをしやすくするため
	FindAll() ([]models.ZenRecord, error)
	Create(log *models.ZenRecord) error
	Delete(id string) error
}

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

func (r *logRepository) FindAll() ([]models.ZenRecord, error) {
	var logs []models.ZenRecord
	err := r.db.Order("created_at desc").Find(&logs).Error
	return logs, err
}

func (r *logRepository) Create(log *models.ZenRecord) error {
	return r.db.Create(log).Error
}

func (r *logRepository) Delete(id string) error {
	return r.db.Delete(&models.ZenRecord{}, id).Error
}
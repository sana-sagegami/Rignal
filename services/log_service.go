package services

import (
	"auto-zen-backend/models"
	"auto-zen-backend/repositories"
)

type LogService interface {
	GetAllLogs() ([]models.ZenRecord, error)
	SaveLog(task string, duration int) error
	DeleteLog(id string) error
}

type logService struct {
	repo repositories.LogRepository
}

func NewLogService(repo repositories.LogRepository) LogService {
	return &logService{repo: repo}
}

func (s *logService) GetAllLogs() ([]models.ZenRecord, error) {
	return s.repo.FindAll()
}

func (s *logService) SaveLog(task string, duration int) error {
	log := &models.ZenRecord{
		Task:     task,
		Duration: duration,
	}
	return s.repo.Create(log)
}

func (s *logService) DeleteLog(id string) error {
	return s.repo.Delete(id)
}
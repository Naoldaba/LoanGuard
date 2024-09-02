package repository_interface

import (
	"LoanGuard/internal/domain/models"
)

type ILogRepository interface {
	CreateLog(log *models.SystemLog) error
	GetAllLogs() ([]models.SystemLog, error)
}

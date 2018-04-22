package money

import "github.com/haziba/theplanner/models/money"

type MoneyLogService interface {
	CreateMoneyLog(models.MoneyLog) (models.MoneyLog, error)
	GetAllMoneyLogs() ([]models.MoneyLog, error)
}

package models

import "time"

type MoneyLog struct {
	ID       string    `json:"id"`
	Category string    `json:"category"`
	Label    string    `json:"label"`
	Amount   float32   `json:"amount"`
	When     time.Time `json:"time"`
}

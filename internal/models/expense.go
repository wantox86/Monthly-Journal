package models

import "time"

type Expense struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Sender      string    `json:"sender"`
	MonthYear   string    `json:"month_year"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Expense) TableName() string {
	return "expenses"
}

package models

import (
	"time"
)

type Account struct {
	AccountID string
	UserID    string
	Name      string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (Account) TableName() string {
	return "accounts"
}

type Accounts []*Account

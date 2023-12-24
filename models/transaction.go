package models

import (
	"time"
)

type Transaction struct {
	TransactionID string
	AccountID     string
	UserID        string
	Amount        float64
	Balance       float64
	Type          string
	Status        string
	Metadata      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

func (Transaction) TableName() string {
	return "transactions"
}

type Transactions []*Transaction

type TransactionMetadata struct {
	FromAccountID string `json:"from_account_id,omitempty"`
	ToAccountID   string `json:"to_account_id,omitempty"`
}

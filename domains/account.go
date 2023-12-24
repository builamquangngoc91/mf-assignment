package domains

import (
	"errors"
	"time"
)

type (
	CreateAccountRequest struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
	}

	Account struct {
		AccountID string
		UserID    string
		Name      string
		Balance   float64
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	GetAccountsResponse struct {
		Accounts   []*Account `json:"accounts"`
		NextCursor string     `json:"next_cursor"`
	}

	DepositAccountRequest struct {
		Amount float64 `json:"amount"`
	}

	DepositAccountResponse struct {
		TransactionID string `json:"transaction_id"`
	}

	WithdrawAccountRequest struct {
		Amount float64 `json:"amount"`
	}

	WithdrawAccountResponse struct {
		TransactionID string `json:"transaction_id"`
	}

	TransferAccountRequest struct {
		ToAccountID string  `json:"to_account_id"`
		Amount      float64 `json:"amount"`
	}

	TransferAccountResponse struct {
		TransactionID string `json:"transaction_id"`
	}
)

func (r *CreateAccountRequest) Validate() error {
	if r.Name == "" {
		return errors.New("missing name")
	}

	return nil
}

func (r *DepositAccountRequest) Validate() error {
	if r.Amount <= 0 {
		return errors.New("insufficient amount")
	}

	return nil
}

func (r *WithdrawAccountRequest) Validate() error {
	if r.Amount <= 0 {
		return errors.New("insufficient amount")
	}

	return nil
}

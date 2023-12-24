package domains

import "time"

type (
	Transaction struct {
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
	}

	GetTransactionsResp struct {
		Transactions []*Transaction `json:"transactions"`
		NextCursor   string         `json:"next_cursor"`
	}
)

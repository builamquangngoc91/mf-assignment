package repositories

import (
	"banking-service/models"
	"context"

	"gorm.io/gorm"
)

var _ TransactionRepositoryI = TransactionRepository{}

type (
	TransactionRepository struct{}

	GetTransactionArgs struct {
		ID string
	}
	GetTransactionsArgs struct {
		UserID    string
		AccountID string
	}

	TransactionRepositoryI interface {
		GetTransaction(context.Context, *gorm.DB, *GetTransactionArgs) (*models.Transaction, error)
		GetTransactions(context.Context, *gorm.DB, *GetTransactionsArgs) (models.Transactions, error)
		Create(context.Context, *gorm.DB, *models.Transaction) error
	}
)

func NewTransactionRepository() TransactionRepositoryI {
	return &TransactionRepository{}
}

func (TransactionRepository) Create(ctx context.Context, db *gorm.DB, account *models.Transaction) error {
	return db.
		WithContext(ctx).
		Table("transactions").
		Create(account).
		Error
}

func (TransactionRepository) GetTransaction(ctx context.Context, db *gorm.DB, args *GetTransactionArgs) (*models.Transaction, error) {
	db = db.
		WithContext(ctx).
		Table("transactions")

	if args.ID != "" {
		db.Where("id = ?", args.ID)
	}

	var transaction *models.Transaction
	result := db.First(transaction)

	return transaction, result.Error
}

func (TransactionRepository) GetTransactions(ctx context.Context, db *gorm.DB, args *GetTransactionsArgs) (transactions models.Transactions, err error) {
	db = db.
		WithContext(ctx).
		Table("transactions")

	if args.UserID != "" {
		db.Where("user_id = ?", args.UserID)
	}
	if args.AccountID != "" {
		db.Where("account_id = ?", args.AccountID)
	}
	err = db.Find(&transactions).Error

	return
}

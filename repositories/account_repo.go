package repositories

import (
	"banking-service/enums"
	"banking-service/models"
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ AccountRepositoryI = accountRepository{}

type (
	accountRepository struct{}

	GetAccountArgs struct {
		UserID    string
		AccountID string
		ForUpdate bool
	}

	GetAccountsArgs struct {
		UserID string
		Cursor string
		Limit  int
	}

	GetAccountIDsArgs struct {
		UserID string
	}

	AccountRepositoryI interface {
		GetAccount(context.Context, *gorm.DB, *GetAccountArgs) (*models.Account, error)
		GetAccounts(context.Context, *gorm.DB, *GetAccountsArgs) (models.Accounts, error)
		GetAccountIDs(context.Context, *gorm.DB, *GetAccountIDsArgs) ([]string, error)
		Create(context.Context, *gorm.DB, *models.Account) error
		Update(context.Context, *gorm.DB, *models.Account) error
	}
)

func NewAccountRepository() AccountRepositoryI {
	return &accountRepository{}
}

func (accountRepository) Create(ctx context.Context, db *gorm.DB, account *models.Account) error {
	return db.
		WithContext(ctx).
		Create(account).
		Error
}

func (accountRepository) GetAccount(ctx context.Context, db *gorm.DB, args *GetAccountArgs) (_ *models.Account, err error) {
	db = db.
		WithContext(ctx).
		Table("accounts")

	if args.AccountID != "" {
		db = db.Where("account_id = ?", args.AccountID)
	}
	if args.UserID != "" {
		db = db.Where("user_id = ?", args.UserID)
	}
	if args.ForUpdate {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	db = db.Where("deleted_at IS NULL")

	var account models.Account
	result := db.First(&account)

	return &account, result.Error
}

func (accountRepository) GetAccounts(ctx context.Context, db *gorm.DB, args *GetAccountsArgs) (accounts models.Accounts, err error) {
	db = db.
		WithContext(ctx).
		Table("accounts")

	if args.UserID != "" {
		db.Where("user_id = ?", args.UserID)
	}
	if args.Cursor != "" {
		db.Where("account_id < ?", args.Cursor)
	}
	if args.Limit == 0 {
		args.Limit = 100
	}
	db.Limit(args.Limit)
	db.Order("account_id DESC")
	db.Where("deleted_at IS NULL")
	err = db.Find(&accounts).Error

	return
}

func (accountRepository) GetAccountIDs(ctx context.Context, db *gorm.DB, args *GetAccountIDsArgs) ([]string, error) {
	db = db.
		WithContext(ctx).
		Table("accounts")

	if args.UserID != "" {
		db = db.Where("user_id = ?", args.UserID)
	}
	db = db.Where("deleted_at IS NULL")

	var accountIDs []string
	result := db.Select("account_id").Find(&accountIDs)

	return accountIDs, result.Error
}

func (accountRepository) Update(ctx context.Context, db *gorm.DB, account *models.Account) (err error) {
	db = db.
		WithContext(ctx).
		Table("accounts").
		Where("account_id = ?", account.AccountID).
		Updates(account)
	if err = db.Error; err != nil {
		return err
	}

	if db.RowsAffected == 0 {
		return errors.New(enums.NotRowsAffected)
	}

	return nil
}

func (accountRepository) Delete(ctx context.Context, db *gorm.DB, account *models.Account) (err error) {
	db = db.
		WithContext(ctx).
		Table("accounts").
		Where("account_id = ?", account.AccountID).
		Where("user_id = ?", account.UserID).
		Update("deleted_at", "NOW()")
	if err = db.Error; err != nil {
		return err
	}

	if db.RowsAffected == 0 {
		return errors.New(enums.NotRowsAffected)
	}

	return nil
}

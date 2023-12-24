package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"banking-service/domains"
	"banking-service/enums"
	"banking-service/models"
	"banking-service/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	_ AccountHandlers = &accountHandlers{}
)

type AccountHandlers interface {
	RouteGroup(r *gin.Engine)

	CreateAccountHandler(*gin.Context)
	GetAccountsHandler(*gin.Context)
	GetAccountHandler(*gin.Context)
	DepositAccountHandler(*gin.Context)
	WithdrawAccountHandler(*gin.Context)
	TransferAmountHandler(*gin.Context)
}

type AccountHandlersDeps struct {
	DB *gorm.DB
}

type accountHandlers struct {
	db                    *gorm.DB
	accountRepository     repositories.AccountRepositoryI
	transactionRepository repositories.TransactionRepositoryI
}

func NewAccountHandlers(deps *AccountHandlersDeps) AccountHandlers {
	if deps == nil {
		return nil
	}

	return &accountHandlers{
		db:                    deps.DB,
		accountRepository:     repositories.NewAccountRepository(),
		transactionRepository: repositories.NewTransactionRepository(),
	}
}

func (u *accountHandlers) RouteGroup(rg *gin.Engine) {
	rg.POST("/accounts", u.CreateAccountHandler)
	rg.GET("/accounts", u.GetAccountsHandler)
	rg.GET("/accounts/:accountID", u.GetAccountHandler)
	rg.POST("/accounts/:accountID/deposit", u.DepositAccountHandler)
	rg.POST("/accounts/:accountID/withdraw", u.WithdrawAccountHandler)
	rg.POST("/accounts/:accountID/transfer", u.TransferAmountHandler)
}

func (u *accountHandlers) CreateAccountHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req domains.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	account := &models.Account{
		AccountID: uuid.NewString(),
		UserID:    req.UserID,
		Name:      req.Name,
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.accountRepository.Create(ctx, u.db, account); err != nil {
		c.JSON(http.StatusInternalServerError, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domains.Account{
		AccountID: account.AccountID,
		UserID:    account.UserID,
		Name:      account.Name,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	})
}

func (u *accountHandlers) GetAccountsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	accounts, err := u.accountRepository.GetAccounts(ctx, u.db, &repositories.GetAccountsArgs{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	accountsResp := make([]*domains.Account, 0, len(accounts))

	for _, account := range accounts {
		accountsResp = append(accountsResp, &domains.Account{
			AccountID: account.AccountID,
			UserID:    account.UserID,
			Name:      account.Name,
			Balance:   account.Balance,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, &domains.GetAccountsResponse{
		Accounts: accountsResp,
	})
}

func (u *accountHandlers) GetAccountHandler(c *gin.Context) {
	ctx := c.Request.Context()
	accountID := c.Param("accountID")

	account, err := u.accountRepository.GetAccount(ctx, u.db, &repositories.GetAccountArgs{
		AccountID: accountID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, domains.ErrorResp{
				Message: fmt.Sprintf("account_id %s not found", accountID),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &domains.Account{
		AccountID: account.AccountID,
		UserID:    account.UserID,
		Name:      account.Name,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	})
}

func (u *accountHandlers) DepositAccountHandler(c *gin.Context) {
	ctx := c.Request.Context()
	accountID := c.Param("accountID")

	var (
		req           domains.DepositAccountRequest
		transactionID string
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	err := u.db.Transaction(func(tx *gorm.DB) error {
		account, err := u.accountRepository.GetAccount(ctx, tx, &repositories.GetAccountArgs{
			AccountID: accountID,
			ForUpdate: true,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domains.NewXError(fmt.Errorf("account_id %s not found", accountID), enums.BadRequest)
			}
			return domains.NewXError(err, enums.InternalError)
		}

		account.Balance += req.Amount
		if err := u.accountRepository.Update(ctx, tx, account); err != nil {
			return domains.NewXError(err, enums.InternalError)
		}

		transactionID = uuid.NewString()
		transaction := &models.Transaction{
			TransactionID: transactionID,
			UserID:        account.UserID,
			AccountID:     account.AccountID,
			Amount:        req.Amount,
			Balance:       account.Balance,
			Type:          enums.Deposit.String(),
			Status:        enums.Completed.String(),
			Metadata:      "{}",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := u.transactionRepository.Create(ctx, tx, transaction); err != nil {
			return domains.NewXError(err, enums.InternalError)
		}

		return nil
	})
	if err != nil {
		err.(domains.XError).Response(c)
		return
	}

	c.JSON(http.StatusOK, &domains.DepositAccountResponse{
		TransactionID: transactionID,
	})
}

func (u *accountHandlers) WithdrawAccountHandler(c *gin.Context) {
	ctx := c.Request.Context()
	accountID := c.Param("accountID")

	var (
		req           domains.WithdrawAccountRequest
		transactionID string
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	err := u.db.Transaction(func(tx *gorm.DB) error {
		account, err := u.accountRepository.GetAccount(ctx, tx, &repositories.GetAccountArgs{
			AccountID: accountID,
			ForUpdate: true,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domains.NewXError(fmt.Errorf("account_id %s not found", accountID), enums.BadRequest)
			}
			return domains.NewXError(err, enums.InternalError)
		}

		if account.Balance-req.Amount < 0 {
			return domains.NewXError(errors.New("insufficient balance"), enums.BadRequest)
		}

		account.Balance = account.Balance - req.Amount
		if err := u.accountRepository.Update(ctx, tx, account); err != nil {
			return domains.NewXError(err, enums.InternalError)
		}

		transactionID = uuid.NewString()
		transaction := &models.Transaction{
			TransactionID: transactionID,
			UserID:        account.UserID,
			AccountID:     account.AccountID,
			Amount:        -req.Amount,
			Balance:       account.Balance,
			Type:          enums.Withdrawal.String(),
			Status:        enums.Completed.String(),
			Metadata:      "{}",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := u.transactionRepository.Create(ctx, tx, transaction); err != nil {
			return domains.NewXError(err, enums.InternalError)
		}

		return nil
	})
	if err != nil {
		err.(domains.XError).Response(c)
		return
	}

	c.JSON(http.StatusOK, &domains.WithdrawAccountResponse{
		TransactionID: transactionID,
	})
}

func (u *accountHandlers) TransferAmountHandler(c *gin.Context) {
	ctx := c.Request.Context()
	accountID := c.Param("accountID")

	var (
		req           domains.TransferAccountRequest
		transactionID string
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	err := u.db.Transaction(func(tx *gorm.DB) error {
		account, err := u.accountRepository.GetAccount(ctx, tx, &repositories.GetAccountArgs{
			AccountID: accountID,
			ForUpdate: true,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domains.NewXError(fmt.Errorf("account_id %s not found", accountID), enums.BadRequest)
			}
			return domains.NewXError(err, enums.InternalError)
		}

		destinationAccount, err := u.accountRepository.GetAccount(ctx, tx, &repositories.GetAccountArgs{
			AccountID: req.ToAccountID,
			ForUpdate: true,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domains.NewXError(fmt.Errorf("destination account_id %s not found", req.ToAccountID), enums.BadRequest)
			}
			return domains.NewXError(err, enums.InternalError)
		}

		if account.Balance-req.Amount < 0 {
			return domains.NewXError(errors.New("insufficient balance"), enums.BadRequest)
		}

		account.Balance = account.Balance - req.Amount
		if err := u.accountRepository.Update(ctx, tx, account); err != nil {
			return domains.NewXError(err, enums.InternalError)
		}

		destinationAccount.Balance = destinationAccount.Balance + req.Amount
		if err := u.accountRepository.Update(ctx, tx, destinationAccount); err != nil {
			return domains.NewXError(err, enums.InternalError)
		}

		metadata := models.TransactionMetadata{
			FromAccountID: account.AccountID,
			ToAccountID:   destinationAccount.AccountID,
		}

		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return domains.NewXError(err, enums.InternalError)
		}

		transactionID = uuid.NewString()
		transaction := &models.Transaction{
			TransactionID: transactionID,
			UserID:        account.UserID,
			AccountID:     accountID,
			Amount:        -req.Amount,
			Balance:       account.Balance,
			Type:          enums.Transfer.String(),
			Status:        enums.Completed.String(),
			Metadata:      string(metadataBytes),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := u.transactionRepository.Create(ctx, tx, transaction); err != nil {
			return domains.NewXError(err, enums.InternalError)
		}

		transactionDestination := &models.Transaction{
			TransactionID: uuid.NewString(),
			UserID:        destinationAccount.UserID,
			AccountID:     destinationAccount.AccountID,
			Amount:        req.Amount,
			Balance:       destinationAccount.Balance,
			Type:          enums.Transfer.String(),
			Status:        enums.Completed.String(),
			Metadata:      string(metadataBytes),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := u.transactionRepository.Create(ctx, tx, transactionDestination); err != nil {
			return domains.NewXError(err, enums.InternalError)
		}

		return nil
	})
	if err != nil {
		err.(domains.XError).Response(c)
		return
	}

	c.JSON(http.StatusOK, &domains.TransferAccountResponse{
		TransactionID: transactionID,
	})
}

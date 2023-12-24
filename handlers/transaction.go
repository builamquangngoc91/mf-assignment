package handlers

import (
	"net/http"
	"strconv"

	"banking-service/domains"
	"banking-service/repositories"
	"banking-service/utilities"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	_ TransactionHandlers = &transactionHandlers{}
)

type TransactionHandlers interface {
	RouteGroup(r *gin.Engine)

	GetAccountTransactionsHandler(c *gin.Context)
}

type TransactionHandlersDeps struct {
	DB          *gorm.DB
	IDGenerator utilities.SnowflakeIDGenerator
}

type transactionHandlers struct {
	db                    *gorm.DB
	idGenerator           utilities.SnowflakeIDGenerator
	userRepositiory       repositories.UserRepositoryI
	accountRepository     repositories.AccountRepositoryI
	transactionRepository repositories.TransactionRepositoryI
}

func NewTransactionHandlers(deps *TransactionHandlersDeps) TransactionHandlers {
	if deps == nil {
		return nil
	}

	return &transactionHandlers{
		db:                    deps.DB,
		idGenerator:           deps.IDGenerator,
		userRepositiory:       repositories.NewUserRepository(),
		accountRepository:     repositories.NewAccountRepository(),
		transactionRepository: repositories.NewTransactionRepository(),
	}
}

func (u *transactionHandlers) RouteGroup(rg *gin.Engine) {
	rg.GET("/accounts/:accountID/transactions", u.GetAccountTransactionsHandler)
}

func (u *transactionHandlers) GetAccountTransactionsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	accountIDStr := c.Param("accountID")
	limitStr := c.Query("limit")
	cursorStr := c.Query("cursor")

	var (
		limit int
		err   error
	)
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, domains.ErrorResp{
				Message: err.Error(),
			})
			return
		}
	}

	transactions, err := u.transactionRepository.GetTransactions(ctx, u.db, &repositories.GetTransactionsArgs{
		AccountID: accountIDStr,
		Cursor:    cursorStr,
		Limit:     limit,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	transactionsResp := make([]*domains.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		transactionsResp = append(transactionsResp, &domains.Transaction{
			TransactionID: transaction.TransactionID,
			UserID:        transaction.UserID,
			AccountID:     transaction.AccountID,
			Amount:        transaction.Amount,
			Balance:       transaction.Balance,
			Type:          transaction.Type,
			Status:        transaction.Status,
			Metadata:      transaction.Metadata,
			CreatedAt:     transaction.CreatedAt,
			UpdatedAt:     transaction.UpdatedAt,
		})
	}

	var nextCursor string
	if len(transactions) != 0 {
		nextCursor = transactions[len(transactions)-1].TransactionID
	}
	c.JSON(http.StatusOK, &domains.GetTransactionsResp{
		Transactions: transactionsResp,
		NextCursor:   nextCursor,
	})
}

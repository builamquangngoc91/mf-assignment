package handlers

import (
	"net/http"

	"banking-service/domains"
	"banking-service/repositories"

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
	DB *gorm.DB
}

type transactionHandlers struct {
	db                    *gorm.DB
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
	accountID := c.Param("accountID")

	transactions, err := u.transactionRepository.GetTransactions(ctx, u.db, &repositories.GetTransactionsArgs{
		AccountID: accountID,
	})
	if err != nil {
		c.Error(err)
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

	c.JSON(http.StatusOK, &domains.GetTransactionsResp{
		Transactions: transactionsResp,
	})
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"banking-service/configs"
	"banking-service/handlers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type (
	UserID  string
	GroupID string

	UserConnectionIdx int64
	ConnectionInfos   struct {
		UserID            UserID
		ExpiresAt         time.Time
		UserConnectionIdx UserConnectionIdx
	}

	MessageType int64

	MessageData struct {
		FromUserID UserID  `json:"from_user_id,omitempty"`
		ToUserID   UserID  `json:"to_user_id,omitempty"`
		ToGroupID  GroupID `json:"to_group_id,omitempty"`
		Text       string  `json:"text,omitempty"`
	}

	Message struct {
		Type MessageType  `json:"type,omitempty"`
		Data *MessageData `json:"data,omitempty"`
	}
)

const (
	PingMessage MessageType = iota + 1
	TextMessage
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("create logger error: %s", err.Error()))
	}

	configs.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		configs.Cfg.Database.Host,
		configs.Cfg.Database.Username,
		configs.Cfg.Database.Password,
		configs.Cfg.Database.Name,
		configs.Cfg.Database.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		logger.Sugar().Errorf("connect database error: %s", err.Error())
		return
	}

	// gracefull shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	router := gin.Default()

	// TODO: add ping and health
	accountHandlersDeps := &handlers.AccountHandlersDeps{
		DB: db,
	}
	accountHandlers := handlers.NewAccountHandlers(accountHandlersDeps)
	accountHandlers.RouteGroup(router)

	userHandlersDeps := &handlers.UserHandlersDeps{
		DB: db,
	}
	userHandlers := handlers.NewUserHandlers(userHandlersDeps)
	userHandlers.RouteGroup(router)

	transactionHandlersDeps := &handlers.TransactionHandlersDeps{
		DB: db,
	}
	transactionHandlers := handlers.NewTransactionHandlers(transactionHandlersDeps)
	transactionHandlers.RouteGroup(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", configs.Cfg.BankingService.Port),
		Handler: router,
	}
	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Sugar().Errorf("shutdown http.Server error: %s", err.Error())
			return
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Sugar().Errorf("ListenAndServe error: %s", err.Error())
		return
	}
}

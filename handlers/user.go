package handlers

import (
	"net/http"
	"time"

	"banking-service/domains"
	"banking-service/models"
	"banking-service/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	_ UserHandlers = &userHandlers{}
)

type UserHandlers interface {
	RouteGroup(r *gin.Engine)

	CreateUserHandler(*gin.Context)
	GetUsersHandler(*gin.Context)
	GetUserHandler(*gin.Context)
}

type UserHandlersDeps struct {
	DB *gorm.DB
}

type userHandlers struct {
	db                *gorm.DB
	userRepositiory   repositories.UserRepositoryI
	accountRepository repositories.AccountRepositoryI
}

func NewUserHandlers(deps *UserHandlersDeps) UserHandlers {
	if deps == nil {
		return nil
	}

	return &userHandlers{
		db:                deps.DB,
		userRepositiory:   repositories.NewUserRepository(),
		accountRepository: repositories.NewAccountRepository(),
	}
}

func (u *userHandlers) RouteGroup(rg *gin.Engine) {
	rg.POST("/users", u.CreateUserHandler)
	rg.GET("/users", u.GetUsersHandler)
	rg.GET("/users/:userID", u.GetUserHandler)
}

func (u *userHandlers) CreateUserHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req domains.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, domains.ErrorResp{
			Message: "missing name",
		})
	}

	user := &models.User{
		UserID:    uuid.NewString(),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := u.userRepositiory.Create(ctx, u.db, user); err != nil {
		c.JSON(http.StatusInternalServerError, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domains.CreateUserResponse{
		ID:        user.UserID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (u *userHandlers) GetUsersHandler(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := u.userRepositiory.GetUsers(ctx, u.db, &repositories.GetUsersArgs{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	usersResp := make([]*domains.User, 0, len(users))
	for _, user := range users {
		usersResp = append(usersResp, &domains.User{
			ID:        user.UserID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, domains.GetUsersResponse{
		Users: usersResp,
	})
}

func (u *userHandlers) GetUserHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("userID")

	user, err := u.userRepositiory.GetUser(ctx, u.db, &repositories.GetUserArgs{
		UserID: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	accountIDs, err := u.accountRepository.GetAccountIDs(ctx, u.db, &repositories.GetAccountIDsArgs{
		UserID: userID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, domains.ErrorResp{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domains.User{
		ID:         user.UserID,
		Name:       user.Name,
		AccountIDs: accountIDs,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	})
}

package repositories

import (
	"context"

	"banking-service/models"

	"gorm.io/gorm"
)

var _ UserRepositoryI = &userRepository{}

//go:generate mockery --name UserRepository
type UserRepositoryI interface {
	Create(ctx context.Context, db *gorm.DB, user *models.User) error
	GetUser(ctx context.Context, db *gorm.DB, args *GetUserArgs) (*models.User, error)
	GetUsers(ctx context.Context, db *gorm.DB, args *GetUsersArgs) (users []*models.User, _ error)
}

type userRepository struct {
}

func NewUserRepository() UserRepositoryI {
	return &userRepository{}
}

func (u *userRepository) Create(ctx context.Context, db *gorm.DB, user *models.User) error {
	return db.WithContext(ctx).Table("users").Create(user).Error
}

type GetUserArgs struct {
	UserID string
}

func (u *userRepository) GetUser(ctx context.Context, db *gorm.DB, args *GetUserArgs) (*models.User, error) {
	query := db.WithContext(ctx).Table("users")
	if args.UserID != "" {
		query.Where("user_id = ?", args.UserID)
	}

	var user models.User
	result := query.First(&user)

	return &user, result.Error
}

type GetUsersArgs struct {
	Cursor string
	Limit  int
}

func (u *userRepository) GetUsers(ctx context.Context, db *gorm.DB, args *GetUsersArgs) (users []*models.User, _ error) {
	db = db.WithContext(ctx).Table("users")
	if args.Cursor != "" {
		db.Where("user_id < ?", args.Cursor)
	}
	if args.Limit == 0 {
		args.Limit = 100
	}

	db.Order("user_id DESC")
	db.Limit(args.Limit)

	result := db.Find(&users)

	return users, result.Error
}

package user

import (
	"boilerplate-api/api/user/user"
	"boilerplate-api/lib/config"

	"gorm.io/gorm"
)

// UserRepository is the persistence contract for admin user operations.
// The service depends on this interface — not the concrete implementation —
// so unit tests can swap in a mock without a real database.
type UserRepository interface {
	WithTrx(trxHandle *gorm.DB) UserRepository
	Create(u user.CUser) error
	GetAllUsers(pagination Pagination) ([]GetUserResponse, int64, error)
	GetOneUser(id int64) (GetUserResponse, error)
	GetOneUserWithEmail(email string) (user.CUser, error)
	GetOneUserWithPhone(phone string) (user.CUser, error)
}

// gormRepository is the GORM-backed UserRepository implementation.
type gormRepository struct {
	db     config.Database
	logger config.Logger
}

// NewRepository creates a new CUser repository. Returns the interface to keep
// the fx-injected dependency loose at the type level.
func NewRepository(db config.Database, logger config.Logger) UserRepository {
	return gormRepository{db: db, logger: logger}
}

func (c gormRepository) WithTrx(trxHandle *gorm.DB) UserRepository {
	if trxHandle == nil {
		c.logger.Error("Transaction Database not found in gin context. ")
		return c
	}
	c.db.DB = trxHandle
	return c
}

func (c gormRepository) Create(u user.CUser) error {
	return c.db.DB.Create(&u).Error
}

func (c gormRepository) GetAllUsers(pagination Pagination) (users []GetUserResponse, count int64, err error) {
	queryBuilder := c.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset).Order("created_at desc")
	queryBuilder = queryBuilder.Model(&user.CUser{})

	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuilder.Where(c.db.DB.Where("`users`.`name` LIKE ?", searchQuery))
	}

	return users, count, queryBuilder.
		Find(&users).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error
}

func (c gormRepository) GetOneUser(id int64) (userModel GetUserResponse, err error) {
	return userModel, c.db.DB.
		Model(&userModel).
		Where("id = ?", id).
		First(&userModel).
		Error
}

func (c gormRepository) GetOneUserWithEmail(email string) (u user.CUser, err error) {
	return u, c.db.DB.Model(&u).
		Where("email = ?", email).
		First(&u).
		Error
}

func (c gormRepository) GetOneUserWithPhone(phone string) (u user.CUser, err error) {
	return u, c.db.DB.
		First(&u, "phone = ?", phone).
		Error
}

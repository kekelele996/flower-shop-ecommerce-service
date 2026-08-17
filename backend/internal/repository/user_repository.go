package repository

import (
	"errors"

	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
)

// UserRepository 用户仓储。
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return util.WrapAppError(50001, "create user failed: "+err.Error(), err)
	}
	return nil
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find user by username failed", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find user by id failed", err)
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return util.WrapAppError(50001, "update user failed", err)
	}
	return nil
}

func (r *UserRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	if err := r.db.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error; err != nil {
		return util.WrapAppError(50001, "update user fields failed", err)
	}
	return nil
}

func (r *UserRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, util.WrapAppError(50001, "count users failed", err)
	}
	return count, nil
}

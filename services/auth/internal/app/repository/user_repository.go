package repository

import (
	"github.com/kvvvseins/mictoservices/services/auth/internal/app/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindById найти по ID
func (ur *UserRepository) FindById(userId int) (*model.User, error) {
	var user model.User

	result := ur.db.First(&user, userId)

	return ur.responseFoundUser(&user, result.Error)
}

// FindByEmail найти по email
func (ur *UserRepository) FindByEmail(email string, onlyActive bool) (*model.User, error) {
	var user model.User

	var result *gorm.DB

	if onlyActive {
		result = ur.db.First(&user, "email = ? AND is_active = ?", email, true)
	} else {
		result = ur.db.First(&user, "email = ?", email)
	}

	return ur.responseFoundUser(&user, result.Error)
}

// FindByGuid найти по guid
func (ur *UserRepository) FindByGuid(guid string, onlyActive bool) (*model.User, error) {
	var user model.User

	var result *gorm.DB

	if onlyActive {
		result = ur.db.First(&user, "guid = ? AND is_active = ?", guid, true)
	} else {
		result = ur.db.First(&user, "guid = ?", guid)
	}

	return ur.responseFoundUser(&user, result.Error)
}

func (ur *UserRepository) responseFoundUser(user *model.User, err error) (*model.User, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения пользователя")
	}

	return user, nil
}

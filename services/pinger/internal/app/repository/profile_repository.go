package repository

import (
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/pinger/internal/app/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// FindByGuid найти по guid
func (ur *ProfileRepository) FindByGuid(guid uuid.UUID) (*model.Profile, error) {
	var user model.Profile

	var result *gorm.DB

	result = ur.db.First(&user, "guid = ?", guid.String())

	return ur.responseFoundProfile(&user, result.Error)
}

func (ur *ProfileRepository) responseFoundProfile(user *model.Profile, err error) (*model.Profile, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения профиля пользователя")
	}

	return user, nil
}

package repository

import (
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/billing/internal/app/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type BillingRepository struct {
	db *gorm.DB
}

func NewBillingRepository(db *gorm.DB) *BillingRepository {
	return &BillingRepository{db: db}
}

// FindByGuid найти по guid
func (ur *BillingRepository) FindByGuid(userID uuid.UUID) (*model.Billing, error) {
	var billing model.Billing

	var result *gorm.DB

	result = ur.db.First(&billing, "user_id = ?", userID.String())

	return ur.responseFoundProfile(&billing, result.Error)
}

func (ur *BillingRepository) responseFoundProfile(user *model.Billing, err error) (*model.Billing, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения профиля пользователя")
	}

	return user, nil
}

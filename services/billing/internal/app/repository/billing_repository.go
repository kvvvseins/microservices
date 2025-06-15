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

// FindByUserID найти по user id
func (ur *BillingRepository) FindByUserID(userID uuid.UUID) (*model.Billing, error) {
	var billing model.Billing

	var result *gorm.DB

	result = ur.db.First(&billing, "user_id = ?", userID.String())

	return ur.responseFound(&billing, result.Error)
}

func (ur *BillingRepository) responseFound(user *model.Billing, err error) (*model.Billing, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения счета пользователя")
	}

	return user, nil
}

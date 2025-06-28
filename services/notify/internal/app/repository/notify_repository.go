package repository

import (
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/notify/internal/app/model"
	"gorm.io/gorm"
)

type NotifyRepository struct {
	db *gorm.DB
}

func NewNotifyRepository(db *gorm.DB) *NotifyRepository {
	return &NotifyRepository{db: db}
}

// FindByUserID найти по user_id
func (ur *NotifyRepository) FindByUserID(userID uuid.UUID) ([]model.Notify, error) {
	var notifies []model.Notify
	if err := ur.db.Where("user_id = ?", userID.String()).Limit(10).Find(&notifies).Error; err != nil {
		return nil, err
	}

	return notifies, nil
}

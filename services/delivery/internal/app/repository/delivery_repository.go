package repository

import (
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/delivery/internal/app/model"
	"gorm.io/gorm"
)

type DeliveryRepository struct {
	db *gorm.DB
}

func NewDeliveryRepository(db *gorm.DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

// FindAll найти все
func (ur *DeliveryRepository) FindAll() ([]model.Delivery, error) {
	var delivery []model.Delivery
	if err := ur.db.
		Limit(10).
		Find(&delivery).Error; err != nil {
		return nil, err
	}

	return delivery, nil
}

// FindByOrderID найти по order_id
func (ur *DeliveryRepository) FindByOrderID(orderId uuid.UUID) (*model.Delivery, error) {
	var delivery model.Delivery
	if err := ur.db.
		Where("order_id = ? AND status = 2", orderId.String()).
		First(&delivery).Error; err != nil {
		return nil, err
	}

	return &delivery, nil
}

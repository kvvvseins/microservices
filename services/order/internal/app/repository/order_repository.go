package repository

import (
	"log"

	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/order/internal/app/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// FindByUserID найти по user_id
func (ur *OrderRepository) FindByUserID(userID uuid.UUID) ([]model.Order, error) {
	var orders []model.Order
	if err := ur.db.Where("user_id = ?", userID.String()).Limit(10).Find(&orders).Error; err != nil {
		log.Fatalf("error retrieving orders: %v", err)
	}

	return orders, nil
}

func (ur *OrderRepository) responseFound(order *model.Order, err error) (*model.Order, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения заказа пользователя")
	}

	return order, nil
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type Reserve struct {
	ID        uint `gorm:"primary_key;AUTO_INCREMENT;->"`
	StoreId   uint `gorm:"type:uuid;"`
	Quantity  int
	OrderId   *uuid.UUID `gorm:"type:uuid;<-:create;column:order_id"`
	UserID    uuid.UUID  `gorm:"type:uuid;<-:create;column:user_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Reserve) TableName() string {
	return "reserve"
}

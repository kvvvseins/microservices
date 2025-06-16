package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID        uint      `gorm:"primary_key;AUTO_INCREMENT;->"`
	Guid      uuid.UUID `gorm:"type:uuid;<-:create"`
	UserID    uuid.UUID `gorm:"type:uuid;column:user_id"`
	Price     uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate : hook before an order is created
func (u *Order) BeforeCreate(_ *gorm.DB) (err error) {
	if u.Guid == uuid.Nil {
		u.Guid = uuid.New()
	}

	return
}

func (*Order) TableName() string {
	return "orders"
}

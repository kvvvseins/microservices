package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Delivery struct {
	ID               uint      `gorm:"primary_key;AUTO_INCREMENT;->"`
	OrderId          uuid.UUID `gorm:"type:uuid;<-:create;column:order_id"`
	UserID           uuid.UUID `gorm:"type:uuid;<-:create;column:user_id"`
	Status           uint
	PlannedDateStart time.Time
	PlannedDateEnd   time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// BeforeCreate : hook before a delivery is created
func (u *Delivery) BeforeCreate(_ *gorm.DB) (err error) {
	if u.Status == 0 {
		u.Status = 2
	}

	return
}

func (*Delivery) TableName() string {
	return "delivery"
}

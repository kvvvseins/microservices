package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notify struct {
	ID        uint      `gorm:"primary_key;AUTO_INCREMENT;->"`
	UserID    uuid.UUID `gorm:"type:uuid;column:user_id"`
	Email     string
	Message   string
	Status    uint8
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Notify) TableName() string {
	return "notify"
}

// BeforeCreate : hook before notify is created
func (u *Notify) BeforeCreate(_ *gorm.DB) (err error) {
	if u.Status == 0 {
		u.Status = 1
	}

	return
}

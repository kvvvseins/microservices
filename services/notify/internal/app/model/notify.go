package model

import (
	"time"

	"github.com/google/uuid"
)

type Notify struct {
	ID        uint      `gorm:"primary_key;AUTO_INCREMENT;->"`
	UserID    uuid.UUID `gorm:"type:uuid;column:user_id"`
	Email     string
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Notify) TableName() string {
	return "notify"
}

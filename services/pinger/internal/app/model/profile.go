package model

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID        uint      `gorm:"primary_key;AUTO_INCREMENT;->"`
	Guid      uuid.UUID `gorm:"type:uuid;<-:create"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Profile) TableName() string {
	return "profile"
}

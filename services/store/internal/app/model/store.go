package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	ID        uint      `gorm:"primary_key;AUTO_INCREMENT;->"`
	Guid      uuid.UUID `gorm:"type:uuid;<-:create"`
	Name      string
	Price     uint
	Quantity  int `gorm:"->"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate : hook before model is created
func (u *Store) BeforeCreate(_ *gorm.DB) (err error) {
	if u.Guid == uuid.Nil {
		u.Guid = uuid.New()
	}

	return
}

func (*Store) TableName() string {
	return "store"
}

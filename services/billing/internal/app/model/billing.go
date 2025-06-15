package model

import (
	"time"

	"github.com/google/uuid"
)

type Billing struct {
	ID        uint      `gorm:"primary_key;AUTO_INCREMENT;->"`
	UserID    uuid.UUID `gorm:"type:uuid;column:user_id"`
	Value     uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Currency  string `gorm:"-:all"`
}

func (*Billing) TableName() string {
	//@todo тут Currency пустое — видимо этот метод вызывается у пустой структуры, или отказываться от gorm или нет
	return "rub"
}

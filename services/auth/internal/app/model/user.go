package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primary_key;AUTO_INCREMENT;->"`
	Guid      uuid.UUID `gorm:"type:uuid;<-:create"`
	Password  string
	Email     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*User) TableName() string {
	return "users"
}

// BeforeCreate : hook before a user is created
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	if u.Password != "" {
		var hash string

		hash, err = hashPassword(u.Password)
		if err != nil {
			return err
		}

		u.Password = hash
	} else {
		return errors.New("password is empty")
	}

	if u.Guid == uuid.Nil {
		u.Guid = uuid.New()
	}

	return
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

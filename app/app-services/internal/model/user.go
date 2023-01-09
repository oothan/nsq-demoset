package model

import (
	"context"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"nsq-demoset/app/app-services/cmd/front_api/criteria"
	"time"
)

type User struct {
	Id           uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserName     string         `gorm:"column:user_name;size:191;" json:"user_name"`
	Email        string         `gorm:"column:email;size:191;unique;not null;" json:"email"`
	Password     string         `gorm:"column:password;size:191;not null;" json:"password"`
	PasswordSalt string         `gorm:"column:password_salt;size:191;" json:"password_salt"`
	UserType     string         `gorm:"column:user_type" json:"user_type"`
	Activated    bool           `gorm:"column:activated;not null;default:0" json:"activated"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"_"`

	Conn map[string]*websocket.Conn `gorm:"-" json:"-"`
}

type UserService interface {
	FindAll(ctx context.Context, crits criteria.Criteria) ([]*User, error)
	FindById(ctx context.Context, userId uint64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, userId string) (*User, error)
}

type UserRepository interface {
	FindAll(ctx context.Context, crits criteria.Criteria) ([]*User, error)
	FindById(ctx context.Context, userId uint64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, userId string) (*User, error)
}

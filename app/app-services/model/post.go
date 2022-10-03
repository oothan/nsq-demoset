package model

import (
	"context"
	"gorm.io/gorm"
	"nsq-demoset/app/app-services/cmd/front_api/criteria"
	"time"
)

type Post struct {
	Id           uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title        string         `gorm:"column:title;size:191" json:"title"`
	Content      string         `gorm:"column:content;type:text" json:"content"`
	Status       string         `gorm:"column:status" json:"status"`
	UserId       uint64         `gorm:"column:user_id" json:"user_id"`
	ReadCount    int64          `gorm:"column:read_count;default:0" json:"read_count"`
	LikeCount    int64          `gorm:"column:like_count;default:0" json:"like_count"`
	CommentCount int64          `gorm:"column:comment_count;default:0" json:"comment_count"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"_"`

	User *User `gorm:"foreignKey:UserId;references:Id" json:"user"`
}

type PostService interface {
	FindAll(ctx context.Context, crits criteria.Criteria) ([]*Post, error)
	FindById(ctx context.Context, id uint64) (*Post, error)
	Create(ctx context.Context, post *Post) (*Post, error)
	Update(ctx context.Context, post *Post) (*Post, error)
	Delete(ctx context.Context, id uint64) (*Post, error)
	Count(ctx context.Context, crits criteria.Criteria) int64
}

type PostRepository interface {
	FindAll(ctx context.Context, crits criteria.Criteria) ([]*Post, error)
	FindById(ctx context.Context, id uint64) (*Post, error)
	Create(ctx context.Context, post *Post) (*Post, error)
	Update(ctx context.Context, post *Post) (*Post, error)
	Delete(ctx context.Context, id uint64) (*Post, error)
	Count(ctx context.Context, crits criteria.Criteria) int64
}

package model

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenDate struct {
	SS        string
	ID        uuid.UUID
	ExpiresIn time.Duration
}

type TokenService interface {
	GenerateTokenPair(ctx context.Context, user *User, prevToken string) (*TokenPair, error)
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, user *User, token *RefreshTokenDate) error
	DeleteRefreshToken(ctx context.Context, user *User, token string) error
}

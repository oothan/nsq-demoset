package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"nsq-demoset/app/app-services/internal/ds"
	"nsq-demoset/app/app-services/internal/model"
)

type tokenRepository struct {
	RDB *redis.Client
}

func (r *tokenRepository) StoreRefreshToken(ctx context.Context, user *model.User, token *model.RefreshTokenDate) error {
	key := fmt.Sprintf("%s:%v", token.ID, user.Id)
	if err := r.RDB.Set(ctx, key, 0, token.ExpiresIn).Err(); err != nil {
		return err
	}
	return nil
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, user *model.User, token string) error {
	key := fmt.Sprintf("%s:%v", token, user.Id)
	if err := r.RDB.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func NewTokenRepository(ds *ds.DataSource) model.TokenRepository {
	return &tokenRepository{
		RDB: ds.RDB,
	}
}

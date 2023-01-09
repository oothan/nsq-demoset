package service

import (
	"context"
	"nsq-demoset/app/app-services/conf"
	"nsq-demoset/app/app-services/internal/model"
	"nsq-demoset/app/app-services/internal/utils"
)

type tokenService struct {
	TokenRepo model.TokenRepository
}

func (t *tokenService) GenerateTokenPair(ctx context.Context, user *model.User, prevToken string) (*model.TokenPair, error) {
	if prevToken != "" {
		if err := t.TokenRepo.DeleteRefreshToken(ctx, user, prevToken); err != nil {
			return nil, err
		}
	}

	// generate access token
	accessToken, err := utils.GenerateAccessToken(user, conf.PrivateKey)
	if err != nil {
		return nil, err
	}

	refreshData, err := utils.GenerateRefreshToken(user.Id, conf.RefreshSecret)
	if err != nil {
		return nil, err
	}

	if err := t.TokenRepo.StoreRefreshToken(ctx, user, refreshData); err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshData.SS,
	}, nil
}

type TokenConfig struct {
	TokenRepo model.TokenRepository
}

func NewTokenService(c *TokenConfig) model.TokenService {
	return &tokenService{
		TokenRepo: c.TokenRepo,
	}
}

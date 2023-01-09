package service

import (
	"context"
	"nsq-demoset/app/app-services/cmd/front_api/criteria"
	"nsq-demoset/app/app-services/internal/model"
)

type userService struct {
	UserRepo model.UserRepository
}

func (u *userService) FindAll(ctx context.Context, crits criteria.Criteria) ([]*model.User, error) {
	return u.UserRepo.FindAll(ctx, crits)
}

func (u *userService) FindById(ctx context.Context, userId uint64) (*model.User, error) {
	return u.UserRepo.FindById(ctx, userId)
}

func (u *userService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return u.UserRepo.FindByEmail(ctx, email)
}

func (u *userService) Create(ctx context.Context, user *model.User) (*model.User, error) {
	return u.UserRepo.Create(ctx, user)
}

func (u *userService) Update(ctx context.Context, user *model.User) (*model.User, error) {
	return u.UserRepo.Update(ctx, user)
}

func (u *userService) Delete(ctx context.Context, userId string) (*model.User, error) {
	return u.UserRepo.Delete(ctx, userId)
}

type UserConfig struct {
	UserRepo model.UserRepository
}

func NewUserService(c *UserConfig) model.UserRepository {
	return &userService{
		UserRepo: c.UserRepo,
	}
}

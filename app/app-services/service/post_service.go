package service

import (
	"context"
	"nsq-demoset/app/app-services/cmd/front_api/criteria"
	"nsq-demoset/app/app-services/model"
)

type postService struct {
	PostRepo model.PostRepository
}

func (p *postService) FindAll(ctx context.Context, crits criteria.Criteria) ([]*model.Post, error) {
	return p.PostRepo.FindAll(ctx, crits)
}

func (p *postService) FindById(ctx context.Context, id uint64) (*model.Post, error) {
	return p.PostRepo.FindById(ctx, id)
}

func (p *postService) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	return p.PostRepo.Create(ctx, post)
}

func (p *postService) Update(ctx context.Context, post *model.Post) (*model.Post, error) {
	return p.PostRepo.Update(ctx, post)
}

func (p *postService) Delete(ctx context.Context, id uint64) (*model.Post, error) {
	return p.PostRepo.Delete(ctx, id)
}

func (p *postService) Count(ctx context.Context, crits criteria.Criteria) int64 {
	return p.PostRepo.Count(ctx, crits)
}

type PostConfig struct {
	PostRepo model.PostRepository
}

func NewPostService(c *PostConfig) model.PostService {
	return &postService{
		PostRepo: c.PostRepo,
	}
}

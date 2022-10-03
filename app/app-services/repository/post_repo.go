package repository

import (
	"context"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
	"nsq-demoset/app/app-services/cmd/front_api/criteria"
	"nsq-demoset/app/app-services/model"
	"nsq-demoset/app/nsq-services/ds"
)

type postRepository struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func (r *postRepository) FindAll(ctx context.Context, crits criteria.Criteria) ([]*model.Post, error) {
	posts := make([]*model.Post, 0)
	db, err := crits.Build(r.DB.WithContext(ctx).Model(&model.Post{}))
	if err != nil {
		return nil, err
	}
	if err := db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) FindById(ctx context.Context, id uint64) (*model.Post, error) {
	post := &model.Post{}
	db := r.DB.Model(&model.Post{})
	db = db.Where("id = ?", id)
	if err := db.First(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (r *postRepository) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	db := r.DB.Model(&model.Post{})
	if err := db.Create(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (r *postRepository) Update(ctx context.Context, post *model.Post) (*model.Post, error) {
	db := r.DB.Model(&model.Post{})
	db = db.Where("id = ?", post.Id)
	if err := db.Save(&post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (r *postRepository) Delete(ctx context.Context, id uint64) (*model.Post, error) {
	db := r.DB.Model(&model.Post{})
	db = db.Where("id = ?", id)
	if err := db.Delete(&model.Post{}).Error; err != nil {
		return nil, err
	}
	return &model.Post{
		Id: id,
	}, nil
}

func (r *postRepository) Count(ctx context.Context, crits criteria.Criteria) int64 {
	db, err := crits.Build(r.DB.Model(&model.Post{}).WithContext(ctx))
	if err != nil {
		return 0
	}
	total := int64(0)
	if err := db.Count(&total).Error; err != nil {
		return 0
	}
	return total
}

func NewPostRepository(ds *ds.DataSource) model.PostRepository {
	return &postRepository{
		DB:  ds.DB,
		RDB: ds.RDB,
	}
}

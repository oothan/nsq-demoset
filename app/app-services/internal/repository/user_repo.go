package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
	logger "nsq-demoset/app/_applib"
	"nsq-demoset/app/app-services/cmd/front_api/criteria"
	"nsq-demoset/app/app-services/internal/ds"
	"nsq-demoset/app/app-services/internal/model"
	"time"
)

type userRepository struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func (r *userRepository) FindAll(ctx context.Context, crits criteria.Criteria) ([]*model.User, error) {
	users := make([]*model.User, 0)
	db, err := crits.Build(r.DB.WithContext(ctx).Model(&model.User{}))
	if err != nil {
		return nil, err
	}
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) FindById(ctx context.Context, userId uint64) (*model.User, error) {
	user := &model.User{}

	key := fmt.Sprintf("cache:user:%v", userId)
	res, err := r.RDB.Get(ctx, key).Result()
	if err != nil || err == redis.Nil {
		db := r.DB.Model(&model.User{})
		if err := db.First(&user, userId).Error; err != nil {
			return nil, err
		}

		logger.Sugar.Debug(user.Activated, " model check activated")
		data, err := json.Marshal(&user)
		if err != nil {
			return nil, err
		}

		if err := r.RDB.Set(ctx, key, string(data), time.Hour*1).Err(); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal([]byte(res), &user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	key := fmt.Sprintf("cache:user:%v", email)
	res, err := r.RDB.Get(ctx, key).Result()
	if err != nil || err == redis.Nil {
		db := r.DB.Model(&model.User{})
		db = db.Where("email = ?", email)
		if err := db.First(&user).Error; err != nil {
			return nil, err
		}

		data, err := json.Marshal(&user)
		if err != nil {
			return nil, err
		}

		if err := r.RDB.Set(ctx, key, string(data), time.Hour*1).Err(); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal([]byte(res), &user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	keyId := fmt.Sprintf("cache:user:%v", user.Id)
	keyMail := fmt.Sprintf("cache:user:%v", user.Email)

	db := r.DB.Model(&model.User{})
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	if _, err := r.RDB.Del(ctx, keyId, keyMail).Result(); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	key := fmt.Sprintf("cache:user:%v", user.Id)

	db := r.DB.Model(&model.User{})
	db = db.Where("id = ?", user.Id)
	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	if _, err := r.RDB.Del(ctx, key).Result(); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, userId string) (*model.User, error) {
	key := fmt.Sprintf("cache:user:%v", userId)

	db := r.DB.Model(&model.User{})
	db = db.Where("id = ?", userId)
	if err := db.Delete(&model.User{}).Error; err != nil {
		return nil, err
	}

	if _, err := r.RDB.Del(ctx, key).Result(); err != nil {
		return nil, err
	}

	user := &model.User{}
	return user, nil
}

func NewUserRepository(ds *ds.DataSource) model.UserRepository {
	return &userRepository{
		DB:  ds.DB,
		RDB: ds.RDB,
	}
}

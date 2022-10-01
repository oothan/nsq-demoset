package ds

import (
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type DataSource struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewDataSource() {

}

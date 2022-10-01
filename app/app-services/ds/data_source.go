package ds

import (
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type DataSource struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewDataSource() *DataSource {
	db, err := LoadDB()
	if err != nil {
		return nil
	}

	rdb, err := LoadRDB()
	if err != nil {
		return nil
	}

	return &DataSource{
		DB:  db,
		RDB: rdb,
	}
}

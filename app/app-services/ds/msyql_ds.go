package ds

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger "nsq-demoset/app/_applib"
	"nsq-demoset/app/app-services/model"
	"os"
)

func LoadDB() (*gorm.DB, error) {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")
	name := os.Getenv("MYSQL_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	logger.Sugar.Info("Successfully connected to database")

	// migrate DB
	err = db.AutoMigrate(
		&model.User{},
		&model.Post{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

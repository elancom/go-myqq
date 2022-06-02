package model

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var NotFound = logger.ErrRecordNotFound

func One[T any](db *gorm.DB) (*T, error) {
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, db.Error
	}
	return db.Statement.Dest.(*T), nil
}

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:123456@/qq?parseTime=true&loc=Local")
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}

func NewGDb(db *sql.DB, conf ...*gorm.Config) (*gorm.DB, error) {
	config := gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	if len(conf) > 0 {
		config.Logger = conf[0].Logger
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &config)
	if err != nil {
		return nil, err
	}
	return gdb, nil
}

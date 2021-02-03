package database

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github/demo/config"
	"github/demo/database/dialects"
	"github/demo/utils/log"
)

type IDatabase interface {
	SetPool(int, int, time.Duration) bool
	GetDB() *gorm.DB
	IsConnected() bool
	Close() bool
}

type database struct {
	gormDB *gorm.DB
}

func (db *database) SetPool(idleConns int, openConns int, duration time.Duration) bool {
	sqlDB := db.gormDB.DB()
	sqlDB.SetMaxIdleConns(idleConns)
	sqlDB.SetMaxOpenConns(openConns)
	sqlDB.SetConnMaxLifetime(duration)
	return true
}

func (db *database) GetDB() *gorm.DB {
	return db.gormDB
}

func (db *database) IsConnected() bool {
	if db.gormDB == nil {
		return false
	}
	return true
}

func (db *database) Close() bool {
	db.gormDB.Close()
	return true
}

func NewDatabase(c *config.Database) (IDatabase, error) {
	db := &database{}
	switch dialects.Dialect(c.Dialect) {
	case dialects.Postgres:
		db.gormDB = dialects.PostgreSQL(c)
	case dialects.Sqlite:
		db.gormDB = dialects.SqliteDB(c)
	default:
		return nil, fmt.Errorf("Database not support: %q", c.Dialect)
	}

	log.Info("Create database success")
	return db, nil
}

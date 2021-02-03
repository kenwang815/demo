package dialects

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"

	"github/demo/config"
	"github/demo/utils/log"
)

func SqliteDB(c *config.Database) *gorm.DB {
	db, err := gorm.Open("sqlite3", c.Host)
	if err != nil {
		log.Errorf("gorm open fail => %+v", err)
		return nil
	}
	db.LogMode(true)

	return db
}

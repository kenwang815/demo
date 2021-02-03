package dialects

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github/demo/config"
	"github/demo/utils/log"
)

func PostgreSQL(c *config.Database) *gorm.DB {
	connect := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", c.Host, c.Port, c.User, c.Name, c.Password)

	pDB, err := gorm.Open("postgres", connect)
	if err != nil {
		log.Errorf("failed to connect database, %+v", err)
		return nil
	}

	return pDB
}

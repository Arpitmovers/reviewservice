package db

import (
	"fmt"
	"log"

	"github.com/Arpitmovers/reviewservice/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDBConnect(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DbUser, cfg.DbPwd, cfg.DbHost, cfg.DbPort, cfg.DbName)
	fmt.Println("db string", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect database: %v", err)
		return nil
	}

	log.Println(cfg.DbName, " connected to database successfully")

	return db
}

package db

import (
	"fmt"

	"github.com/Arpitmovers/reviewservice/internal/config"
	"go.uber.org/zap"

	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDBConnect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DbUser, cfg.DbPwd, cfg.DbHost, cfg.DbPort, cfg.DbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Logger.Error("error in  db connect", zap.Error(err))

		return nil, err
	}

	logger.Logger.Info("connected to database successfully", zap.String("dbName", cfg.DbName))

	return db, err
}

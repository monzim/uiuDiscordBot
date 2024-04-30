package db

import (
	zLog "github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabaseConnection(DB_URI string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DB_URI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		return nil, err
	}

	zLog.Info().Msg("Database connection established")
	return db, nil

}

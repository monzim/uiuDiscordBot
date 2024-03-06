package db

import (
	"github.com/monzim/uiuBot/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Info().Msg("Migrating the database...")
	db.Migrator().DropTable()

	err := db.AutoMigrate(
		&models.Exam{},
	)

	if err != nil {
		return err
	}

	log.Info().Msg("Database migration completed")

	// setup the database
	log.Info().Msg("Setting up the database with bulk data...")
	err = setup(db)
	if err != nil {
		return err
	}

	log.Info().Msg("Database setup completed")
	return nil
}

func setup(_ *gorm.DB) error {

	return nil
}

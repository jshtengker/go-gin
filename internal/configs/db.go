package configs

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnPostgres(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		return nil, fmt.Errorf("unnable to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil{
		return nil, fmt.Errorf("unable to get the underlying sql.DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil{
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return db, nil
}
package database

import (
	"fmt"

	"building-report-backend/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
	// &entity.User{},
	// &entity.Report{},
	// &entity.ReportPhoto{},
	// &entity.SpatialPlanningReport{},
	// &entity.SpatialPlanningPhoto{},
	// &entity.WaterResourcesReport{},
	// &entity.WaterResourcesPhoto{},
	// &entity.BinaMargaReport{},
	// &entity.BinaMargaPhoto{},
	// &entity.AgricultureReport{},
	// &entity.AgriculturePhoto{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

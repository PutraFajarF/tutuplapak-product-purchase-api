package postgresql

import (
	"fmt"
	"log"

	"time"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/config"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	loggers "gorm.io/gorm/logger"
)

func New(cfg *config.Config, l *logger.Logger) *gorm.DB {
	// Format DSN Postgres
	// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Jakarta"
	dsn := cfg.POSTGRESQL.URL

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: loggers.Default.LogMode(loggers.Info),
	})
	if err != nil {
		log.Fatalf(fmt.Sprintf("couldn't connect to database: %v", err))
	}

	// Set connection pool (pakai DB dari gorm.DB)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf(fmt.Sprintf("couldn't get sql.DB from gorm: %v", err))
	}

	sqlDB.SetMaxIdleConns(cfg.POSTGRESQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.POSTGRESQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.POSTGRESQL.MaxLifetimeConns) * time.Second)

	seedCategories(db)

	return db
}

func seedCategories(db *gorm.DB) {
	query := `
        INSERT INTO categories (type) VALUES
            ('Food'),
            ('Beverage'),
            ('Clothes'),
            ('Furniture'),
            ('Tools')
        ON CONFLICT (type) DO NOTHING;
    `
	if err := db.Exec(query).Error; err != nil {
		log.Printf("failed to seed categories: %v", err)
	}

	log.Println("Categories seeded successfully (if not already present).")
}

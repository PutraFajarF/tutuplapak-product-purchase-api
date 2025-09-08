package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/config"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/httpserver"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/logger"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/postgresql"

	"github.com/gorilla/mux"
)

func Run(cfg *config.Config) {
	fmt.Println("Running Service TutupLapak Product Purchase API")

	var err error
	l := logger.New(cfg)

	// Postgresql
	db := postgresql.New(cfg, l)
	// / Untuk tutup koneksi, ambil sql.DB dari gorm
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("couldn't get sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Repository
	// consumerRepository := postgresql_repository.NewConsumerMysqlRepository(l, cfg, db)

	// Usecase
	// consumerUsecase := consumer.NewConsumerUsecase(l, cfg, consumerRepository)

	// HTTP Server
	handler := mux.NewRouter()
	// v1.NewRouter(handler, l, cfg, consumerUsecase)
	httpServer := httpserver.New(handler, cfg, httpserver.Port(cfg.HTTPServer.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

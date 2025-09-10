package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/config"
	purchaseHandler "github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/purchase/delivery"
	purchaseRepo "github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/purchase/repository"
	purchaseUsecase "github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/purchase/usecase"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/httpserver"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/logger"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/postgresql"
	validator "github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// Repositories
	purchRepo := purchaseRepo.NewPurchaseRepository(db)

	// Usecases
	purchUc := purchaseUsecase.NewPurchaseUsecase(purchRepo)

	// Delivery
	purchHandler := purchaseHandler.NewPurchaseDelivery(purchUc)

	// Echo HTTP Server
	e := echo.New()
	e.HideBanner = true
	e.Validator = validator.NewValidator()
	e.Use(middleware.Recover(), middleware.Logger())
	e.GET("/healthz", func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	// Routers (v1)
	api := e.Group("/api/v1")
	purchaseHandler.RegisterPurchaseRoutes(api, purchHandler)

	// v1.NewRouter(handler, l, cfg, consumerUsecase)
	httpServer := httpserver.New(e, cfg, httpserver.Port(cfg.HTTPServer.Port))

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

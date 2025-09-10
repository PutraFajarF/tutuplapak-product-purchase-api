package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/config"
	file_repo "github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/file/repository"
	product_handler "github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/product/delivery"
	product_repo "github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/product/repository"
	product_usecase "github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/product/usecase"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/logger"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/postgresql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(cfg *config.Config) {
	fmt.Println("Running Service TutupLapak Product Purchase API")

	var err error
	l := logger.New(cfg)

	// Postgresql
	db := postgresql.New(cfg, l)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("couldn't get sql.DB: %v", err)
	}
	defer sqlDB.Close()
	// Echo
	e := echo.New()
	e.Validator = NewCustomValidator()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/api/v1")

	// Repository
	productRepository := product_repo.NewProductRepository(*db)
	fileRepository := file_repo.NewRepositoryFile(*db)

	// Usecase
	productUsecase := product_usecase.NewProductUsecase(productRepository, fileRepository)

	// Routes
	product_handler.RegisterProductRoutes(v1, productUsecase)

	// Start server in goroutine
	go func() {
		if err := e.Start(":" + cfg.HTTPServer.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}

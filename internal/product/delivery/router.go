package delivery

import (
	"time"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/config"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/product"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterProductRoutes(e *echo.Group, u product.IUsecaseProduct, cfg *config.Config) {
	handler := NewProductHandler(u)

	api := e.Group("/product")
	api.POST("", handler.CreateProduct, middleware.JWTMiddlewareHS256(cfg, 5*time.Second))
	api.GET("", handler.ProductList)
	api.PUT("/:productId", handler.UpdateProduct, middleware.JWTMiddlewareHS256(cfg, 5*time.Second))
	api.DELETE("/:productId", handler.DeleteProduct, middleware.JWTMiddlewareHS256(cfg, 5*time.Second))
}

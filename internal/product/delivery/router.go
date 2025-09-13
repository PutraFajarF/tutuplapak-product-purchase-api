package delivery

import (
	"time"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/config"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterProductRoutes(e *echo.Group, u ProductHandler, cfg *config.Config) {

	api := e.Group("/product")
	api.POST("", u.CreateProduct, middleware.JWTMiddlewareHS256(cfg, 5*time.Second))
	api.GET("", u.ProductList)
	api.PUT("/:productId", u.UpdateProduct, middleware.JWTMiddlewareHS256(cfg, 5*time.Second))
	api.DELETE("/:productId", u.DeleteProduct, middleware.JWTMiddlewareHS256(cfg, 5*time.Second))
}

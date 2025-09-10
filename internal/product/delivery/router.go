package delivery

import (
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/product"
	"github.com/labstack/echo/v4"
)

func RegisterProductRoutes(e *echo.Group, u product.IUsecaseProduct) {
	handler := NewProductHandler(u)

	api := e.Group("/product")
	api.POST("", handler.CreateProduct)
}

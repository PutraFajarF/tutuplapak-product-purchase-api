package delivery

import (
	"net/http"
	"strings"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/product"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	usecaseProduct product.IUsecaseProduct
}

func NewProductHandler(usecaseProduct product.IUsecaseProduct) ProductHandler {
	return ProductHandler{usecaseProduct: usecaseProduct}
}

func (d ProductHandler) CreateProduct(c echo.Context) error {
	var req product.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	res, err := d.usecaseProduct.CreateProduct(c.Request().Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "unique_auth_sku") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Product with this SKU already exists for your account"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create product"})
	}

	return c.JSON(http.StatusCreated, res)
}

func (d ProductHandler) UpdateProduct(c echo.Context) error {
	var req product.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	req.ProductId = c.Param("productId")
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	res, err := d.usecaseProduct.UpdateProduct(c.Request().Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "unique_auth_sku") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Product with this SKU already exists for your account"})
		}
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	var req product.DeleteProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	req.ProductId = c.Param("productId")
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	err := h.usecaseProduct.DeleteProduct(c.Request().Context(), req.ProductId)
	if err != nil {
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "product not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete product",
		})
	}

	return c.NoContent(http.StatusOK)
}

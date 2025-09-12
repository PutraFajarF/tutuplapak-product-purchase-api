package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/product"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/middleware"
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

	authId, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(401, "no user id in token")
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	req.AuthId = authId
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

	authId, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(401, "no user id in token")
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	req.AuthId = authId
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

	authId, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(401, "no user id in token")
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	req.ProductId = c.Param("productId")
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	err := h.usecaseProduct.DeleteProduct(c.Request().Context(), authId, req.ProductId)
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

func (d ProductHandler) ProductList(c echo.Context) error {
	var req product.ProductListRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	req = product.NewProductListRequest(req)
	data, _ := json.Marshal(req)

	log.Println("DATA REQ", string(data))

	res, err := d.usecaseProduct.GetProducts(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get product list"})
	}

	return c.JSON(http.StatusOK, res)
}

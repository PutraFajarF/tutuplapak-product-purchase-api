package delivery

import (
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
		return c.JSON(http.StatusBadRequest, map[string]any{"message": "Invalid request"})
	}

	req.AuthId = authId
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}

	res, err := d.usecaseProduct.CreateProduct(c.Request().Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "unique_auth_sku") {
			return c.JSON(http.StatusConflict, map[string]any{"message": "Product with this SKU already exists for your account"})
		}
		if strings.Contains(err.Error(), "file not found") {
			return c.JSON(http.StatusBadRequest, map[string]any{"message": "fileId does not exist"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{"message": "Failed to create product"})
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
		return c.JSON(http.StatusBadRequest, map[string]any{"message": "Invalid request"})
	}

	req.AuthId = authId
	req.ProductId = c.Param("productId")
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}

	res, err := d.usecaseProduct.UpdateProduct(c.Request().Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "unique_auth_sku") {
			return c.JSON(http.StatusConflict, map[string]any{"message": "Product with this SKU already exists for your account"})
		}
		if strings.Contains(err.Error(), "product not found") || strings.Contains(err.Error(), "invalid productId") {
			return c.JSON(http.StatusNotFound, map[string]any{"message": "Product not found"})
		}
		if strings.Contains(err.Error(), "file not found") {
			return c.JSON(http.StatusBadRequest, map[string]any{"message": "fileId does not exist"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{"message": "Failed to update product"})
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
		return c.JSON(http.StatusBadRequest, map[string]any{"message": "Invalid request"})
	}

	req.ProductId = c.Param("productId")
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}

	err := h.usecaseProduct.DeleteProduct(c.Request().Context(), authId, req.ProductId)
	if err != nil {
		if strings.Contains(err.Error(), "product not found") || strings.Contains(err.Error(), "invalid productId") {
			return c.JSON(http.StatusNotFound, map[string]any{
				"message": "Product not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": "Failed to delete product",
		})
	}

	return c.NoContent(http.StatusOK)
}

func (d ProductHandler) ProductList(c echo.Context) error {
	var req product.ProductListRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"message": "Invalid request"})
	}

	req = product.NewProductListRequest(req)

	res, err := d.usecaseProduct.GetProducts(c.Request().Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid cursor") {
			return c.JSON(http.StatusBadRequest, map[string]any{"message": "Invalid cursor"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{"message": "Failed to get product list"})
	}

	return c.JSON(http.StatusOK, res)
}

func (d ProductHandler) GetProduct(c echo.Context) error {
	productId := c.Param("productId")

	res, err := d.usecaseProduct.GetProductByID(c.Request().Context(), productId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "invalid productId") {
			return c.JSON(http.StatusNotFound, map[string]any{"message": "Product not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{"message": "Failed to get product"})
	}

	return c.JSON(http.StatusOK, res)
}

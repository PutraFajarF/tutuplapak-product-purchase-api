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

// CreateProduct godoc
//
//	@Summary		Create a new product
//	@Description	Create a new product with the provided details
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			request	body		product.CreateProductRequest	true	"Product creation request"
//	@Success		201		{object}	product.CreateProdutResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		Bearer
//	@Router			/api/v1/product [post]
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

// UpdateProduct godoc
//
//	@Summary		Update an existing product
//	@Description	Update a product with the provided details
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			productId	path		string						true	"Product ID"
//	@Param			request		body		product.UpdateProductRequest	true	"Product update request"
//	@Success		200			{object}	product.UpdateProdutResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		409			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		Bearer
//	@Router			/api/v1/product/{productId} [put]
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

// DeleteProduct godoc
//
//	@Summary		Delete a product
//	@Description	Delete a product by its ID
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			productId	path	string	true	"Product ID"
//	@Success		200
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		Bearer
//	@Router			/api/v1/product/{productId} [delete]
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

// ProductList godoc
//
//	@Summary		Get list of products
//	@Description	Get a paginated list of products with optional filters
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			limit		query		int		false	"Limit"		default(5)
//	@Param			offset		query		int		false	"Offset"	default(0)
//	@Param			productId	query		string	false	"Product ID"
//	@Param			sku			query		string	false	"SKU"
//	@Param			category	query		string	false	"Category"
//	@Param			sortBy		query		string	false	"Sort by"	Enums(newest,oldest,cheapest,expensive)
//	@Success		200			{array}		product.ProdutListResponse
//	@Failure		500			{object}	map[string]string
//	@Router			/api/v1/product [get]
func (d ProductHandler) ProductList(c echo.Context) error {
	var req product.ProductListRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	req = product.NewProductListRequest(req)

	res, err := d.usecaseProduct.GetProducts(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get product list"})
	}

	return c.JSON(http.StatusOK, res)
}

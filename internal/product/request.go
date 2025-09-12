package product

import (
	"strconv"
	"strings"
)

type CreateProductRequest struct {
	Name     string `json:"name" validate:"required,min=4,max=32"`
	Category string `json:"category" validate:"required,oneof=Food Beverage Clothes Furniture Tools"`
	Qty      int    `json:"qty" validate:"required,gte=1"`
	Price    int    `json:"price" validate:"required,gte=100"`
	SKU      string `json:"sku" validate:"required,max=32"`
	FileID   string `json:"fileId" validate:"required"`
}
type UpdateProductRequest struct {
	ProductId string `json:"productId" validate:"required"`
	Name      string `json:"name" validate:"required,min=4,max=32"`
	Category  string `json:"category" validate:"required,oneof=Food Beverage Clothes Furniture Tools"`
	Qty       int    `json:"qty" validate:"required,gte=1"`
	Price     int    `json:"price" validate:"required,gte=100"`
	SKU       string `json:"sku" validate:"required,max=32"`
	FileID    string `json:"fileId" validate:"required"`
}

type ProductListRequest struct {
	Limit        int    `query:"limit"`
	Offset       int    `query:"offset"`
	ProductId    string `query:"productId"`
	Sku          string `query:"sku"`
	Category     string `query:"category"`
	SortBy       string `query:"sortBy"`
	ProductIdInt int    `query:"-"`
}

func NewProductListRequest(req ProductListRequest) (res ProductListRequest) {
	res = req

	if res.Limit <= 0 {
		res.Limit = 5
	}
	if res.Offset < 0 {
		res.Offset = 0
	}

	productIdInt, err := strconv.Atoi(res.ProductId)
	if err == nil {
		res.ProductIdInt = productIdInt
	} else {
		res.ProductIdInt = 0
	}

	sort := strings.ToLower(res.SortBy)
	if sort != "newest" && sort != "oldest" &&
		sort != "cheapest" && sort != "expensive" {
		res.SortBy = ""
	}

	cat := res.Category
	if cat != "Food" && cat != "Beverage" &&
		cat != "Clothes" && cat != "Furniture" &&
		cat != "Tools" {
		res.Category = ""
	}

	return
}

type DeleteProductRequest struct {
	ProductId string `json:"productId" validate:"required"`
}

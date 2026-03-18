package product

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CreateProductRequest struct {
	Name     string `json:"name" validate:"required,min=4,max=32"`
	Category string `json:"category" validate:"required,oneof=Food Beverage Clothes Furniture Tools"`
	Qty      int    `json:"qty" validate:"required,gte=1"`
	Price    int    `json:"price" validate:"required,gte=100"`
	SKU      string `json:"sku" validate:"required,max=32"`
	FileID   string `json:"fileId" validate:"required"`
	AuthId   string `json:"-"`
}
type UpdateProductRequest struct {
	ProductId string `json:"productId" validate:"required"`
	Name      string `json:"name" validate:"required,min=4,max=32"`
	Category  string `json:"category" validate:"required,oneof=Food Beverage Clothes Furniture Tools"`
	Qty       int    `json:"qty" validate:"required,gte=1"`
	Price     int    `json:"price" validate:"required,gte=100"`
	SKU       string `json:"sku" validate:"required,max=32"`
	FileID    string `json:"fileId" validate:"required"`
	AuthId    string `json:"-"`
}

type ProductListRequest struct {
	Limit        int    `query:"limit"`
	Cursor       string `query:"cursor"`
	ProductId    string `query:"productId"`
	Sku          string `query:"sku"`
	Name         string `query:"name"`
	Category     string `query:"category"`
	SortBy       string `query:"sortBy"`
	ProductIdInt int    `query:"-"`
}

// CursorData holds the decoded cursor position for keyset pagination.
// ID is always the tiebreaker. CreatedAt/Price are set based on sortBy.
type CursorData struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"ca,omitempty"`
	Price     int       `json:"p,omitempty"`
}

func EncodeCursor(c CursorData) string {
	b, _ := json.Marshal(c)
	return base64.URLEncoding.EncodeToString(b)
}

func DecodeCursor(encoded string) (CursorData, error) {
	var c CursorData
	b, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return c, fmt.Errorf("invalid cursor")
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return c, fmt.Errorf("invalid cursor")
	}
	if c.ID <= 0 {
		return c, fmt.Errorf("invalid cursor: missing id")
	}
	return c, nil
}

func NewProductListRequest(req ProductListRequest) (res ProductListRequest) {
	res = req

	if res.Limit <= 0 {
		res.Limit = 5
	}
	if res.Limit > 100 {
		res.Limit = 100
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

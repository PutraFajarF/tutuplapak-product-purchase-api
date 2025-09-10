package product

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
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	ProductId string `json:"productId"`
	Sku       string `json:"sku"`
	Category  string `json:"category"`
	SortBy    string `json:"sortBy"`
}

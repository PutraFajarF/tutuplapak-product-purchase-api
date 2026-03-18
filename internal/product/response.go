package product

type ProductResponse struct {
	ProductId        string `json:"productId"`
	Category         string `json:"category"`
	Name             string `json:"name"`
	Qty              int    `json:"qty"`
	Price            int    `json:"price"`
	SKU              string `json:"sku"`
	FileID           string `json:"fileId"`
	FileUri          string `json:"fileUri"`
	FileThumbnailUri string `json:"fileThumbnailUri"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

// Aliases for backward compatibility
type CreateProdutResponse = ProductResponse
type UpdateProdutResponse = ProductResponse
type ProdutListResponse = ProductResponse

type ProductListPaginatedResponse struct {
	Data   []ProductResponse `json:"data"`
	Total  int64             `json:"total"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

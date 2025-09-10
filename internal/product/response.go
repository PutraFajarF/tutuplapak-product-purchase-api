package product

type CreateProdutResponse struct {
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

type UpdateProdutResponse struct {
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

type ProdutListResponse struct {
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

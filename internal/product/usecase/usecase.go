package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/file"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/product"
	"gorm.io/gorm"
)

type ProductUsecase struct {
	productRepository product.IRepositoryProduct
	fileRepository    file.IRepositoryFile
}

func NewProductUsecase(productRepository product.IRepositoryProduct, fileRepository file.IRepositoryFile) ProductUsecase {
	return ProductUsecase{productRepository: productRepository, fileRepository: fileRepository}
}

func (u ProductUsecase) CreateProduct(ctx context.Context, req product.CreateProductRequest) (res product.CreateProdutResponse, err error) {
	// Validate file exists before creating product
	fileData, err := u.fileRepository.GetFileByID(ctx, req.FileID)
	if err != nil {
		return res, fmt.Errorf("file not found: %s", req.FileID)
	}

	createProduct := entity.Product{
		AuthID:       req.AuthId,
		TypeCategory: req.Category,
		Name:         req.Name,
		Qty:          req.Qty,
		Price:        req.Price,
		Sku:          req.SKU,
		FileId:       req.FileID,
	}

	savedProduct, err := u.productRepository.CreateProduct(ctx, createProduct)
	if err != nil {
		return res, err
	}

	res.ProductId = intToString(savedProduct.ID)
	res.Category = savedProduct.TypeCategory
	res.Name = savedProduct.Name
	res.Qty = savedProduct.Qty
	res.Price = savedProduct.Price
	res.SKU = savedProduct.Sku
	res.FileID = req.FileID
	res.FileUri = fileData.FileUri
	res.FileThumbnailUri = fileData.FileThumbnailUri
	res.CreatedAt = savedProduct.CreatedAt.Format(time.RFC3339)
	res.UpdatedAt = savedProduct.UpdatedAt.Format(time.RFC3339)

	return
}

func (u ProductUsecase) UpdateProduct(ctx context.Context, req product.UpdateProductRequest) (res product.UpdateProdutResponse, err error) {
	prdId, err := stringToInt(req.ProductId)
	if err != nil {
		return res, fmt.Errorf("invalid productId: %s", req.ProductId)
	}

	// Validate file exists before updating product
	fileData, err := u.fileRepository.GetFileByID(ctx, req.FileID)
	if err != nil {
		return res, fmt.Errorf("file not found: %s", req.FileID)
	}

	updateProductReq := entity.Product{
		ID:           prdId,
		AuthID:       req.AuthId,
		TypeCategory: req.Category,
		Name:         req.Name,
		Qty:          req.Qty,
		Price:        req.Price,
		Sku:          req.SKU,
		FileId:       req.FileID,
	}

	updateProduct, err := u.productRepository.UpdateProduct(ctx, updateProductReq)
	if err != nil {
		return res, err
	}

	res.ProductId = intToString(updateProduct.ID)
	res.Category = updateProduct.TypeCategory
	res.Name = updateProduct.Name
	res.Qty = updateProduct.Qty
	res.Price = updateProduct.Price
	res.SKU = updateProduct.Sku
	res.FileID = updateProduct.FileId
	res.FileUri = fileData.FileUri
	res.FileThumbnailUri = fileData.FileThumbnailUri
	res.CreatedAt = updateProduct.CreatedAt.Format(time.RFC3339)
	res.UpdatedAt = updateProduct.UpdatedAt.Format(time.RFC3339)

	return
}

func (u ProductUsecase) DeleteProduct(ctx context.Context, authId, productId string) (err error) {
	prdId, err := stringToInt(productId)
	if err != nil {
		return fmt.Errorf("invalid productId: %s", productId)
	}
	err = u.productRepository.DeleteProduct(ctx, authId, prdId)
	if err != nil {
		return err
	}
	return nil
}

func (u ProductUsecase) GetProducts(ctx context.Context, req product.ProductListRequest) (res product.ProductListPaginatedResponse, err error) {
	res = product.ProductListPaginatedResponse{
		Data:  []product.ProductResponse{},
		Limit: req.Limit,
	}

	// Decode cursor if provided
	var cursor *product.CursorData
	if req.Cursor != "" {
		decoded, err := product.DecodeCursor(req.Cursor)
		if err != nil {
			return res, fmt.Errorf("invalid cursor")
		}
		cursor = &decoded
	}

	products, err := u.productRepository.GetProducts(ctx, req, cursor)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res, nil
		}
		return res, err
	}

	// For backward ("prev") direction the DB returns items in reversed order,
	// so we reverse them back to the natural display order.
	if req.Direction == "prev" {
		for i, j := 0, len(products)-1; i < j; i, j = i+1, j-1 {
			products[i], products[j] = products[j], products[i]
		}
	}

	for _, p := range products {
		prd := product.ProductResponse{
			ProductId:        intToString(p.ID),
			Category:         p.TypeCategory,
			Name:             p.Name,
			Qty:              p.Qty,
			Price:            p.Price,
			SKU:              p.Sku,
			FileID:           p.File.ID,
			FileUri:          p.File.FileUri,
			FileThumbnailUri: p.File.FileThumbnailUri,
			CreatedAt:        p.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        p.UpdatedAt.Format(time.RFC3339),
		}
		res.Data = append(res.Data, prd)
	}

	if len(products) == 0 {
		return res, nil
	}

	first := products[0]
	last := products[len(products)-1]

	firstCursor := product.EncodeCursor(product.CursorData{
		ID:        first.ID,
		CreatedAt: first.CreatedAt,
		Price:     first.Price,
	})
	lastCursor := product.EncodeCursor(product.CursorData{
		ID:        last.ID,
		CreatedAt: last.CreatedAt,
		Price:     last.Price,
	})

	switch req.Direction {
	case "prev":
		// Going backward: there's always a next page (we came from there).
		res.NextCursor = lastCursor
		// There's a previous page only if we got a full page of results.
		if len(products) == req.Limit {
			res.PrevCursor = firstCursor
		}
	default: // "next" or first page
		// There's a next page only if we got a full page of results.
		if len(products) == req.Limit {
			res.NextCursor = lastCursor
		}
		// There's a previous page only if a cursor was provided (not first page).
		if cursor != nil {
			res.PrevCursor = firstCursor
		}
	}

	return
}

func (u ProductUsecase) GetProductByID(ctx context.Context, productId string) (res product.ProductResponse, err error) {
	prdId, err := stringToInt(productId)
	if err != nil {
		return res, fmt.Errorf("invalid productId: %s", productId)
	}

	p, err := u.productRepository.GetProductByID(ctx, prdId)
	if err != nil {
		return res, err
	}

	res.ProductId = intToString(p.ID)
	res.Category = p.TypeCategory
	res.Name = p.Name
	res.Qty = p.Qty
	res.Price = p.Price
	res.SKU = p.Sku
	res.FileID = p.File.ID
	res.FileUri = p.File.FileUri
	res.FileThumbnailUri = p.File.FileThumbnailUri
	res.CreatedAt = p.CreatedAt.Format(time.RFC3339)
	res.UpdatedAt = p.UpdatedAt.Format(time.RFC3339)

	return
}

func intToString(n int) string {
	return fmt.Sprintf("%d", n)
}

func stringToInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

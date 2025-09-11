package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/file"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/product"
)

type ProductUsecase struct {
	productRepository product.IRepositoryProduct
	fileRepository    file.IRepositoryFile
}

func NewProductUsecase(productRepository product.IRepositoryProduct, fileRepository file.IRepositoryFile) ProductUsecase {
	return ProductUsecase{productRepository: productRepository, fileRepository: fileRepository}
}

func (u ProductUsecase) CreateProduct(ctx context.Context, req product.CreateProductRequest) (res product.CreateProdutResponse, err error) {
	// get auth id from jwt token
	authId := "authid"

	createProduct := entity.Product{
		AuthId:       authId,
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

	fileData, err := u.fileRepository.GetFileByID(ctx, req.FileID)
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
	res.FileUri = fileData.Uri
	res.FileThumbnailUri = fileData.ThumbnailUri
	res.CreatedAt = savedProduct.CreatedAt.Format(time.RFC3339)
	res.UpdatedAt = savedProduct.UpdatedAt.Format(time.RFC3339)

	return
}

func (u ProductUsecase) UpdateProduct(ctx context.Context, req product.UpdateProductRequest) (res product.UpdateProdutResponse, err error) {
	// get auth id from jwt token
	authId := "authid"

	prdId, _ := stringToInt(req.ProductId)
	updateProductReq := entity.Product{
		ID:           prdId,
		AuthId:       authId,
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

	fileData, err := u.fileRepository.GetFileByID(ctx, req.FileID)
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
	res.FileUri = fileData.Uri
	res.FileThumbnailUri = fileData.ThumbnailUri
	res.CreatedAt = updateProduct.CreatedAt.Format(time.RFC3339)
	res.UpdatedAt = updateProduct.UpdatedAt.Format(time.RFC3339)

	return
}

func (u ProductUsecase) DeleteProduct(ctx context.Context, productId string) (err error) {
	// ambil auth id dari token
	authId := "authid"

	prdId, _ := stringToInt(productId)
	err = u.productRepository.DeleteProduct(ctx, authId, prdId)
	if err != nil {
		return err
	}
	return nil
}

func (u ProductUsecase) GetProducts(ctx context.Context, req product.ProductListRequest) (res []product.ProdutListResponse, err error) {
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

package repository

import (
	"context"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/product"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db gorm.DB
}

func NewProductRepository(db gorm.DB) ProductRepository {
	return ProductRepository{db: db}
}

func (r ProductRepository) CreateProduct(ctx context.Context, req entity.Product) (entity.Product, error) {
	err := r.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return entity.Product{}, err
	}
	return req, nil
}

func (r ProductRepository) UpdateProduct(ctx context.Context, req product.UpdateProductRequest) (res product.UpdateProdutResponse, err error) {
	return
}

func (r ProductRepository) DeleteProduct(ctx context.Context, authId string, productId string) (err error) {
	return
}

func (r ProductRepository) GetProducts(ctx context.Context, req product.ProductListRequest) (res []product.ProdutListResponse, err error) {
	return
}

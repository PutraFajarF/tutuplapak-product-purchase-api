package repository

import (
	"context"
	"fmt"
	"time"

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

func (r ProductRepository) UpdateProduct(ctx context.Context, req entity.Product) (entity.Product, error) {
	result := r.db.WithContext(ctx).Model(&entity.Product{}).Where("auth_id = ? AND id = ?", req.AuthID, req.ID).Updates(&req)
	if result.Error != nil {
		return entity.Product{}, result.Error
	}
	if result.RowsAffected == 0 {
		return entity.Product{}, fmt.Errorf("product not found")
	}

	var res entity.Product
	if err := r.db.WithContext(ctx).
		Where("id = ? AND auth_id = ?", req.ID, req.AuthID).
		First(&res).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Product{}, fmt.Errorf("product not found")
		}
		return entity.Product{}, err
	}

	return res, nil
}

func (r ProductRepository) DeleteProduct(ctx context.Context, authId string, productId int) error {
	result := r.db.WithContext(ctx).
		Model(&entity.Product{}).
		Where("auth_id = ? AND id = ? AND deleted_at IS NULL", authId, productId).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r ProductRepository) GetProducts(ctx context.Context, req product.ProductListRequest) (res []entity.Product, err error) {
	query := r.db.Model(&entity.Product{}).Preload("File")

	if req.ProductIdInt != 0 {
		query = query.Where("id = ?", req.ProductIdInt)
	}

	if req.Sku != "" {
		query = query.Where("sku = ?", req.Sku)
	}

	if req.Category != "" {
		query = query.Where("type = ?", req.Category)
	}

	switch req.SortBy {
	case "newest":
		query = query.Order("created_at DESC").Order("updated_at DESC")
	case "oldest":
		query = query.Order("created_at ASC").Order("updated_at ASC")
	case "cheapest":
		query = query.Order("price ASC")
	case "expensive":
		query = query.Order("price DESC")
	}

	err = query.Limit(req.Limit).Offset(req.Offset).Find(&res).Error
	if err != nil {
		return res, err
	}

	return
}

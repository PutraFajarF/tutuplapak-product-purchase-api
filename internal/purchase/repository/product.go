package repository

import (
	"context"
	"errors"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
)

func (r *purchaseRepository) GetProductByIDs(ctx context.Context, ids []int64) ([]entity.Product, error) {
	var products []entity.Product
	if len(ids) == 0 {
		return products, nil
	}
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("File").
		Where("id IN ?", ids).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, errors.New("no products found")
	}
	return products, nil
}

func (r *purchaseRepository) DecrementProductQty(ctx context.Context, productID int64, by int) error {
	res := r.db.WithContext(ctx).Exec("UPDATE products SET qty = qty - ?, updated_at = NOW() WHERE id = ?", by, productID)
	return res.Error
}

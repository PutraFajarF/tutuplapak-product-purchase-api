package repository

import (
	"context"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
)

func (r *purchaseRepository) ExistsAll(ctx context.Context, ids []int64) (bool, error) {
	if len(ids) == 0 {
		return false, nil
	}
	var cnt int64
	if err := r.db.WithContext(ctx).Model(&entity.File{}).Where("id IN ?", ids).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt == int64(len(ids)), nil
}

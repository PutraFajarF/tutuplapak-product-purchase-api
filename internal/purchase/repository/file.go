package repository

import (
	"context"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
)

func (r *purchaseRepository) ExistsAllByFileID(ctx context.Context, fileIDs []string) (bool, error) {
	if len(fileIDs) == 0 {
		return false, nil
	}
	var cnt int64
	if err := r.db.WithContext(ctx).Model(&entity.File{}).Where("\"fileId\" IN ?", fileIDs).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt == int64(len(fileIDs)), nil
}

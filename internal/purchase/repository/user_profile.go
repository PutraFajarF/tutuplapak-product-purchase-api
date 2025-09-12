package repository

import (
	"context"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
)

func (r *purchaseRepository) GetProfilesByIDs(ctx context.Context, ids []string) (map[string]entity.Profile, error) {
	m := make(map[string]entity.Profile)
	if len(ids) == 0 {
		return m, nil
	}
	var profiles []entity.Profile
	if err := r.db.WithContext(ctx).Where("user_id IN ?", ids).Find(&profiles).Error; err != nil {
		return nil, err
	}
	for _, p := range profiles {
		m[p.UserID] = p
	}
	return m, nil
}

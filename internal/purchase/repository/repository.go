package repository

import (
	"context"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/purchase"
	"gorm.io/gorm"
)

type purchaseRepository struct {
	db gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) purchase.IPurchaseRepository {
	return &purchaseRepository{db: *db}
}

func (r *purchaseRepository) CreatePurchaseWithItems(ctx context.Context, p *entity.Purchase, items []entity.PurchaseItem) error {
	return nil
}

func (r *purchaseRepository) GetPurchaseByIDWithItems(ctx context.Context, id int64) (*entity.Purchase, []entity.PurchaseItem, error) {
	return nil, nil, nil
}

func (r *purchaseRepository) SetPurchasePaidWithProofsAndDecrement(ctx context.Context, id int64, proofs []entity.PurchasePaymentProof) error {
	return nil
}

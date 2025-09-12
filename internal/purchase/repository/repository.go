package repository

import (
	"context"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/purchase"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type purchaseRepository struct {
	db gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) purchase.IPurchaseRepository {
	return &purchaseRepository{db: *db}
}

func (r *purchaseRepository) CreatePurchaseWithItems(ctx context.Context, p *entity.Purchase, items []entity.PurchaseItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].PurchaseID = p.ID
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *purchaseRepository) GetPurchaseByIDWithItems(ctx context.Context, id int64) (entity.Purchase, []entity.PurchaseItem, error) {
	var p entity.Purchase
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return entity.Purchase{}, nil, err
	}
	var items []entity.PurchaseItem
	if err := r.db.WithContext(ctx).Where("purchase_id = ?", id).Find(&items).Error; err != nil {
		return entity.Purchase{}, nil, err
	}
	return p, items, nil
}

func (r *purchaseRepository) SetPurchasePaidWithProofsAndDecrement(ctx context.Context, id int64, proofs []entity.PurchasePaymentProof) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var p entity.Purchase
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&p, "id = ?", id).Error; err != nil {
			return err
		}
		if p.Status == "PAID" {
			return gorm.ErrInvalidData
		}

		if err := tx.Create(&proofs).Error; err != nil {
			return err
		}
		var items []entity.PurchaseItem
		if err := tx.Where("purchase_id = ?", id).Find(&items).Error; err != nil {
			return err
		}
		for _, it := range items {
			if err := tx.Exec("UPDATE products SET qty = qty - ?, updated_at = NOW() WHERE id = ?", it.BuyQty, it.ProductID).Error; err != nil {
				return err
			}
		}
		if err := tx.Exec("UPDATE purchases SET status = 'PAID', updated_at = NOW() WHERE id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

package purchase

import (
	"context"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
)

type IPurchaseUsecase interface {
	CreatePurchase(ctx context.Context, req CreatePurchaseReq) (*CreatePurchaseResp, error)
	UploadPurchaseProofs(ctx context.Context, purchaseID int64, req UploadProofReq) error
}

type IPurchaseRepository interface {
	CreatePurchaseWithItems(ctx context.Context, p *entity.Purchase, items []entity.PurchaseItem) error
	GetPurchaseByIDWithItems(ctx context.Context, id int64) (*entity.Purchase, []entity.PurchaseItem, error)
	SetPurchasePaidWithProofsAndDecrement(ctx context.Context, id int64, proofs []entity.PurchasePaymentProof) error
	IProductRepository
	IUserProfileRepository
	IFileRepository
}

type IProductRepository interface {
	GetProductByIDs(ctx context.Context, ids []int64) ([]entity.Product, error)
	DecrementProductQty(ctx context.Context, productID int64, by int) error
}

type IUserProfileRepository interface {
	GetProfilesByIDs(ctx context.Context, ids []string) (map[string]*entity.UserProfile, error)
}

type IFileRepository interface {
	ExistsAll(ctx context.Context, ids []int64) (bool, error)
}

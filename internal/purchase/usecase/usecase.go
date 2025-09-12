package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/entity"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/purchase"
	"github.com/PutraFajarF/tutuplapak-product-purchase-api/pkg/helper"
)

type purchaseUsecase struct {
	purchaseRepo purchase.IPurchaseRepository
}

func NewPurchaseUsecase(purchaseRepo purchase.IPurchaseRepository) purchase.IPurchaseUsecase {
	return &purchaseUsecase{purchaseRepo: purchaseRepo}
}

func (uc *purchaseUsecase) CreatePurchase(ctx context.Context, req purchase.CreatePurchaseReq) (*purchase.CreatePurchaseResp, error) {
	// kumpulkan product ids & qty
	buyQty := map[int64]int{}
	order := make([]int64, 0, len(req.PurchasedItems))
	for _, it := range req.PurchasedItems {
		pid, err := strconv.ParseInt(it.ProductID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid productId: %s", it.ProductID)
		}
		buyQty[pid] += it.Qty
		order = append(order, pid)
	}
	uniqSet := map[int64]struct{}{}
	uniq := make([]int64, 0)
	for _, id := range order {
		if _, ok := uniqSet[id]; !ok {
			uniqSet[id] = struct{}{}
			uniq = append(uniq, id)
		}
	}

	products, err := uc.purchaseRepo.GetProductByIDs(ctx, uniq)
	if err != nil {
		return nil, err
	}

	// build snapshot items (no pointers)
	items := make([]entity.PurchaseItem, 0, len(products))
	var total int64
	sellerSet := map[string]struct{}{}

	for _, p := range products {
		bq := buyQty[p.ID]
		if p.Qty < bq {
			return nil, fmt.Errorf("product %d qty not enough (have %d, want %d)", p.ID, p.Qty, bq)
		}

		it := entity.PurchaseItem{
			ProductID:        p.ID,
			SellerID:         p.AuthID,
			Name:             p.Name,
			CategoryCode:     p.TypeCategory,
			Price:            p.Price,
			SKU:              p.Sku,
			FileID:           sql.NullString{String: p.FileId, Valid: p.FileId != ""},
			FileURI:          sql.NullString{String: p.File.FileUri, Valid: p.File.FileUri != ""},
			FileThumbnailURI: sql.NullString{String: p.File.FileThumbnailUri, Valid: p.File.FileThumbnailUri != ""},
			ProductCreatedAt: p.CreatedAt,
			ProductUpdatedAt: p.UpdatedAt,
			QtyBefore:        p.Qty,
			BuyQty:           bq,
		}
		items = append(items, it)
		total += int64(bq) * p.Price
		sellerSet[p.AuthID] = struct{}{}
	}

	purch := &entity.Purchase{
		SenderName:          req.SenderName,
		SenderContactType:   req.SenderContactType,
		SenderContactDetail: req.SenderContactDetail,
		TotalPrice:          total,
		Status:              "PENDING",
	}
	if err := uc.purchaseRepo.CreatePurchaseWithItems(ctx, purch, items); err != nil {
		return nil, err
	}

	// response: convert ke bentuk non-pointer & empty string
	respItems := make([]purchase.ProductSnapshotResp, 0, len(items))
	for _, it := range items {
		respItems = append(respItems, purchase.ProductSnapshotResp{
			ProductID:        strconv.FormatInt(it.ProductID, 10),
			Name:             it.Name,
			Category:         it.CategoryCode,
			Qty:              it.QtyBefore,
			Price:            it.Price,
			SKU:              it.SKU,
			FileID:           helper.NullString(it.FileID),
			FileURI:          helper.NullString(it.FileURI),
			FileThumbnailURI: helper.NullString(it.FileThumbnailURI),
			CreatedAt:        it.ProductCreatedAt,
			UpdatedAt:        it.ProductUpdatedAt,
		})
	}

	// payment details per seller
	sellerIDs := make([]string, 0, len(sellerSet))
	for id := range sellerSet {
		sellerIDs = append(sellerIDs, id)
	}
	profiles, err := uc.purchaseRepo.GetProfilesByIDs(ctx, sellerIDs)
	if err != nil {
		return nil, err
	}

	totalPerSeller := map[string]int64{}
	for _, it := range items {
		totalPerSeller[it.SellerID] += int64(it.BuyQty) * it.Price
	}

	sort.Strings(sellerIDs)
	pay := make([]purchase.PaymentDetailResp, 0, len(sellerIDs))
	for _, sid := range sellerIDs {
		p := profiles[sid]
		pay = append(pay, purchase.PaymentDetailResp{
			BankAccountName:   helper.NullString(p.BankAccountName),
			BankAccountHolder: helper.NullString(p.BankAccountHolder),
			BankAccountNumber: helper.NullString(p.BankAccountNumber),
			TotalPrice:        totalPerSeller[sid],
		})
	}

	return &purchase.CreatePurchaseResp{
		PurchaseID:     strconv.FormatInt(purch.ID, 10),
		PurchasedItems: respItems,
		TotalPrice:     total,
		PaymentDetails: pay,
	}, nil
}

func (uc *purchaseUsecase) UploadPurchaseProofs(ctx context.Context, purchaseID int64, req purchase.UploadProofReq) error {
	// get purchase + items by value
	p, items, err := uc.purchaseRepo.GetPurchaseByIDWithItems(ctx, purchaseID)
	if err != nil {
		return err
	}
	if p.Status == "PAID" {
		return errors.New("purchase already paid")
	}

	// stable sellers order
	sellers := map[string]struct{}{}
	ordered := []string{}
	for _, it := range items {
		if _, ok := sellers[it.SellerID]; !ok {
			sellers[it.SellerID] = struct{}{}
			ordered = append(ordered, it.SellerID)
		}
	}
	if len(req.FileIDs) != len(ordered) {
		return fmt.Errorf("fileIds must equal number of sellers: %d", len(ordered))
	}

	// validate fileIds (files.fileId)
	exists, err := uc.purchaseRepo.ExistsAllByFileID(ctx, req.FileIDs)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("some fileIds do not exist")
	}

	// build proofs
	proofs := make([]entity.PurchasePaymentProof, 0, len(ordered))
	for i, sid := range ordered {
		proofs = append(proofs, entity.PurchasePaymentProof{
			PurchaseID: p.ID,
			SellerID:   sid,
			FileID:     req.FileIDs[i],
			CreatedAt:  time.Now(),
		})
	}

	return uc.purchaseRepo.SetPurchasePaidWithProofsAndDecrement(ctx, p.ID, proofs)
}

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
	quantityByProductID := make(map[int64]int)
	inputOrderIDs := make([]int64, 0, len(req.PurchasedItems))
	for _, purchasedItem := range req.PurchasedItems {
		productID, err := strconv.ParseInt(purchasedItem.ProductID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid productId: %s", purchasedItem.ProductID)
		}
		quantityByProductID[productID] += purchasedItem.Qty
		inputOrderIDs = append(inputOrderIDs, productID)
	}
	uniqueIDsSet := map[int64]struct{}{}
	uniqueIDs := make([]int64, 0)
	for _, productID := range inputOrderIDs {
		if _, ok := uniqueIDsSet[productID]; !ok {
			uniqueIDsSet[productID] = struct{}{}
			uniqueIDs = append(uniqueIDs, productID)
		}
	}

	products, err := uc.purchaseRepo.GetProductByIDs(ctx, uniqueIDs)
	if err != nil {
		return nil, err
	}

	// Build snapshot items
	snapshotItems := make([]entity.PurchaseItem, 0, len(products))
	var grandTotal int64
	uniqueSellerIDs := map[string]struct{}{}

	for _, product := range products {
		requestedQty := quantityByProductID[int64(product.ID)]
		if product.Qty < requestedQty {
			return nil, fmt.Errorf("product %d qty not enough (have %d, want %d)", product.ID, product.Qty, requestedQty)
		}

		snapshotItem := entity.PurchaseItem{
			ProductID:        int64(product.ID),
			SellerID:         product.AuthID,
			Name:             product.Name,
			CategoryCode:     product.TypeCategory,
			Price:            int64(product.ID),
			SKU:              product.Sku,
			FileID:           sql.NullString{String: product.FileId, Valid: product.FileId != ""},
			FileURI:          sql.NullString{String: product.File.FileUri, Valid: product.File.FileUri != ""},
			FileThumbnailURI: sql.NullString{String: product.File.FileThumbnailUri, Valid: product.File.FileThumbnailUri != ""},
			ProductCreatedAt: product.CreatedAt,
			ProductUpdatedAt: product.UpdatedAt,
			QtyBefore:        product.Qty,
			BuyQty:           requestedQty,
		}
		snapshotItems = append(snapshotItems, snapshotItem)
		grandTotal += int64(requestedQty) * int64(product.Price)
		uniqueSellerIDs[product.AuthID] = struct{}{}
	}

	purchaseRecord := &entity.Purchase{
		SenderName:          req.SenderName,
		SenderContactType:   req.SenderContactType,
		SenderContactDetail: req.SenderContactDetail,
		TotalPrice:          grandTotal,
		Status:              "PENDING",
	}
	if err := uc.purchaseRepo.CreatePurchaseWithItems(ctx, purchaseRecord, snapshotItems); err != nil {
		return nil, err
	}

	// Build response
	responseItems := make([]purchase.ProductSnapshotResp, 0, len(snapshotItems))
	for _, snapshotItem := range snapshotItems {
		responseItems = append(responseItems, purchase.ProductSnapshotResp{
			ProductID:        strconv.FormatInt(snapshotItem.ProductID, 10),
			Name:             snapshotItem.Name,
			Category:         snapshotItem.CategoryCode,
			Qty:              snapshotItem.QtyBefore,
			Price:            snapshotItem.Price,
			SKU:              snapshotItem.SKU,
			FileID:           helper.NullString(snapshotItem.FileID),
			FileURI:          helper.NullString(snapshotItem.FileURI),
			FileThumbnailURI: helper.NullString(snapshotItem.FileThumbnailURI),
			CreatedAt:        snapshotItem.ProductCreatedAt,
			UpdatedAt:        snapshotItem.ProductUpdatedAt,
		})
	}

	// Payment details per seller
	sellerIDs := make([]string, 0, len(uniqueSellerIDs))
	for sellerID := range uniqueSellerIDs {
		sellerIDs = append(sellerIDs, sellerID)
	}
	profilesBySeller, err := uc.purchaseRepo.GetProfilesByIDs(ctx, sellerIDs)
	if err != nil {
		return nil, err
	}

	totalPerSeller := map[string]int64{}
	for _, snapshotItem := range snapshotItems {
		totalPerSeller[snapshotItem.SellerID] += int64(snapshotItem.BuyQty) * snapshotItem.Price
	}

	sort.Strings(sellerIDs)
	paymentDetails := make([]purchase.PaymentDetailResp, 0, len(sellerIDs))
	for _, sellerID := range sellerIDs {
		profile := profilesBySeller[sellerID]
		paymentDetails = append(paymentDetails, purchase.PaymentDetailResp{
			BankAccountName:   helper.NullString(profile.BankAccountName),
			BankAccountHolder: helper.NullString(profile.BankAccountHolder),
			BankAccountNumber: helper.NullString(profile.BankAccountNumber),
			TotalPrice:        totalPerSeller[sellerID],
		})
	}

	return &purchase.CreatePurchaseResp{
		PurchaseID:     strconv.FormatInt(purchaseRecord.ID, 10),
		PurchasedItems: responseItems,
		TotalPrice:     grandTotal,
		PaymentDetails: paymentDetails,
	}, nil
}

func (uc *purchaseUsecase) UploadPurchaseProofs(ctx context.Context, purchaseID int64, req purchase.UploadProofReq) error {
	// get purchase + items by value
	purchaseRecord, purchaseItems, err := uc.purchaseRepo.GetPurchaseByIDWithItems(ctx, purchaseID)
	if err != nil {
		return err
	}
	if purchaseRecord.Status == "PAID" {
		return errors.New("purchase already paid")
	}

	// Build ordered list of unique sellers
	uniqueSellerIDsSet := map[string]struct{}{}
	orderedSellerIDs := []string{}
	for _, item := range purchaseItems {
		if _, exists := uniqueSellerIDsSet[item.SellerID]; !exists {
			uniqueSellerIDsSet[item.SellerID] = struct{}{}
			orderedSellerIDs = append(orderedSellerIDs, item.SellerID)
		}
	}
	if len(req.FileIDs) != len(orderedSellerIDs) {
		return fmt.Errorf("fileIds must equal number of sellers: %d", len(orderedSellerIDs))
	}

	// Validate files
	exists, err := uc.purchaseRepo.ExistsAllByFileID(ctx, req.FileIDs)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("some fileIds do not exist")
	}

	// Build proofs
	paymentProofs := make([]entity.PurchasePaymentProof, 0, len(orderedSellerIDs))
	for i, sellerID := range orderedSellerIDs {
		paymentProofs = append(paymentProofs, entity.PurchasePaymentProof{
			PurchaseID: purchaseRecord.ID,
			SellerID:   sellerID,
			FileID:     req.FileIDs[i],
			CreatedAt:  time.Now(),
		})
	}

	return uc.purchaseRepo.SetPurchasePaidWithProofsAndDecrement(ctx, purchaseRecord.ID, paymentProofs)
}

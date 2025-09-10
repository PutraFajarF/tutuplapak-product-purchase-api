package purchase

import "time"

type ProductSnapshotResp struct {
	ProductID        string     `json:"productId"`
	Name             string     `json:"name"`
	Category         string     `json:"category"`
	Qty              int        `json:"qty"` // qty before bought
	Price            int64      `json:"price"`
	SKU              string     `json:"sku"`
	FileID           *int64     `json:"fileId"`
	FileURI          *string    `json:"fileUri"`
	FileThumbnailURI *string    `json:"fileThumbnailUri"`
	CreatedAt        *time.Time `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}

type PaymentDetailResp struct {
	BankAccountName   *string `json:"bankAccountName"`
	BankAccountHolder *string `json:"bankAccountHolder"`
	BankAccountNumber *string `json:"bankAccountNumber"`
	TotalPrice        int64   `json:"totalPrice"`
}

type CreatePurchaseResp struct {
	PurchaseID     string                `json:"purchaseId"`
	PurchasedItems []ProductSnapshotResp `json:"purchasedItems"`
	TotalPrice     int64                 `json:"totalPrice"`
	PaymentDetails []PaymentDetailResp   `json:"paymentDetails"`
}

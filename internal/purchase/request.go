package purchase

type CreatePurchaseItemReq struct {
	ProductID string `json:"productId" validate:"required,number"`
	Qty       int    `json:"qty" validate:"required,gte=2"`
}

type CreatePurchaseReq struct {
	PurchasedItems      []CreatePurchaseItemReq `json:"purchasedItems" validate:"required,min=1,dive"`
	SenderName          string                  `json:"senderName" validate:"required,min=4,max=55"`
	SenderContactType   string                  `json:"senderContactType" validate:"required,oneof=email phone"`
	SenderContactDetail string                  `json:"senderContactDetail" validate:"required"`
}

type UploadProofReq struct {
	FileIDs []string `json:"fileIds" validate:"required,min=1,dive,printascii"`
}

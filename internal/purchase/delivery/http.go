package delivery

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/PutraFajarF/tutuplapak-product-purchase-api/internal/purchase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PurchaseDelivery struct {
	uc purchase.IPurchaseUsecase
}

func NewPurchaseDelivery(uc purchase.IPurchaseUsecase) *PurchaseDelivery {
	return &PurchaseDelivery{uc: uc}
}

// Create godoc
//
//	@Summary		Create a purchase
//	@Description	Create a new purchase with items and sender details
//	@Tags			purchases
//	@Accept			json
//	@Produce		json
//	@Param			request	body		purchase.CreatePurchaseReq	true	"Purchase creation request"
//	@Success		201		{object}	purchase.CreatePurchaseResp
//	@Failure		400		{object}	map[string]any
//	@Failure		404		{object}	map[string]any
//	@Failure		500		{object}	map[string]any
//	@Router			/api/v1/purchase [post]
func (h *PurchaseDelivery) Create(c echo.Context) error {
	var req purchase.CreatePurchaseReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}

	// conditional validation for senderContactDetail
	switch req.SenderContactType {
	case "phone":
		if err := c.Validate(&struct {
			Phone string `json:"senderContactDetail" validate:"e164"`
		}{Phone: req.SenderContactDetail}); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		}
	case "email":
		if err := c.Validate(&struct {
			Email string `json:"senderContactDetail" validate:"email"`
		}{Email: req.SenderContactDetail}); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		}
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}

	resp, err := h.uc.CreatePurchase(c.Request().Context(), req)
	if err != nil {
		// simple mapping without helper
		switch {
		case errors.Is(err, echo.ErrNotFound):
			return c.JSON(http.StatusNotFound, map[string]any{"message": "not found"})
		case errors.Is(err, gorm.ErrInvalidData), strings.Contains(err.Error(), "qty not enough"):
			return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]any{"message": err.Error()})
		}
	}
	return c.JSON(http.StatusCreated, resp)
}

// UploadProof godoc
//
//	@Summary		Upload payment proof for a purchase
//	@Description	Upload proof of payment files for an existing purchase
//	@Tags			purchases
//	@Accept			json
//	@Produce		json
//	@Param			purchaseId	path		string					true	"Purchase ID"
//	@Param			request		body		purchase.UploadProofReq	true	"Upload proof request"
//	@Success		201
//	@Failure		400	{object}	map[string]any
//	@Failure		404	{object}	map[string]any
//	@Failure		500	{object}	map[string]any
//	@Router			/api/v1/purchase/{purchaseId} [post]
func (h *PurchaseDelivery) UploadProof(c echo.Context) error {
	sid := c.Param("purchaseId")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}

	var req purchase.UploadProofReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}

	if err := h.uc.UploadPurchaseProofs(c.Request().Context(), id, req); err != nil {
		switch {
		case errors.Is(err, echo.ErrNotFound):
			return c.JSON(http.StatusNotFound, map[string]any{"message": "not found"})
		case errors.Is(err, gorm.ErrInvalidData),
			strings.Contains(err.Error(), "already paid"),
			strings.Contains(err.Error(), "fileIds must equal"),
			strings.Contains(err.Error(), "do not exist"):
			return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]any{"message": err.Error()})
		}
	}
	return c.NoContent(http.StatusCreated)
}

package delivery

import (
	"github.com/labstack/echo/v4"
)

// RegisterPurchaseRoutes mendaftarkan semua endpoint Purchase di prefix /v1
func RegisterPurchaseRoutes(g *echo.Group, h *PurchaseDelivery) {
	g.POST("/purchase", h.Create)
	g.GET("/purchase/:purchaseId", h.GetPurchase)
	g.POST("/purchase/:purchaseId", h.UploadProof)
}

package handler

import (
	"net/http"
	
	"github.com/google/uuid"

	"github.com/datnguyen/life_capital/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type PriceSyncHandler struct {
	syncService service.PriceSyncService
}

func NewPriceSyncHandler(syncService service.PriceSyncService) *PriceSyncHandler {
	return &PriceSyncHandler{
		syncService: syncService,
	}
}

// @Summary Sync asset prices
// @Description Fetch current prices for assets and update them in the database
// @Tags PriceSync
// @Accept json
// @Produce json
// @Success 200 {object} service.SyncReport
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/price-sync [post]
func (h *PriceSyncHandler) SyncPrices(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing user ID")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	report, err := h.syncService.SyncByUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, report)
}

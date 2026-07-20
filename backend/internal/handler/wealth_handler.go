package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/service"
)

type WealthHandler struct {
	wealthService service.WealthService
}

func NewWealthHandler(wealthService service.WealthService) *WealthHandler {
	return &WealthHandler{wealthService: wealthService}
}

// GetNetWorthSummary
// @Summary Get net worth summary
// @Description Get total assets, total liabilities and net worth for the current user
// @Tags Wealth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} model.NetWorthSummary
// @Router /api/v1/wealth/net-worth [get]
func (h *WealthHandler) GetNetWorthSummary(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	summary, err := h.wealthService.GetNetWorthSummary(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, summary)
}

// CreateAsset
// @Summary Create a new asset
// @Tags Wealth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateAssetRequest true "Asset info"
// @Success 201 {object} model.Asset
// @Router /api/v1/wealth/assets [post]
func (h *WealthHandler) CreateAsset(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	req := new(model.CreateAssetRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	asset, err := h.wealthService.CreateAsset(c.Request().Context(), userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, asset)
}

// GetAssets
// @Summary Get all assets for current user
// @Tags Wealth
// @Produce json
// @Security Bearer
// @Param category query string false "Category filter"
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} model.PaginatedAssets
// @Router /api/v1/wealth/assets [get]
func (h *WealthHandler) GetAssets(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	category := c.QueryParam("category")
	sort := c.QueryParam("sort")
	limit := 0
	offset := 0
	echo.QueryParamsBinder(c).Int("limit", &limit).Int("offset", &offset)

	assets, err := h.wealthService.GetAssets(c.Request().Context(), userID, category, sort, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, assets)
}

// UpdateAsset
// @Summary Update an asset
// @Tags Wealth
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Asset ID"
// @Param request body model.UpdateAssetRequest true "Asset info"
// @Success 200 {object} model.Asset
// @Router /api/v1/wealth/assets/{id} [put]
func (h *WealthHandler) UpdateAsset(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid asset id")
	}

	req := new(model.UpdateAssetRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	asset, err := h.wealthService.UpdateAsset(c.Request().Context(), assetID, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, asset)
}

// DeleteAsset
// @Summary Delete an asset
// @Tags Wealth
// @Param id path string true "Asset ID"
// @Security Bearer
// @Success 204
// @Router /api/v1/wealth/assets/{id} [delete]
func (h *WealthHandler) DeleteAsset(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid asset id")
	}

	if err := h.wealthService.DeleteAsset(c.Request().Context(), assetID, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// CreateLiability
// @Summary Create a new liability
// @Tags Wealth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateLiabilityRequest true "Liability info"
// @Success 201 {object} model.Liability
// @Router /api/v1/wealth/liabilities [post]
func (h *WealthHandler) CreateLiability(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	req := new(model.CreateLiabilityRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	liability, err := h.wealthService.CreateLiability(c.Request().Context(), userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, liability)
}

// GetLiabilities
// @Summary Get all liabilities for current user
// @Tags Wealth
// @Produce json
// @Security Bearer
// @Param category query string false "Category filter"
// @Param limit query int false "Limit for pagination"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} model.PaginatedLiabilities
// @Router /api/v1/wealth/liabilities [get]
func (h *WealthHandler) GetLiabilities(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	category := c.QueryParam("category")
	sort := c.QueryParam("sort")
	limit := 0
	offset := 0
	echo.QueryParamsBinder(c).Int("limit", &limit).Int("offset", &offset)

	liabilities, err := h.wealthService.GetLiabilities(c.Request().Context(), userID, category, sort, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, liabilities)
}

// UpdateLiability
// @Summary Update a liability
// @Tags Wealth
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Liability ID"
// @Param request body model.UpdateLiabilityRequest true "Liability info"
// @Success 200 {object} model.Liability
// @Router /api/v1/wealth/liabilities/{id} [put]
func (h *WealthHandler) UpdateLiability(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	liabilityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid liability id")
	}

	req := new(model.UpdateLiabilityRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	liability, err := h.wealthService.UpdateLiability(c.Request().Context(), liabilityID, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, liability)
}

// DeleteLiability
// @Summary Delete a liability
// @Tags Wealth
// @Param id path string true "Liability ID"
// @Security Bearer
// @Success 204
// @Router /api/v1/wealth/liabilities/{id} [delete]
func (h *WealthHandler) DeleteLiability(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	liabilityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid liability id")
	}

	if err := h.wealthService.DeleteLiability(c.Request().Context(), liabilityID, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/service"
)

type IPSHandler struct {
	ipsService service.IPSService
}

func NewIPSHandler(ipsService service.IPSService) *IPSHandler {
	return &IPSHandler{ipsService: ipsService}
}

// GetLatestIPS godoc
// @Summary Get latest IPS
// @Description Retrieve the latest Investment Policy Statement for the authenticated user
// @Tags IPS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.InvestmentPolicy
// @Router /api/v1/ips/latest [get]
func (h *IPSHandler) GetLatestIPS(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	policy, err := h.ipsService.GetLatestIPS(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if policy == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "No IPS found"})
	}

	return c.JSON(http.StatusOK, policy)
}

// UpdateIPS godoc
// @Summary Update IPS
// @Description Update the current Investment Policy Statement
// @Tags IPS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param policy body model.InvestmentPolicy true "Updated Policy"
// @Success 200 {object} map[string]string
// @Router /api/v1/ips/latest [put]
func (h *IPSHandler) UpdateIPS(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	var req model.InvestmentPolicy
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	err := h.ipsService.UpdateIPS(c.Request().Context(), userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

// GenerateIPS godoc
// @Summary Generate IPS via AI
// @Description Let AI analyze current profile and assets to generate a new IPS strategy
// @Tags IPS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.InvestmentPolicy
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ips/generate [post]
func (h *IPSHandler) GenerateIPS(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	type GenerateIPSRequest struct {
		PreferredAssetClasses []string `json:"preferred_asset_classes"`
	}

	var req GenerateIPSRequest
	if err := c.Bind(&req); err != nil {
		// Ignore bind error, just proceed with empty if payload is invalid/missing
	}

	policy, err := h.ipsService.GenerateIPS(c.Request().Context(), userID, req.PreferredAssetClasses)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, policy)
}

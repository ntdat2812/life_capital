package handler

import (
	"net/http"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type FinancialGoalHandler struct {
	goalService service.FinancialGoalService
}

func NewFinancialGoalHandler(goalService service.FinancialGoalService) *FinancialGoalHandler {
	return &FinancialGoalHandler{goalService: goalService}
}

// CreateGoal godoc
// @Summary Create a new financial goal
// @Description Creates a new financial goal for the authenticated user
// @Tags goals
// @Accept json
// @Produce json
// @Param request body model.CreateFinancialGoalRequest true "Goal creation request"
// @Success 201 {object} model.FinancialGoal
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/goals [post]
// @Security BearerAuth
func (h *FinancialGoalHandler) CreateGoal(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user ID")
	}

	var req model.CreateFinancialGoalRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	goal, err := h.goalService.CreateGoal(c.Request().Context(), userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, goal)
}

// ListGoals godoc
// @Summary List financial goals
// @Description Lists all financial goals for the authenticated user
// @Tags goals
// @Produce json
// @Success 200 {array} model.FinancialGoal
// @Failure 401 {object} map[string]string
// @Router /api/v1/goals [get]
// @Security BearerAuth
func (h *FinancialGoalHandler) ListGoals(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user ID")
	}

	goals, err := h.goalService.ListGoals(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// if goals is nil, return empty array instead of null
	if goals == nil {
		goals = []model.FinancialGoal{}
	}
	return c.JSON(http.StatusOK, goals)
}

// UpdateGoal godoc
// @Summary Update a financial goal
// @Description Updates an existing financial goal for the authenticated user
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "Goal ID"
// @Param request body model.UpdateFinancialGoalRequest true "Goal update request"
// @Success 200 {object} model.FinancialGoal
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/goals/{id} [put]
// @Security BearerAuth
func (h *FinancialGoalHandler) UpdateGoal(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user ID")
	}

	goalIDStr := c.Param("id")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid goal ID")
	}

	var req model.UpdateFinancialGoalRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	goal, err := h.goalService.UpdateGoal(c.Request().Context(), userID, goalID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, goal)
}

// DeleteGoal godoc
// @Summary Delete a financial goal
// @Description Deletes a financial goal for the authenticated user
// @Tags goals
// @Produce json
// @Param id path string true "Goal ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/goals/{id} [delete]
// @Security BearerAuth
func (h *FinancialGoalHandler) DeleteGoal(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user ID")
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid goal ID")
	}

	if err := h.goalService.DeleteGoal(c.Request().Context(), userID, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// LinkAsset godoc
// @Summary Link an asset to a financial goal
// @Description Links a specific asset to a financial goal
// @Tags goals
// @Accept json
// @Produce json
// @Param id path string true "Goal ID"
// @Param request body model.LinkAssetToGoalRequest true "Link asset request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/goals/{id}/allocations [post]
// @Security BearerAuth
func (h *FinancialGoalHandler) LinkAsset(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user ID")
	}

	goalIDStr := c.Param("id")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid goal ID")
	}

	var req model.LinkAssetToGoalRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.goalService.LinkAsset(c.Request().Context(), userID, goalID, req.AssetID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "asset linked successfully"})
}

// UnlinkAsset godoc
// @Summary Unlink an asset from a financial goal
// @Description Removes an asset link from a financial goal
// @Tags goals
// @Produce json
// @Param id path string true "Goal ID"
// @Param asset_id path string true "Asset ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/goals/{id}/allocations/{asset_id} [delete]
// @Security BearerAuth
func (h *FinancialGoalHandler) UnlinkAsset(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user ID")
	}

	goalIDStr := c.Param("id")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid goal ID")
	}

	assetIDStr := c.Param("asset_id")
	assetID, err := uuid.Parse(assetIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid asset ID")
	}

	if err := h.goalService.UnlinkAsset(c.Request().Context(), userID, goalID, assetID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

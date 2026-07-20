package handler

import (
	"net/http"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CashflowHandler struct {
	cashflowService service.CashflowService
}

func NewCashflowHandler(cashflowService service.CashflowService) *CashflowHandler {
	return &CashflowHandler{cashflowService: cashflowService}
}

// CreateIncomeStream
// @Summary Create an income stream
// @Tags Cashflow
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateIncomeStreamRequest true "Income Stream info"
// @Success 201 {object} model.IncomeStream
// @Router /api/v1/cashflow/income [post]
func (h *CashflowHandler) CreateIncomeStream(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	req := new(model.CreateIncomeStreamRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	income, err := h.cashflowService.CreateIncomeStream(c.Request().Context(), userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, income)
}

// GetIncomeStreams
// @Summary Get all income streams
// @Tags Cashflow
// @Produce json
// @Security Bearer
// @Success 200 {array} model.IncomeStream
// @Router /api/v1/cashflow/income [get]
func (h *CashflowHandler) GetIncomeStreams(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	incomes, err := h.cashflowService.GetIncomeStreams(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, incomes)
}

// UpdateIncomeStream
// @Summary Update an income stream
// @Tags Cashflow
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Income ID"
// @Param request body model.UpdateIncomeStreamRequest true "Income Stream info"
// @Success 200 {object} model.IncomeStream
// @Router /api/v1/cashflow/income/{id} [put]
func (h *CashflowHandler) UpdateIncomeStream(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid income id")
	}

	req := new(model.UpdateIncomeStreamRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	income, err := h.cashflowService.UpdateIncomeStream(c.Request().Context(), id, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, income)
}

// DeleteIncomeStream
// @Summary Delete an income stream
// @Tags Cashflow
// @Security Bearer
// @Param id path string true "Income ID"
// @Success 204
// @Router /api/v1/cashflow/income/{id} [delete]
func (h *CashflowHandler) DeleteIncomeStream(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid income id")
	}

	if err := h.cashflowService.DeleteIncomeStream(c.Request().Context(), id, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// CreateDependent
// @Summary Create a dependent
// @Tags Cashflow
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body model.CreateDependentRequest true "Dependent info"
// @Success 201 {object} model.Dependent
// @Router /api/v1/cashflow/dependents [post]
func (h *CashflowHandler) CreateDependent(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	req := new(model.CreateDependentRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dep, err := h.cashflowService.CreateDependent(c.Request().Context(), userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, dep)
}

// GetDependents
// @Summary Get all dependents
// @Tags Cashflow
// @Produce json
// @Security Bearer
// @Success 200 {array} model.Dependent
// @Router /api/v1/cashflow/dependents [get]
func (h *CashflowHandler) GetDependents(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	deps, err := h.cashflowService.GetDependents(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, deps)
}

// UpdateDependent
// @Summary Update a dependent
// @Tags Cashflow
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Dependent ID"
// @Param request body model.UpdateDependentRequest true "Dependent info"
// @Success 200 {object} model.Dependent
// @Router /api/v1/cashflow/dependents/{id} [put]
func (h *CashflowHandler) UpdateDependent(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid dependent id")
	}

	req := new(model.UpdateDependentRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dep, err := h.cashflowService.UpdateDependent(c.Request().Context(), id, userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dep)
}

// DeleteDependent
// @Summary Delete a dependent
// @Tags Cashflow
// @Security Bearer
// @Param id path string true "Dependent ID"
// @Success 204
// @Router /api/v1/cashflow/dependents/{id} [delete]
func (h *CashflowHandler) DeleteDependent(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid dependent id")
	}

	if err := h.cashflowService.DeleteDependent(c.Request().Context(), id, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

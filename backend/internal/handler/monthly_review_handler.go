package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/datnguyen/life_capital/backend/internal/service"
	"github.com/datnguyen/life_capital/backend/internal/model"
)

type MonthlyReviewHandler struct {
	monthlyReviewService service.MonthlyReviewService
}

func NewMonthlyReviewHandler(monthlyReviewService service.MonthlyReviewService) *MonthlyReviewHandler {
	return &MonthlyReviewHandler{
		monthlyReviewService: monthlyReviewService,
	}
}

type GenerateReviewRequest struct {
	NewInvestmentAmount float64 `json:"new_investment_amount"`
}

// GenerateReview godoc
// @Summary Generate Monthly Review
// @Description Analyzes portfolio, watchlist, liabilities and new capital to generate a recommendation
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GenerateReviewRequest true "New Investment Amount"
// @Success 200 {object} model.MonthlyReview
// @Router /api/v1/reviews/generate [post]
func (h *MonthlyReviewHandler) GenerateReview(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	var req GenerateReviewRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	review, err := h.monthlyReviewService.GenerateReview(c.Request().Context(), userID, req.NewInvestmentAmount)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, review)
}

// SaveReview godoc
// @Summary Save Monthly Review
// @Description Saves a generated monthly review
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.MonthlyReview true "Monthly Review Data"
// @Success 200 {object} model.MonthlyReview
// @Router /api/v1/reviews [post]
func (h *MonthlyReviewHandler) SaveReview(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	var review model.MonthlyReview
	if err := c.Bind(&review); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	
	// Ensure the user ID matches
	if review.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "user ID mismatch")
	}

	if err := h.monthlyReviewService.SaveReview(c.Request().Context(), &review); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, review)
}

// GetReviewHistory godoc
// @Summary Get Review History
// @Description Get all past monthly reviews for the user
// @Tags Reviews
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.MonthlyReview
// @Router /api/v1/reviews [get]
func (h *MonthlyReviewHandler) GetReviewHistory(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	reviews, err := h.monthlyReviewService.GetReviewHistory(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if reviews == nil {
		reviews = make([]model.MonthlyReview, 0)
	}

	return c.JSON(http.StatusOK, reviews)
}

// GetReviewByMonth godoc
// @Summary Get Review By Month
// @Description Get a specific monthly review
// @Tags Reviews
// @Produce json
// @Security BearerAuth
// @Param month path string true "Month (YYYY-MM-01)"
// @Success 200 {object} model.MonthlyReview
// @Router /api/v1/reviews/{month} [get]
func (h *MonthlyReviewHandler) GetReviewByMonth(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	month := c.Param("month")
	if month == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "month is required")
	}

	review, err := h.monthlyReviewService.GetReviewByMonth(c.Request().Context(), userID, month)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "review not found")
	}

	return c.JSON(http.StatusOK, review)
}

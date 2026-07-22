package handler

import (
	"github.com/datnguyen/life_capital/backend/internal/ai"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/service"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PortfolioHandler struct {
	portfolioService service.PortfolioService
	aiProviders      []ai.AIProvider
}

func NewPortfolioHandler(portfolioService service.PortfolioService, aiProviders []ai.AIProvider) *PortfolioHandler {
	return &PortfolioHandler{
		portfolioService: portfolioService,
		aiProviders:      aiProviders,
	}
}

// GetInvestableAssets godoc
// @Summary Get user's investable assets
// @Description Get a list of assets that are considered investments (stocks, crypto, real_estate, gold)
// @Tags portfolio
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Asset
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/portfolio [get]
func (h *PortfolioHandler) GetInvestableAssets(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	assets, err := h.portfolioService.GetInvestableAssets(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, assets)
}

// GetWatchlist godoc
// @Summary Get user's watchlist
// @Tags watchlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.WatchlistItem
// @Router /api/v1/watchlist [get]
func (h *PortfolioHandler) GetWatchlist(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	items, err := h.portfolioService.GetWatchlistByUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, items)
}

// CreateWatchlistItem godoc
// @Summary Create a watchlist item
// @Tags watchlist
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.WatchlistItem true "Watchlist Item"
// @Success 201 {object} model.WatchlistItem
// @Router /api/v1/watchlist [post]
func (h *PortfolioHandler) CreateWatchlistItem(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	var item model.WatchlistItem
	if err := c.Bind(&item); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	item.UserID = userID
	if err := h.portfolioService.CreateWatchlistItem(c.Request().Context(), &item); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, item)
}

// DeleteWatchlistItem godoc
// @Summary Delete a watchlist item
// @Tags watchlist
// @Security BearerAuth
// @Param id path string true "Watchlist ID"
// @Success 204
// @Router /api/v1/watchlist/{id} [delete]
func (h *PortfolioHandler) DeleteWatchlistItem(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ID format")
	}

	if err := h.portfolioService.DeleteWatchlistItem(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// GetTheses godoc
// @Summary Get user's investment theses
// @Tags thesis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.InvestmentThesis
// @Router /api/v1/theses [get]
func (h *PortfolioHandler) GetTheses(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	theses, err := h.portfolioService.GetThesesByUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, theses)
}

// GetThesisByTicker godoc
// @Summary Get thesis by ticker
// @Tags thesis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ticker path string true "Ticker"
// @Success 200 {object} model.InvestmentThesis
// @Router /api/v1/theses/{ticker} [get]
func (h *PortfolioHandler) GetThesisByTicker(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	ticker := c.Param("ticker")
	thesis, err := h.portfolioService.GetThesisByTicker(c.Request().Context(), userID, ticker)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "thesis not found")
	}
	return c.JSON(http.StatusOK, thesis)
}

// CreateThesis godoc
// @Summary Create an investment thesis
// @Tags thesis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.InvestmentThesis true "Thesis Data"
// @Success 201 {object} model.InvestmentThesis
// @Router /api/v1/theses [post]
func (h *PortfolioHandler) CreateThesis(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	var thesis model.InvestmentThesis
	if err := c.Bind(&thesis); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	thesis.UserID = userID
	clampScores(&thesis)
	if err := h.portfolioService.CreateThesis(c.Request().Context(), &thesis); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, thesis)
}

// UpdateThesis godoc
// @Summary Update an investment thesis
// @Tags thesis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Thesis ID"
// @Param request body model.InvestmentThesis true "Thesis Data"
// @Success 200 {object} model.InvestmentThesis
// @Router /api/v1/theses/{id} [put]
func (h *PortfolioHandler) UpdateThesis(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user token")
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ID format")
	}

	var thesis model.InvestmentThesis
	if err := c.Bind(&thesis); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	thesis.ID = id
	thesis.UserID = userID
	clampScores(&thesis)
	if err := h.portfolioService.UpdateThesis(c.Request().Context(), &thesis); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, thesis)
}

// DeleteThesis godoc
// @Summary Delete an investment thesis
// @Tags thesis
// @Security BearerAuth
// @Param id path string true "Thesis ID"
// @Success 204
// @Router /api/v1/theses/{id} [delete]
func (h *PortfolioHandler) DeleteThesis(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ID format")
	}

	if err := h.portfolioService.DeleteThesis(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

// GenerateThesisAI godoc
// @Summary Generate thesis using AI
// @Tags thesis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.ThesisGenerationRequest true "Generation Request"
// @Success 200 {object} model.InvestmentThesis
// @Router /api/v1/theses/generate [post]
func (h *PortfolioHandler) GenerateThesisAI(c echo.Context) error {
	var req model.ThesisGenerationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.Ticker == "" || req.CompanyName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Thiếu thông tin Ticker hoặc Tên công ty/tài sản")
	}

	thesis, err := ai.ExecuteWithFallback(h.aiProviders, func(p ai.AIProvider) (*model.InvestmentThesis, error) {
		return p.GenerateThesis(c.Request().Context(), &req)
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, thesis)
}
func clampScores(thesis *model.InvestmentThesis) {
	if thesis.ConvictionScore < 1 {
		thesis.ConvictionScore = 5
	} else if thesis.ConvictionScore > 10 {
		thesis.ConvictionScore = 10
	}
	if thesis.QualityScore < 1 {
		thesis.QualityScore = 5
	} else if thesis.QualityScore > 10 {
		thesis.QualityScore = 10
	}
	if thesis.ValuationScore < 1 {
		thesis.ValuationScore = 5
	} else if thesis.ValuationScore > 10 {
		thesis.ValuationScore = 10
	}
}

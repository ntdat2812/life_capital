package handler

import (
	"net/http"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type TimelineHandler struct {
	service *service.TimelineService
}

func NewTimelineHandler(service *service.TimelineService) *TimelineHandler {
	return &TimelineHandler{service: service}
}

// AnalyzeEvent godoc
// @Summary Analyze a new life event
// @Description Send event text to AI to analyze its impact on the user's financial profile
// @Tags timeline
// @Accept json
// @Produce json
// @Param req body model.LogEventRequest true "Event Text"
// @Success 200 {object} ai.LifeEventAnalysisResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/profile/timeline/analyze [post]
func (h *TimelineHandler) AnalyzeEvent(c echo.Context) error {
	userID := c.Get("user_id").(string)
	
	var req model.LogEventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if req.EventText == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "event_text is required"})
	}

	result, err := h.service.AnalyzeEvent(c.Request().Context(), userID, req.EventText)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// ConfirmEvent godoc
// @Summary Confirm and save a life event
// @Description Apply the analyzed (and potentially edited) life event to the database, cascading profile changes
// @Tags timeline
// @Accept json
// @Produce json
// @Param req body model.ConfirmEventRequest true "Event Confirmation"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/profile/timeline/confirm [post]
func (h *TimelineHandler) ConfirmEvent(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req model.ConfirmEventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := h.service.ConfirmEvent(c.Request().Context(), userID, &req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Life event saved and cascaded successfully"})
}

// GetTimeline godoc
// @Summary Get user's life timeline
// @Description Retrieve a list of all life events logged by the user
// @Tags timeline
// @Produce json
// @Success 200 {array} model.LifeEvent
// @Failure 500 {object} map[string]string
// @Router /api/v1/profile/timeline [get]
func (h *TimelineHandler) GetTimeline(c echo.Context) error {
	userID := c.Get("user_id").(string)

	events, err := h.service.GetTimeline(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// if empty, return empty array instead of null
	if events == nil {
		events = []model.LifeEvent{}
	}

	return c.JSON(http.StatusOK, events)
}

package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthHandler handles system health endpoints.
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler instance.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary      Health Check
// @Description  Check backend health status
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/v1/health [get]
func (h *HealthHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"version": "v1.0.0",
	})
}

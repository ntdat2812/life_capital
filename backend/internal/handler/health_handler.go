package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// HealthHandler handles system health endpoints.
type HealthHandler struct {
	db *pgxpool.Pool
}

// NewHealthHandler creates a new HealthHandler instance.
func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck godoc
// @Summary      Health Check
// @Description  Check backend health status
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/v1/health [get]
func (h *HealthHandler) HealthCheck(c echo.Context) error {
	// Create a timeout context for the DB ping
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "unhealthy",
			"message": "database connection failed",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":   "healthy",
		"database": "connected",
		"version":  "v1.0.0",
	})
}

package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/datnguyen/life_capital/backend/docs"
	"github.com/datnguyen/life_capital/backend/internal/handler"
)

// @title Life Capital (WealthOS) API
// @version 1.0
// @description API server for Personal Wealth Operating System.
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize Echo instance
	e := echo.New()

	// Default Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Swagger documentation endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()

	// Register routes
	registerRoutes(e, healthHandler)

	// Get PORT from environment variable (default: 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Life Capital server on port %s...", port)

	// Start server
	if err := e.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Shutting down server with error: %v", err)
	}
}

// registerRoutes maps all API routes to their handlers.
func registerRoutes(e *echo.Echo, healthHandler *handler.HealthHandler) {
	api := e.Group("/api/v1")

	// System
	api.GET("/health", healthHandler.HealthCheck)
}

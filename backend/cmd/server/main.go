package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

	// Health check endpoint
	e.GET("/api/v1/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "healthy",
			"version": "v1.0.0",
		})
	})

	// Get PORT from environment variable (default: 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Life Capital server on port %s...", port)
	
	// Start server
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Shutting down server with error: %v", err)
	}
}

// Ensure context propagation for cancellation
func processRequest(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

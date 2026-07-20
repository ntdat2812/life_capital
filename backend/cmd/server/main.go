package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/datnguyen/life_capital/backend/docs"
	"github.com/datnguyen/life_capital/backend/internal/handler"
	"github.com/datnguyen/life_capital/backend/internal/repository"
	"github.com/datnguyen/life_capital/backend/internal/service"
	customMiddleware "github.com/datnguyen/life_capital/backend/internal/middleware"
)

// @title Life Capital (WealthOS) API
// @version 1.0
// @description API server for Personal Wealth Operating System.
// @host localhost:8080
// @BasePath /
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Database Connection
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	dbPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

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

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbPool)

	// Initialize services
	jwtSecret := os.Getenv("JWT_SECRET")
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	authService := service.NewAuthService(userRepo, jwtSecret, googleClientID)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)

	// Register routes
	registerRoutes(e, healthHandler, authHandler)

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
func registerRoutes(e *echo.Echo, healthHandler *handler.HealthHandler, authHandler *handler.AuthHandler) {
	api := e.Group("/api/v1")

	// Public routes
	api.GET("/health", healthHandler.HealthCheck)
	
	auth := api.Group("/auth")
	auth.POST("/signup", authHandler.Signup)
	auth.POST("/login", authHandler.Login)
	auth.POST("/google", authHandler.GoogleLogin)

	// Protected routes
	protected := api.Group("")
	protected.Use(customMiddleware.JWTMiddleware(os.Getenv("JWT_SECRET")))
	protected.GET("/auth/me", authHandler.Me)
}

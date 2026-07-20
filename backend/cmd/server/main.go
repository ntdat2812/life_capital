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
	"github.com/datnguyen/life_capital/backend/internal/ai"
	"github.com/datnguyen/life_capital/backend/internal/handler"
	customMiddleware "github.com/datnguyen/life_capital/backend/internal/middleware"
	"github.com/datnguyen/life_capital/backend/internal/repository"
	"github.com/datnguyen/life_capital/backend/internal/service"
)

// @title Life Capital (WealthOS) API
// @version 1.0
// @description API server for Personal Wealth Operating System.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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
	investorProfileRepo := repository.NewInvestorProfileRepository(dbPool)
	assetRepo := repository.NewAssetRepository(dbPool)
	liabilityRepo := repository.NewLiabilityRepository(dbPool)
	incomeRepo := repository.NewIncomeRepository(dbPool)
	dependentRepo := repository.NewDependentRepository(dbPool)

	// Initialize AI Provider
	geminiProvider, err := ai.NewGeminiProvider()
	if err != nil {
		log.Fatalf("Failed to initialize Gemini provider: %v", err)
	}

	// Initialize services
	jwtSecret := os.Getenv("JWT_SECRET")
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	authService := service.NewAuthService(userRepo, jwtSecret, googleClientID)
	profileService := service.NewInvestorProfileService(investorProfileRepo, incomeRepo, geminiProvider)
	wealthService := service.NewWealthService(assetRepo, liabilityRepo, userRepo)
	cashflowService := service.NewCashflowService(incomeRepo, dependentRepo)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	profileHandler := handler.NewProfileHandler(profileService)
	wealthHandler := handler.NewWealthHandler(wealthService)
	cashflowHandler := handler.NewCashflowHandler(cashflowService)

	// Register routes
	registerRoutes(e, healthHandler, authHandler, profileHandler, wealthHandler, cashflowHandler)

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
func registerRoutes(e *echo.Echo, healthHandler *handler.HealthHandler, authHandler *handler.AuthHandler, profileHandler *handler.ProfileHandler, wealthHandler *handler.WealthHandler, cashflowHandler *handler.CashflowHandler) {
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

	// Profile routes
	profile := protected.Group("/profile")
	profile.POST("/onboarding", profileHandler.Onboarding)
	profile.GET("/me", profileHandler.GetMe)
	profile.PUT("/me", profileHandler.UpdateProfile)

	// Cashflow routes
	cashflowGroup := protected.Group("/cashflow")
	cashflowGroup.POST("/income", cashflowHandler.CreateIncomeStream)
	cashflowGroup.GET("/income", cashflowHandler.GetIncomeStreams)
	cashflowGroup.PUT("/income/:id", cashflowHandler.UpdateIncomeStream)
	cashflowGroup.DELETE("/income/:id", cashflowHandler.DeleteIncomeStream)

	cashflowGroup.POST("/dependents", cashflowHandler.CreateDependent)
	cashflowGroup.GET("/dependents", cashflowHandler.GetDependents)
	cashflowGroup.PUT("/dependents/:id", cashflowHandler.UpdateDependent)
	cashflowGroup.DELETE("/dependents/:id", cashflowHandler.DeleteDependent)

	// Wealth routes
	wealthGroup := protected.Group("/wealth")
	wealthGroup.GET("/net-worth", wealthHandler.GetNetWorthSummary)
	wealthGroup.POST("/assets", wealthHandler.CreateAsset)
	wealthGroup.GET("/assets", wealthHandler.GetAssets)
	wealthGroup.PUT("/assets/:id", wealthHandler.UpdateAsset)
	wealthGroup.DELETE("/assets/:id", wealthHandler.DeleteAsset)

	wealthGroup.POST("/liabilities", wealthHandler.CreateLiability)
	wealthGroup.GET("/liabilities", wealthHandler.GetLiabilities)
	wealthGroup.PUT("/liabilities/:id", wealthHandler.UpdateLiability)
	wealthGroup.DELETE("/liabilities/:id", wealthHandler.DeleteLiability)
}

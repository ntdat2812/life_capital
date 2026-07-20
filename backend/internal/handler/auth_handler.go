package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/service"
)

type AuthHandler struct {
	authService         *service.AuthService
	localAuthenticator  *service.LocalAuthenticator
	googleAuthenticator *service.GoogleAuthenticator
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		localAuthenticator:  &service.LocalAuthenticator{Service: authService},
		googleAuthenticator: &service.GoogleAuthenticator{Service: authService},
	}
}

// Signup godoc
// @Summary      Local Account Signup
// @Description  Register a new user using Email and Password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.SignupRequest true "Signup information"
// @Success      201  {object}  model.User
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /api/v1/auth/signup [post]
func (h *AuthHandler) Signup(c echo.Context) error {
	req := new(model.SignupRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	}

	user, err := h.authService.Signup(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return echo.NewHTTPError(http.StatusConflict, "Email already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary      Local Account Login
// @Description  Login using Email and Password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.LoginRequest true "Login credentials"
// @Success      200  {object}  model.AuthResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	req := new(model.LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	res, err := h.localAuthenticator.Authenticate(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	token, err := h.authService.GenerateJWT(res.User.ID.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(http.StatusOK, model.AuthResponse{
		Token: token,
		User:  res.User,
	})
}

// GoogleLogin godoc
// @Summary      Google Login
// @Description  Authenticate using Google ID Token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.GoogleLoginRequest true "Google ID Token"
// @Success      200  {object}  model.AuthResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/v1/auth/google [post]
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	req := new(model.GoogleLoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	res, err := h.googleAuthenticator.Authenticate(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Google token: " + err.Error())
	}

	token, err := h.authService.GenerateJWT(res.User.ID.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(http.StatusOK, model.AuthResponse{
		Token: token,
		User:  res.User,
	})
}

// Me godoc
// @Summary      Account Information
// @Description  Get current account information via JWT Token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  model.User
// @Failure      401  {object}  map[string]string
// @Router       /api/v1/auth/me [get]
func (h *AuthHandler) Me(c echo.Context) error {
	userID := c.Get("user_id").(string)
	user, err := h.authService.GetUserByID(c.Request().Context(), userID)
	if err != nil || user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
	}
	return c.JSON(http.StatusOK, user)
}

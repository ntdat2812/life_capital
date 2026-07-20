package handler

import (
	"net/http"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	profileService *service.InvestorProfileService
}

func NewProfileHandler(profileService *service.InvestorProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

// Onboarding godoc
// @Summary      Submit AI Onboarding Interview
// @Description  Submit the 8-step chat history to let AI analyze and create an Investor Profile.
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body model.OnboardingRequest true "Chat history data"
// @Success      201  {object}  model.InvestorProfile
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/profile/onboarding [post]
func (h *ProfileHandler) Onboarding(c echo.Context) error {
	userID := c.Get("user_id").(string)

	req := new(model.OnboardingRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	profile, err := h.profileService.ProcessOnboarding(c.Request().Context(), userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, profile)
}

// GetMe godoc
// @Summary      Get Active Profile
// @Description  Get the current active investor profile for the logged in user.
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  model.InvestorProfile
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/profile/me [get]
func (h *ProfileHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(string)

	profile, err := h.profileService.GetActiveProfile(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if profile == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Profile not found")
	}

	return c.JSON(http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary      Update Active Profile
// @Description  Update the current active investor profile for the logged in user.
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body model.UpdateProfileRequest true "Profile Update Data"
// @Success      200  {object}  model.InvestorProfile
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/profile/me [put]
func (h *ProfileHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(string)

	req := new(model.UpdateProfileRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	profile, err := h.profileService.UpdateProfile(c.Request().Context(), userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, profile)
}

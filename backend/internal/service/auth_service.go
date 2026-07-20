package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"

	"github.com/datnguyen/life_capital/backend/internal/model"
	"github.com/datnguyen/life_capital/backend/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthService struct {
	repo           *repository.UserRepository
	jwtSecret      string
	googleClientID string
}

func NewAuthService(repo *repository.UserRepository, jwtSecret, googleClientID string) *AuthService {
	return &AuthService{
		repo:           repo,
		jwtSecret:      jwtSecret,
		googleClientID: googleClientID,
	}
}

// GenerateJWT creates a new JWT token for a user
func (s *AuthService) GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// Signup handles local user registration
func (s *AuthService) Signup(ctx context.Context, req *model.SignupRequest) (*model.User, error) {
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedStr := string(hashed)

	user := &model.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: &hashedStr,
		AuthProvider: "local",
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID fetches a user by their UUID
func (s *AuthService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

// -- Strategies --

type LocalAuthenticator struct {
	Service *AuthService
}

func (a *LocalAuthenticator) Authenticate(ctx context.Context, payload interface{}) (*AuthenticateResult, error) {
	req, ok := payload.(*model.LoginRequest)
	if !ok {
		return nil, fmt.Errorf("invalid payload type for local authenticator")
	}

	user, err := a.Service.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil || user.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return &AuthenticateResult{User: user}, nil
}

type GoogleAuthenticator struct {
	Service *AuthService
}

func (a *GoogleAuthenticator) Authenticate(ctx context.Context, payload interface{}) (*AuthenticateResult, error) {
	req, ok := payload.(*model.GoogleLoginRequest)
	if !ok {
		return nil, fmt.Errorf("invalid payload type for google authenticator")
	}

	// Verify ID token
	payloadClaim, err := idtoken.Validate(ctx, req.IDToken, a.Service.googleClientID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	email, ok := payloadClaim.Claims["email"].(string)
	if !ok {
		return nil, errors.New("email not found in google token")
	}

	name, _ := payloadClaim.Claims["name"].(string)

	user, err := a.Service.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// Auto-signup via Google
		googleID := payloadClaim.Subject
		user = &model.User{
			Email:        email,
			Name:         name,
			AuthProvider: "google",
			GoogleID:     &googleID,
		}
		err = a.Service.repo.CreateUser(ctx, user)
		if err != nil {
			return nil, err
		}
	} else {
		// User exists, maybe local, maybe google. Ensure linking.
		if user.GoogleID == nil || *user.GoogleID != payloadClaim.Subject {
			err = a.Service.repo.UpdateGoogleID(ctx, user.ID.String(), payloadClaim.Subject)
			if err != nil {
				return nil, err
			}
			googleID := payloadClaim.Subject
			user.GoogleID = &googleID
		}
	}

	return &AuthenticateResult{User: user}, nil
}

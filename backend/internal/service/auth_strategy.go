package service

import (
	"context"

	"github.com/datnguyen/life_capital/backend/internal/model"
)

type AuthenticateResult struct {
	User *model.User
}

// Authenticator defines the Strategy for login
type Authenticator interface {
	Authenticate(ctx context.Context, payload interface{}) (*AuthenticateResult, error)
}

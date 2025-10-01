package services

import (
	"context"

	"github.com/shravanskie/vegetable-backend/internal/apperrors"
	"github.com/shravanskie/vegetable-backend/internal/models"
)

type UserService interface {
	SignupWithEmail(ctx context.Context, firstName, lastName, email, phone, password string) (*models.User, *apperrors.AppError)
	SignupOrLoginWithGoogle(ctx context.Context, idToken string) (*models.User, string, *apperrors.AppError) // returns user, token, error
	Login(ctx context.Context, username, password string) (*models.LoginResult, *apperrors.AppError)
	Logout(ctx context.Context, userID uint) *apperrors.AppError
}

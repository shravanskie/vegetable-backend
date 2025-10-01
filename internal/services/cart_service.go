package services

import (
	"github.com/shravanskie/vegetable-backend/internal/apperrors"
	"github.com/shravanskie/vegetable-backend/internal/models"
)

type CartService interface {
	AddToCart(userID uint, vegetableID uint, quantity int) (*models.Cart, *apperrors.AppError)
	// ListCart(userID uint) ([]models.Cart, *apperrors.AppError)
	// RemoveFromCart(userID uint, vegetableID uint, quantity int) (*models.Cart, *apperrors.AppError)x
}

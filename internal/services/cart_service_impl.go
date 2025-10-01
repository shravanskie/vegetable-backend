package services

import (
	"errors"
	"net/http"

	"github.com/shravanskie/vegetable-backend/internal/apperrors"
	"github.com/shravanskie/vegetable-backend/internal/db"
	"github.com/shravanskie/vegetable-backend/internal/models"
	"gorm.io/gorm"
)

type CartServiceImpl struct {
}

func NewCartService() CartService {
	return &CartServiceImpl{}
}

// AddToCart inserts or updates the cart item
func (s *CartServiceImpl) AddToCart(userID uint, vegetableID uint, quantity int) (*models.Cart, *apperrors.AppError) {
	if quantity <= 0 {
		return nil, apperrors.New(http.StatusBadRequest, "quantity must be greater than zero")
	}

	var cart *models.Cart
	err := db.DB.Where("user_id = ? AND vegetable_id = ?", userID, vegetableID).First(cart).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Insert new cart row
		cart = &models.Cart{
			UserID:      userID,
			VegetableID: vegetableID,
			Quantity:    quantity,
		}
		err := db.DB.Create(cart).Error
		if err != nil {
			return nil, apperrors.New(http.StatusInternalServerError, "failed to add to cart")
		}
		return cart, nil
	} else if err != nil {
		return nil, apperrors.New(http.StatusInternalServerError, err.Error())
	}

	// Update existing quantity
	cart.Quantity += quantity
	err = db.DB.Create(&cart).Error
	if err != nil {
		return nil, apperrors.New(http.StatusInternalServerError, "failed to update to cart")
	}
	return cart, nil
}

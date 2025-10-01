package services

import (
	"context"

	"github.com/shravanskie/vegetable-backend/internal/apperrors"
	"github.com/shravanskie/vegetable-backend/internal/models"
)

type VegetableService interface {
	ListVegetables(ctx context.Context) ([]models.Vegetable, *apperrors.AppError)
	GetVegetable(ctx context.Context, id uint) (*models.Vegetable, *apperrors.AppError)
	AddVegetable(ctx context.Context, veg *models.Vegetable) *apperrors.AppError
	UpdateVegetable(ctx context.Context, veg *models.Vegetable) *apperrors.AppError
	DeleteVegetable(ctx context.Context, id uint) *apperrors.AppError
}

package services

import (
	"context"
	"net/http"

	"github.com/shravanskie/vegetable-backend/internal/apperrors"
	"github.com/shravanskie/vegetable-backend/internal/db"
	"github.com/shravanskie/vegetable-backend/internal/models"
)

type VegetableServiceImpl struct {
	// Add any dependencies like DB connection here
}

func NewVegetableService() VegetableService {
	return &VegetableServiceImpl{}
}

func (s *VegetableServiceImpl) ListVegetables(ctx context.Context) ([]models.Vegetable, *apperrors.AppError) {
	// Implement logic to list vegetables from DB
	var vegetables []models.Vegetable
	result := db.DB.Find(&vegetables)
	if result.Error != nil {
		return nil, apperrors.New(http.StatusInternalServerError, "failed to fetch vegetables")
	}
	return vegetables, nil
}

func (s *VegetableServiceImpl) GetVegetable(ctx context.Context, id uint) (*models.Vegetable, *apperrors.AppError) {
	// Implement logic to get a single vegetable by ID from DB
	return nil, nil
}

func (s *VegetableServiceImpl) AddVegetable(ctx context.Context, veg *models.Vegetable) *apperrors.AppError {
	// Implement logic to add a new vegetable to DB
	var existing models.Vegetable
	// check if vegetable already exists
	if err := db.DB.Where("name = ?", veg.Name).First(&existing).Error; err == nil {
		return apperrors.New(apperrors.ErrVegetableExists, "vegetable with this name already exists")
	}
	if err := db.DB.Create(&veg).Error; err != nil {
		return apperrors.New(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (s *VegetableServiceImpl) UpdateVegetable(ctx context.Context, veg *models.Vegetable) *apperrors.AppError {
	// Implement logic to update an existing vegetable in DB
	return nil
}

func (s *VegetableServiceImpl) DeleteVegetable(ctx context.Context, id uint) *apperrors.AppError {
	// Implement logic to delete a vegetable by ID from DB
	return nil
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shravanskie/vegetable-backend/internal/models"
	"github.com/shravanskie/vegetable-backend/internal/services"
)

type VegetableHandler struct {
	VegetableService services.VegetableService
}

func NewVegetableHandler(vegetableService services.VegetableService) *VegetableHandler {
	return &VegetableHandler{VegetableService: vegetableService}
}

func (s *VegetableHandler) AddVegetable(c *gin.Context) {
	var input models.VegetableInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false})
		return
	}

	vegetable := &models.Vegetable{
		Name:         input.Name,
		Unit:         input.Unit,
		UnitQuantity: input.UnitQuantity,
		ImageURL:     input.ImagePath,
	}

	err := s.VegetableService.AddVegetable(c.Request.Context(), vegetable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Message, "success": false, "errorCode": err.Code})
		return
	}

	c.JSON(200, gin.H{"success": true, "vegetable": vegetable})
}

// GET /api/vegetables
func (h *VegetableHandler) ListVegetables(c *gin.Context) {
	vegetables, err := h.VegetableService.ListVegetables(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Message, "success": false, "errorCode": err.Code})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": vegetables})
}

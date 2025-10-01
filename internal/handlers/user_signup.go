package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shravanskie/vegetable-backend/internal/apperrors"
	"github.com/shravanskie/vegetable-backend/internal/models"
	"github.com/shravanskie/vegetable-backend/internal/services"
	"github.com/shravanskie/vegetable-backend/internal/utils"
)

type SignupHandler struct {
	UserService services.UserService
}

func NewSignupHandler(userService services.UserService) *SignupHandler {
	return &SignupHandler{UserService: userService}
}

// Email/Password Signup
func (h *SignupHandler) Signup(c *gin.Context) {
	var req struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Phone     string `json:"phone" binding:"required"`
		Password  string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "success": false, "errorCode": apperrors.ErrUserSignupFailed})
		return
	}

	user, svcErr := h.UserService.SignupWithEmail(c.Request.Context(), req.FirstName, req.LastName, req.Email, req.Phone, req.Password)
	if svcErr != nil {
		// If service returned an AppError, use its code and message
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":   false,
			"error":     svcErr.Message,
			"errorCode": svcErr.Code,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"userId":  user.ID,
		"message": "Signup successful",
	})
}

// Google Signup/Login
func (h *SignupHandler) GoogleSignup(c *gin.Context) {
	var req struct {
		IdToken string `json:"idToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	user, token, svcErr := h.UserService.SignupOrLoginWithGoogle(c.Request.Context(), req.IdToken)
	if svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":   false,
			"error":     svcErr.Message,
			"errorCode": svcErr.Code,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"userId":  user.ID,
		"token":   token,
		"message": "Google signup/login successful",
	})
}

func (h *SignupHandler) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request"})
		return
	}
	loginResult, err := h.UserService.Login(c.Request.Context(), input.Identifier, input.Password)
	if err != nil {
		log.Println("Login error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Message, "errorCode": err.Code})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"token":     loginResult.Token,
		"expiresIn": loginResult.ExpiresIn,
		"userId":    loginResult.User.ID,
	})
}

func (h *SignupHandler) ValidateToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}

	claims, err := utils.ValidateJWT(strings.TrimPrefix(tokenString, "Bearer "))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "claims": claims.Claims})
}

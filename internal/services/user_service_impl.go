package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shravanskie/vegetable-backend/internal/apperrors"
	"github.com/shravanskie/vegetable-backend/internal/db"
	"github.com/shravanskie/vegetable-backend/internal/models"
	"github.com/shravanskie/vegetable-backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

var jwtSecret = []byte("supersecretkey") // 🔒 replace with env var

type userServiceImpl struct {
	googleClientID string
}

func NewUserService(googleClientID string) UserService {
	return &userServiceImpl{googleClientID: googleClientID}
}

// Email/Password Signup
func (s *userServiceImpl) SignupWithEmail(ctx context.Context, firstName, lastName, email, phone, password string) (*models.User, *apperrors.AppError) {
	// Check if user exists
	var existing models.User
	if err := db.DB.Where("email = ? OR phone = ?", email, phone).First(&existing).Error; err == nil {
		return nil, apperrors.New(apperrors.ErrUserExists, "user with email or phone already exists")
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashedPassword),
	}

	if err := db.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, apperrors.Wrap(apperrors.ErrFailedToCreateUser, "failed to create user", err)
	}

	return user, nil
}

// Google Signup/Login
func (s *userServiceImpl) SignupOrLoginWithGoogle(ctx context.Context, idToken string) (*models.User, string, *apperrors.AppError) {
	// 1. Verify Google ID token
	payload, err := idtoken.Validate(context.Background(), idToken, s.googleClientID)
	if err != nil {
		return nil, "", apperrors.New(400, "invalid Google token")
	}

	email := payload.Claims["email"].(string)

	// Extract firstName, lastName
	firstName, _ := payload.Claims["given_name"].(string)
	lastName, _ := payload.Claims["family_name"].(string)
	fullName, _ := payload.Claims["name"].(string)

	// Fallback: split full name if given_name / family_name are missing
	if firstName == "" && fullName != "" {
		parts := strings.SplitN(fullName, " ", 2)
		firstName = parts[0]
		if len(parts) > 1 {
			lastName = parts[1]
		}
	}

	var user models.User
	// 2. Check if user exists
	if err := db.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 3. Create new Google user
			user = models.User{
				FirstName:    firstName,
				LastName:     lastName,
				Email:        email,
				IsGoogleUser: true,
			}
			if err := db.DB.WithContext(ctx).Create(&user).Error; err != nil {
				return nil, "", apperrors.Wrap(500, "failed to create google user", err)
			}
		} else {
			return nil, "", apperrors.Wrap(apperrors.ErrFailedToCreateUser, "db error", err)
		}
	}

	// 4. Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // token valid for 24h
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, "", apperrors.Wrap(500, "could not generate token", err)
	}

	return &user, tokenString, nil
}

func (s *userServiceImpl) Login(ctx context.Context, identifier, password string) (*models.LoginResult, *apperrors.AppError) {
	var user models.User

	// 🔍 Try to find user by email OR phone
	if err := db.DB.WithContext(ctx).
		Where("email = ? OR phone = ?", identifier, identifier).
		First(&user).Error; err != nil {
		return nil, apperrors.New(404, "user not found")
	}

	// ✅ Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, apperrors.New(401, "invalid credentials")
	}

	// ✅ Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, apperrors.Wrap(500, "could not generate token", err)
	}

	return &models.LoginResult{
		Token:     token,
		ExpiresIn: 86400, // 24 hours
		User:      user,
	}, nil
}

func (s *userServiceImpl) Logout(ctx context.Context, userID uint) *apperrors.AppError {
	return nil
}

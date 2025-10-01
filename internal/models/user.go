package models

type User struct {
	ID           uint   `gorm:"primaryKey"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `gorm:"unique"`
	Phone        string `gorm:"unique"`
	PasswordHash string
	IsGoogleUser bool  `gorm:"default:false"`
	CreatedAt    int64 `gorm:"autoCreateTime"`
}

type LoginResult struct {
	Token     string
	ExpiresIn int64
	User      User
}

type LoginInput struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

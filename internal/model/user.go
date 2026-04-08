package model

type User struct {
	BaseModel
	Email    string `gorm:"uniqueIndex;not null"`
	Name     string `gorm:"not null"`
	Password string `gorm:"not null"`
}

type UserRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

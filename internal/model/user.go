package model

type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

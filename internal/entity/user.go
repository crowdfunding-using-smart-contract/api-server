package entity

import (
	"fmt"
	"time"

	"fund-o/api-server/pkg/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole int

const (
	Backer UserRole = iota + 1
	Creator
)

type User struct {
	gorm.Model
	Email          string `gorm:"not null;uniqueIndex"`
	HashedPassword string `gorm:"not null"`
	FirstName      string `gorm:"not null"`
	LastName       string `gorm:"not null"`
	ProfileImage   string
	Role           UserRole `gorm:"not null"`
}

type UserDto struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	ProfileImage string `json:"profile_image"`
	Role         string `json:"role"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
} // @name User

// Secondary types

type UserCreatePayload struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role" binding:"required"`
} // @name UserCreatePayload

type UserLoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
} // @name UserLoginPayload

type UserLoginResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiredAt  time.Time `json:"access_token_expired_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`
	User                  *UserDto  `json:"user"`
} // @name UserLoginResponse

// Parse functions

func (u *User) ToUserDto() *UserDto {
	return &UserDto{
		ID:           u.ID,
		Email:        u.Email,
		FullName:     fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		ProfileImage: u.ProfileImage,
		Role:         u.Role.String(),
		CreatedAt:    u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    u.UpdatedAt.Format(time.RFC3339),
	}
}

func (r UserRole) String() string {
	return [...]string{"", "backer", "creator"}[r]
}

var ParseUserRole = func(str string) UserRole {
	mapString := map[string]UserRole{
		"backer":  Backer,
		"creator": Creator,
	}

	role, ok := helper.ParseString[UserRole](mapString, str)
	if !ok {
		return Backer
	}

	return role
}

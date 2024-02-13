package entity

import (
	"fmt"
	"time"

	"fund-o/api-server/pkg/helper"
)

type UserRole int

const (
	Backer UserRole = iota + 1
	Creator
)

type User struct {
	Base
	Email          string `gorm:"not null;uniqueIndex"`
	HashedPassword string `gorm:"not null"`
	Firstname      string `gorm:"not null"`
	Lastname       string `gorm:"not null"`
	PhoneNumber    string `gorm:"not null"`
	ProfileImage   string
}

type UserDto struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	PhoneNumber  string `json:"phone_number"`
	ProfileImage string `json:"profile_image"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
} // @name User

// Secondary types

type UserCreatePayload struct {
	Email                string `json:"email" binding:"required"`
	Firstname            string `json:"firstname" binding:"required"`
	Lastname             string `json:"lastname" binding:"required"`
	PhoneNumber          string `json:"phone_number" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
} // @name UserCreatePayload

type UserLoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
} // @name UserLoginPayload

type UserLoginResponse struct {
	SessionID             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiredAt  time.Time `json:"access_token_expired_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`
	User                  *UserDto  `json:"user"`
} // @name UserLoginResponse

// Parse functions

func (u *User) ToUserDto() *UserDto {
	return &UserDto{
		ID:           u.ID.String(),
		Email:        u.Email,
		FullName:     fmt.Sprintf("%s %s", u.Firstname, u.Lastname),
		PhoneNumber:  u.PhoneNumber,
		ProfileImage: u.ProfileImage,
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

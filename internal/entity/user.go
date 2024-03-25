package entity

import (
	"fmt"
	"mime/multipart"
	"time"

	"fund-o/api-server/pkg/helper"
)

type UserRole int
type Gender int

const (
	Backer UserRole = iota + 1
	Creator
)

const (
	Male Gender = iota + 1
	Female
	NotSay
)

type User struct {
	Base
	Email           string `gorm:"not null;uniqueIndex"`
	HashedPassword  string `gorm:"not null"`
	Firstname       string `gorm:"not null"`
	Lastname        string `gorm:"not null"`
	ProfileImage    string
	BirthDate       time.Time `gorm:"not null"`
	Gender          Gender    `gorm:"not null;default:3"`
	IsEmailVerified bool      `gorm:"not null;default:false"`
}

type UserDto struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	FullName        string `json:"full_name"`
	ProfileImage    string `json:"profile_image"`
	BirthDate       string `json:"birthdate"`
	Gender          string `json:"gender"`
	IsEmailVerified bool   `json:"is_email_verified"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
} // @name User

// Secondary types

type UserCreatePayload struct {
	Email                string `json:"email" binding:"required" example:"someemail@gmail.com"`
	Password             string `json:"password" binding:"required" example:"@Password123"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required" example:"@Password123"`
	Firstname            string `json:"firstname" binding:"required" example:"John"`
	Lastname             string `json:"lastname" binding:"required" example:"Doe"`
	BirthDate            string `json:"birthdate" binding:"required" example:"2002-04-16T00:00:00Z"`
	Gender               string `json:"gender" binding:"required" example:"m"`
} // @name UserCreatePayload

type UserUpdatePayload struct {
	Email           string                `form:"email"`
	ProfileImage    *multipart.FileHeader `form:"profile_image"`
	IsEmailVerified bool                  `form:"is_email_verified"`
} // @name UserUpdatePayload

type UserLoginPayload struct {
	Email    string `json:"email" binding:"required" example:"someemail@gmail.com"`
	Password string `json:"password" binding:"required" example:"@Password123"`
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
		ID:              u.ID.String(),
		Email:           u.Email,
		FullName:        fmt.Sprintf("%s %s", u.Firstname, u.Lastname),
		ProfileImage:    u.ProfileImage,
		BirthDate:       u.BirthDate.Format(time.RFC3339),
		Gender:          u.Gender.String(),
		IsEmailVerified: u.IsEmailVerified,
		CreatedAt:       u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       u.UpdatedAt.Format(time.RFC3339),
	}
}

func (g Gender) String() string {
	return [...]string{"", "m", "f", "ns"}[g]
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

var ParseGender = func(str string) Gender {
	mapString := map[string]Gender{
		"m":  Male,
		"f":  Female,
		"ns": NotSay,
	}

	gender, ok := helper.ParseString(mapString, str)
	if !ok {
		return NotSay
	}

	return gender
}

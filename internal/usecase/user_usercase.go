package usecase

import (
	"fmt"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/password"
	"fund-o/api-server/pkg/uploader"
	"time"

	"github.com/google/uuid"
)

type UserUseCase interface {
	CreateUser(user *entity.UserCreatePayload) (*entity.UserDto, error)
	AuthenticateUser(payload *entity.UserLoginPayload) (*entity.UserDto, error)
	GetUserById(id string) (*entity.UserDto, error)
	GetUserByEmail(email string) (*entity.UserDto, error)
	UpdateUserByID(id string, user *entity.UserUpdatePayload) (*entity.UserDto, error)
}

type userUseCase struct {
	userRepository repository.UserRepository
	imageUploader  uploader.ImageUploader
}

type UserUseCaseOptions struct {
	repository.UserRepository
	uploader.ImageUploader
}

func NewUserUseCase(options *UserUseCaseOptions) UserUseCase {
	return &userUseCase{
		userRepository: options.UserRepository,
		imageUploader:  options.ImageUploader,
	}
}

func (uc *userUseCase) CreateUser(user *entity.UserCreatePayload) (*entity.UserDto, error) {
	if user.Password != user.PasswordConfirmation {
		return nil, fmt.Errorf("password and password confirmation does not match")
	}

	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	birthDate, err := time.Parse(time.RFC3339, user.BirthDate)
	if err != nil {
		return nil, err
	}

	payload := entity.User{
		Email:          user.Email,
		Firstname:      user.Firstname,
		Lastname:       user.Lastname,
		HashedPassword: hashedPassword,
		BirthDate:      birthDate,
		Gender:         entity.ParseGender(user.Gender),
	}

	fmt.Println("payload", payload.Gender)

	newUser, err := uc.userRepository.Create(&payload)
	if err != nil {
		return nil, err
	}

	return newUser.ToUserDto(), nil
}

func (uc *userUseCase) AuthenticateUser(payload *entity.UserLoginPayload) (*entity.UserDto, error) {
	user, err := uc.userRepository.FindByEmail(payload.Email)
	if err != nil {
		return nil, err
	}

	if err := password.CheckPassword(payload.Password, user.HashedPassword); err != nil {
		return nil, err
	}

	return user.ToUserDto(), nil
}

func (uc *userUseCase) GetUserById(id string) (*entity.UserDto, error) {
	userID := uuid.MustParse(id)

	user, err := uc.userRepository.FindById(userID)
	if err != nil {
		return nil, err
	}

	return user.ToUserDto(), nil
}

func (uc *userUseCase) GetUserByEmail(email string) (*entity.UserDto, error) {
	user, err := uc.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return user.ToUserDto(), nil
}

func (uc *userUseCase) UpdateUserByID(id string, user *entity.UserUpdatePayload) (*entity.UserDto, error) {
	userID := uuid.MustParse(id)

	var profileImage string
	if user.ProfileImage != nil {
		imageSource, err := uc.imageUploader.Upload(uploader.ProfileImageFolder, user.ProfileImage)
		if err != nil {
			return nil, err
		}

		profileImage = imageSource
	}

	payload := entity.User{
		Email:           user.Email,
		ProfileImage:    profileImage,
		IsEmailVerified: user.IsEmailVerified,
	}

	updatedUser, err := uc.userRepository.UpdateByID(userID, &payload)
	if err != nil {
		return nil, err
	}

	return updatedUser.ToUserDto(), nil
}

package usecase

import (
	"errors"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/apperrors"
	"fund-o/api-server/pkg/password"
	"fund-o/api-server/pkg/uploader"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserUseCase interface {
	CreateUser(user *entity.User) (*entity.UserDto, error)
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

func (uc *userUseCase) CreateUser(user *entity.User) (*entity.UserDto, error) {
	newUser, err := uc.userRepository.Create(user)
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
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.ErrInvalidUserID
	}

	user, err := uc.userRepository.FindById(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}

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
		Email:             user.Email,
		DisplayName:       user.DisplayName,
		ProfileImage:      profileImage,
		MetaMaskAccountID: user.MetamaskAccountID,
		IsEmailVerified:   user.IsEmailVerified,
	}

	updatedUser, err := uc.userRepository.UpdateByID(userID, &payload)
	if err != nil {
		return nil, err
	}

	return updatedUser.ToUserDto(), nil
}

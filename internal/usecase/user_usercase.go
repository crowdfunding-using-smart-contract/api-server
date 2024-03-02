package usecase

import (
	"fmt"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/password"

	"github.com/google/uuid"
)

type UserUseCase interface {
	CreateUser(user *entity.UserCreatePayload) (*entity.UserDto, error)
	AuthenticateUser(payload *entity.UserLoginPayload) (*entity.UserDto, error)
	GetUserById(id string) (*entity.UserDto, error)
}

type userUseCase struct {
	userRepository repository.UserRepository
}

type UserUseCaseOptions struct {
	UserRepository repository.UserRepository
}

func NewUserUseCase(options *UserUseCaseOptions) UserUseCase {
	return &userUseCase{
		userRepository: options.UserRepository,
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

	payload := entity.User{
		Email:          user.Email,
		Firstname:      user.Firstname,
		Lastname:       user.Lastname,
		PhoneNumber:    user.PhoneNumber,
		HashedPassword: hashedPassword,
	}

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

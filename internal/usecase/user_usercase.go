package usecase

import (
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/password"
)

type UserUsecase interface {
	CreateUser(user *entity.UserCreatePayload) (*entity.UserDto, error)
	AuthenticateUser(payload *entity.UserLoginPayload) (*entity.UserDto, error)
}

type userUsecase struct {
	userRepository repository.UserRepository
}

type UserUsecaseOptions struct {
	UserRepository repository.UserRepository
}

func NewUserUsecase(options *UserUsecaseOptions) UserUsecase {
	return &userUsecase{
		userRepository: options.UserRepository,
	}
}

func (uc *userUsecase) CreateUser(user *entity.UserCreatePayload) (*entity.UserDto, error) {
	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	payload := entity.User{
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		HashedPassword: hashedPassword,
		Role:           entity.ParseUserRole(user.Role),
	}

	newUser, err := uc.userRepository.Create(&payload)
	if err != nil {
		return nil, err
	}

	return newUser.ToUserDto(), nil
}

func (uc *userUsecase) AuthenticateUser(payload *entity.UserLoginPayload) (*entity.UserDto, error) {
	user, err := uc.userRepository.FindByEmail(payload.Email)
	if err != nil {
		return nil, err
	}

	if err := password.CheckPassword(payload.Password, user.HashedPassword); err != nil {
		return nil, err
	}

	return user.ToUserDto(), nil
}

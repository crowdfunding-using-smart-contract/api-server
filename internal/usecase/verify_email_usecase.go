package usecase

import (
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
)

type VerifyEmailUseCase interface {
	CreateVerifyEmail(verifyEmail *entity.VerifyEmailCreatePayload) (*entity.VerifyEmailDto, error)
}

type verifyEmailUseCase struct {
	verifyEmailRepository repository.VerifyEmailRepository
}

type VerifyEmailUseCaseOptions struct {
	repository.VerifyEmailRepository
}

func NewVerifyEmailUseCase(options *VerifyEmailUseCaseOptions) VerifyEmailUseCase {
	return &verifyEmailUseCase{
		verifyEmailRepository: options.VerifyEmailRepository,
	}
}

func (uc *verifyEmailUseCase) CreateVerifyEmail(verifyEmail *entity.VerifyEmailCreatePayload) (*entity.VerifyEmailDto, error) {
	ve, err := uc.verifyEmailRepository.Create(&entity.VerifyEmail{
		Email:      verifyEmail.Email,
		SecretCode: verifyEmail.SecretCode,
	})
	if err != nil {
		return nil, err
	}

	return ve.ToVerifyEmailDto(), nil
}

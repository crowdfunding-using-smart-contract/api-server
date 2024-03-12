package usecase

import (
	"errors"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"time"
)

type VerifyEmailUseCase interface {
	CreateVerifyEmail(verifyEmail *entity.VerifyEmailCreatePayload) (*entity.VerifyEmailDto, error)
	VerifyEmail(payload *entity.VerifyEmailUpdatePayload) (*entity.VerifyEmailDto, error)
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

func (uc *verifyEmailUseCase) VerifyEmail(payload *entity.VerifyEmailUpdatePayload) (*entity.VerifyEmailDto, error) {
	ve, err := uc.verifyEmailRepository.FindByID(payload.ID)
	if err != nil {
		return nil, err
	}

	if ve.SecretCode != payload.SecretCode {
		return nil, errors.New("invalid secret code")
	}

	if ve.IsUsed {
		return nil, errors.New("verify email already used")
	}

	if ve.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("verify email already expired")
	}

	ve.IsUsed = true

	ve, err = uc.verifyEmailRepository.UpdateByID(ve.ID.String(), ve)
	if err != nil {
		return nil, err
	}

	return ve.ToVerifyEmailDto(), nil
}

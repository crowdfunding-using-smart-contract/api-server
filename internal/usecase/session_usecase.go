package usecase

import (
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"

	"github.com/google/uuid"
)

type SessionUsecase interface {
	CreateSession(payload *entity.SessionCreatePayload) (*entity.SessionDto, error)
	GetSessionByID(sessionID uuid.UUID) (*entity.SessionDto, error)
}

type sessionUsecase struct {
	sessionRepository repository.SessionRepository
}

type SessionUsecaseOptions struct {
	repository.SessionRepository
}

func NewSessionUsecase(options *SessionUsecaseOptions) SessionUsecase {
	return &sessionUsecase{
		sessionRepository: options.SessionRepository,
	}
}

func (uc *sessionUsecase) CreateSession(payload *entity.SessionCreatePayload) (*entity.SessionDto, error) {
	session := entity.Session{
		ID:           payload.ID,
		UserID:       uuid.MustParse(payload.UserID),
		RefreshToken: payload.RefreshToken,
		UserAgent:    payload.UserAgent,
		ClientIP:     payload.ClientIP,
		ExpiredAt:    payload.ExpiredAt,
	}

	newSession, err := uc.sessionRepository.Create(&session)
	if err != nil {
		return nil, err
	}

	return newSession.ToSessionDto(), nil
}

func (uc *sessionUsecase) GetSessionByID(sessionID uuid.UUID) (*entity.SessionDto, error) {
	session, err := uc.sessionRepository.FindByID(sessionID)
	if err != nil {
		return nil, err
	}

	return session.ToSessionDto(), nil
}

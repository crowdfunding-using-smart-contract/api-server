package usecase

import (
	"fmt"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
)

type SessionUsecase interface {
	CreateSession(payload *entity.SessionCreatePayload) (*entity.Session, error)
	GetSessionByID(sessionID string) (*entity.Session, error)
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

func (uc *sessionUsecase) CreateSession(payload *entity.SessionCreatePayload) (*entity.Session, error) {
	fmt.Println("payload: ", payload)
	session := entity.Session{
		ID:           payload.ID,
		UserID:       payload.UserID,
		RefreshToken: payload.RefreshToken,
		UserAgent:    payload.UserAgent,
		ClientIP:     payload.ClientIP,
		ExpiredAt:    payload.ExpiredAt,
	}

	newSession, err := uc.sessionRepository.Create(&session)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}

func (uc *sessionUsecase) GetSessionByID(sessionID string) (*entity.Session, error) {
	return uc.sessionRepository.FindByID(sessionID)
}

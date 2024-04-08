package usecase

import (
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/apperrors"
	"github.com/google/uuid"
)

type ChannelUsecase interface {
	CreateChannel(payload *entity.ChannelCreatePayload) (*entity.ChannelDto, error)
	GetExistingChannel(userID string, channelID string) (*entity.ChannelDto, error)
	GetChannelByUserID(userID string) ([]entity.ChannelDto, error)
}

type channelUsecase struct {
	channelRepository repository.ChannelRepository
	userRepository    repository.UserRepository
}

type ChannelUsecaseOptions struct {
	repository.ChannelRepository
}

func NewChannelUsecase(options *ChannelUsecaseOptions) ChannelUsecase {
	return &channelUsecase{
		channelRepository: options.ChannelRepository,
	}
}

func (u *channelUsecase) CreateChannel(payload *entity.ChannelCreatePayload) (*entity.ChannelDto, error) {
	if len(payload.Members) != 2 {
		return nil, apperrors.ErrInvalidMemberChannelLength
	}

	members := []entity.User{
		{
			Base: entity.Base{ID: uuid.MustParse(payload.Members[0])},
		},
		{
			Base: entity.Base{ID: uuid.MustParse(payload.Members[1])},
		},
	}

	channel := entity.Channel{
		Name:     payload.Name,
		Members:  members,
		Messages: []entity.Message{},
	}

	newChannel, err := u.channelRepository.Create(&channel)
	if err != nil {
		return nil, err
	}

	return newChannel.ToChannelDto(), nil
}

func (u *channelUsecase) GetExistingChannel(userID string, channelID string) (*entity.ChannelDto, error) {
	channel, err := u.channelRepository.GetExistingChannel(userID, channelID)
	if err != nil {
		return nil, err
	}

	if channel.ID == uuid.Nil {
		return nil, apperrors.ErrChannelNotFound
	}

	return channel.ToChannelDto(), nil
}

func (u *channelUsecase) GetChannelByUserID(userID string) ([]entity.ChannelDto, error) {
	channels, err := u.channelRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var channelDtos []entity.ChannelDto
	for _, c := range channels {
		channelDtos = append(channelDtos, *c.ToChannelDto())
	}

	return channelDtos, nil
}

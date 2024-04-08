package usecase

import (
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/uploader"
	"github.com/google/uuid"
)

type MessageUsecase interface {
	CreateChannelMessage(channelID uuid.UUID, payload *entity.MessageCreatePayload) (*entity.MessageDto, error)
}

type messageUsecase struct {
	messageRepository repository.MessageRepository
	imageUploader     uploader.ImageUploader
}

type MessageUsecaseOptions struct {
	repository.MessageRepository
	uploader.ImageUploader
}

func NewMessageUsecase(options *MessageUsecaseOptions) MessageUsecase {
	return &messageUsecase{
		messageRepository: options.MessageRepository,
		imageUploader:     options.ImageUploader,
	}
}

func (u *messageUsecase) CreateChannelMessage(channelID uuid.UUID, payload *entity.MessageCreatePayload) (*entity.MessageDto, error) {
	var attachment *string
	if payload.Attachment != nil {
		attachmentURL, err := u.imageUploader.Upload(uploader.PostImageFolder, payload.Attachment)
		if err != nil {
			return nil, err
		}

		attachment = &attachmentURL
	}

	message := entity.Message{
		Text:       payload.Text,
		Attachment: attachment,
		ChannelID:  channelID,
		AuthorID:   payload.AuthorID,
	}

	newMessage, err := u.messageRepository.Create(&message)
	if err != nil {
		return nil, err
	}

	return newMessage.ToMessageDto(), nil
}

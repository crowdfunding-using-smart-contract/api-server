package entity

import (
	"github.com/google/uuid"
	"mime/multipart"
	"time"
)

type Message struct {
	Base
	Text       *string `gorm:"varchar(255)"`
	Attachment *string `gorm:"varchar(255)"`
	ChannelID  uuid.UUID
	AuthorID   uuid.UUID
	Author     User `gorm:"foreignKey:AuthorID"`
}

type MessageDto struct {
	Text       *string  `json:"text"`
	Attachment *string  `json:"attachment"`
	Author     *UserDto `json:"author"`
	CreatedAt  string   `json:"created_at"`
}

// Secondary types

type MessageCreatePayload struct {
	Text       *string               `form:"text"`
	Attachment *multipart.FileHeader `form:"attachment"`
	ChannelID  uuid.UUID             `form:"-"`
	AuthorID   uuid.UUID             `form:"-"`
}

// Parse functions

func (m *Message) ToMessageDto() *MessageDto {
	return &MessageDto{
		Text:       m.Text,
		Attachment: m.Attachment,
		Author:     m.Author.ToUserDto(),
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
	}
}

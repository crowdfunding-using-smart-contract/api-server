package entity

type Channel struct {
	Base
	Name     string    `gorm:"type:varchar(255);not null"`
	Messages []Message `gorm:"foreignKey:ChannelID"`
	Members  []User    `gorm:"many2many:channel_members;"`
}

type ChannelDto struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Messages []MessageDto `json:"messages"`
	Members  []UserDto    `json:"members"`
}

// Secondary types

type GetOwnChannelsResponse struct {
	Receiver    UserDto     `json:"receiver"`
	LastMessage *MessageDto `json:"last_message"`
}

type ChannelCreatePayload struct {
	Name    string `json:"name"`
	Members []string
}

// Parse functions

func (c *Channel) ToChannelDto() *ChannelDto {
	var messages []MessageDto
	for _, m := range c.Messages {
		messages = append(messages, *m.ToMessageDto())
	}

	var members []UserDto
	for _, u := range c.Members {
		members = append(members, *u.ToUserDto())
	}

	return &ChannelDto{
		ID:       c.ID.String(),
		Name:     c.Name,
		Messages: messages,
		Members:  members,
	}
}

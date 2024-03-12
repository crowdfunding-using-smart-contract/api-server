package entity

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	Base
	Content  string `gorm:"type:varchar(255);not null"`
	AuthorID uuid.UUID
	Author   User `gorm:"foreignKey:AuthorID"`
	ForumID  uuid.UUID
	Reply    []Reply `gorm:"foreignKey:CommentID"`
}

type Reply struct {
	Base
	Content   string `gorm:"type:varchar(255);not null"`
	AuthorID  uuid.UUID
	Author    User `gorm:"foreignKey:AuthorID"`
	CommentID uuid.UUID
}

type CommentDto struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	Author    *UserDto   `json:"author"`
	Replies   []ReplyDto `json:"replies"`
	CreatedAt string     `json:"created_at"`
}

type ReplyDto struct {
	ID        string   `json:"id"`
	Content   string   `json:"content"`
	Author    *UserDto `json:"author"`
	CreatedAt string   `json:"created_at"`
}

// Secondary types

type CommentCreatePayload struct {
	Content  string `json:"content" binding:"required"`
	AuthorID string `swaggerignore:"true"`
}

// Parse functions

func (c *Comment) ToCommentDto() *CommentDto {
	replies := make([]ReplyDto, len(c.Reply))
	for i, reply := range c.Reply {
		replies[i] = *reply.ToReplyDto()
	}

	return &CommentDto{
		ID:        c.ID.String(),
		Content:   c.Content,
		Author:    c.Author.ToUserDto(),
		Replies:   replies,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
	}
}

func (r *Reply) ToReplyDto() *ReplyDto {
	return &ReplyDto{
		ID:        r.ID.String(),
		Content:   r.Content,
		Author:    r.Author.ToUserDto(),
		CreatedAt: r.CreatedAt.Format(time.RFC3339),
	}
}

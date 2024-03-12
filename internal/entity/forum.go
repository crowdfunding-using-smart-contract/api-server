package entity

import (
	"github.com/google/uuid"
	"time"
)

type Forum struct {
	Base
	Title    string    `gorm:"size:255;not null"`
	Content  string    `gorm:"not null"`
	AuthorID uuid.UUID `gorm:"not null"`
	Author   User      `gorm:"foreignKey:AuthorID"`
	Comments []Comment `gorm:"foreignKey:ForumID"`
}

type ForumDto struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Author    *UserDto     `json:"author"`
	CreatedAt string       `json:"created_at"`
	Comments  []CommentDto `json:"comments"`
} // @name Forum

// Secondary types

type ForumCreatePayload struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	AuthorID string `swaggerignore:"true"`
}

// Parse functions

func (f *Forum) ToForumDto() *ForumDto {
	comments := make([]CommentDto, len(f.Comments))
	for i, comment := range f.Comments {
		comments[i] = *comment.ToCommentDto()
	}

	return &ForumDto{
		ID:        f.ID.String(),
		Title:     f.Title,
		Content:   f.Content,
		Author:    f.Author.ToUserDto(),
		Comments:  comments,
		CreatedAt: f.CreatedAt.Format(time.RFC3339),
	}
}

package entity

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	Base
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:varchar(255);not null"`
	Content     string    `gorm:"not null"`
	AuthorID    uuid.UUID `gorm:"not null"`
	Author      User      `gorm:"foreignKey:AuthorID"`
	ProjectID   uuid.UUID `gorm:"not null"`
	Project     Project   `gorm:"foreignKey:ProjectID"`
	Comments    []Comment
}

type PostDto struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Content     string       `json:"content"`
	Author      *UserDto     `json:"author"`
	Project     *ProjectDto  `json:"project"`
	Comments    []CommentDto `json:"comments"`
	CreatedAt   string       `json:"created_at"`
} // @name Post

type Comment struct {
	Base
	Content  string `gorm:"type:varchar(255);not null"`
	AuthorID uuid.UUID
	Author   User `gorm:"foreignKey:AuthorID"`
	PostID   uuid.UUID
	Replies  []Reply
}

type CommentDto struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	Author    *UserDto   `json:"author"`
	Replies   []ReplyDto `json:"replies"`
	CreatedAt string     `json:"created_at"`
} // @name Comment

type Reply struct {
	Base
	Content   string `gorm:"type:varchar(255);not null"`
	AuthorID  uuid.UUID
	Author    User `gorm:"foreignKey:AuthorID"`
	CommentID uuid.UUID
}

type ReplyDto struct {
	ID        string   `json:"id"`
	Content   string   `json:"content"`
	Author    *UserDto `json:"author"`
	CreatedAt string   `json:"created_at"`
} // @name Reply

// Secondary types

type PostCreatePayload struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Content     string `json:"content"`
	ProjectID   string `json:"project_id" binding:"required"`
	AuthorID    string `swaggerignore:"true"`
}

type CommentCreatePayload struct {
	Content  string `json:"content" binding:"required"`
	AuthorID string `swaggerignore:"true"`
}

type ReplyCreatePayload struct {
	Content  string `json:"content" binding:"required"`
	AuthorID string `swaggerignore:"true"`
}

// Parse functions

func (f *Post) ToPostDto() *PostDto {
	comments := make([]CommentDto, len(f.Comments))
	for i, comment := range f.Comments {
		comments[i] = *comment.ToCommentDto()
	}

	return &PostDto{
		ID:          f.ID.String(),
		Title:       f.Title,
		Description: f.Description,
		Content:     f.Content,
		Author:      f.Author.ToUserDto(),
		Project:     f.Project.ToProjectDto(),
		Comments:    comments,
		CreatedAt:   f.CreatedAt.Format(time.RFC3339),
	}
}

func (c *Comment) ToCommentDto() *CommentDto {
	replies := make([]ReplyDto, len(c.Replies))
	for i, reply := range c.Replies {
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

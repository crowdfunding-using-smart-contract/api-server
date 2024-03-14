package entity

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	Base
	Title    string    `gorm:"size:255;not null"`
	Content  string    `gorm:"not null"`
	AuthorID uuid.UUID `gorm:"not null"`
	Author   User      `gorm:"foreignKey:AuthorID"`
	Comments []Comment `gorm:"foreignKey:ForumID"`
}

type PostDto struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Author    *UserDto     `json:"author"`
	CreatedAt string       `json:"created_at"`
	Comments  []CommentDto `json:"comments"`
} // @name Post

type Comment struct {
	Base
	Content  string `gorm:"type:varchar(255);not null"`
	AuthorID uuid.UUID
	Author   User `gorm:"foreignKey:AuthorID"`
	ForumID  uuid.UUID
	Reply    []Reply `gorm:"foreignKey:CommentID"`
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
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	AuthorID string `swaggerignore:"true"`
}

type CommentCreatePayload struct {
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
		ID:        f.ID.String(),
		Title:     f.Title,
		Content:   f.Content,
		Author:    f.Author.ToUserDto(),
		Comments:  comments,
		CreatedAt: f.CreatedAt.Format(time.RFC3339),
	}
}

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

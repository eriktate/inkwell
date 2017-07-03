package inkwell

import "time"

// A Blog represents all of the fields that make up an Inkwell blog.
type Blog struct {
	ID              string     `json:"id" dynamodbav:"blog_id"`
	Title           string     `json:"title" dynamodbav:"title""`
	AuthorID        string     `json:"authorId" dynamodbav:"authorID"`
	Content         string     `json:"content" dynamodbav:"content"`
	ContentLocation string     `json:"contentLocation" dynamodbav:"content_loc"`
	Comments        []Comment  `json:"comments" dynamodbav:"comments"`
	Published       bool       `json:"published" dynamodbav:"published"`
	CreatedAt       time.Time  `json:"createdAt" dynamodbav:"created_at"`
	UpdatedAt       time.Time  `json:"updatedAt" dynamodbav:"updated_at"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty" dynamodbav:"deleted_at"`
}

// A Comment represents all of the fields that make up an Inkwell blog comment.
type Comment struct {
	ID         string    `json:"id"`
	AuthorName string    `json:"authorName"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	DeletedAt  time.Time `json:"deletedAt,omitempty"`
}

// BlogService describes all of the actions necessary for an implementation to be
// considered a valid backend for dealing with blogs.
type BlogService interface {
	Get(blogID string) (Blog, error)
	// GetByAuthor(authorID string) ([]Blog, error)
	Write(blog Blog) error
	// Update(blog Blog) error
	// Publish(blogID string) error
	// Delete(blogID string) error
}

// Client wraps concrete implementations of all required inkwell services.
type Client interface {
	BlogService() BlogService
}

package inkwell

import "time"

// A Blog represents all of the fields that make up an Inkwell blog.
type Blog struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	AuthorID        string     `json:"authorId"`
	Content         string     `json:"content"`
	ContentLocation string     `json:"contentLocation"`
	Comments        []Comment  `json:"comments"`
	Published       bool       `json:"published"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty"`
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

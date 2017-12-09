package inkwell

import "time"

// Blog
type Blog struct {
	AuthorID string
	Key      string
	Title    string
	// Content may be something else longer term.
	Content   []byte
	Published bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// A BlogWriter knows how to persist Blog data.
type BlogWriter interface {
	Write(blog Blog) error
	SetContent(authorID, key string, content []byte) error
	SetKey(authorID, key, newKey string) error
	SetTitle(authorID, key, title string) error
	Publish(authorID, key string) error
	Redact(authorID, key string) error
	Delete(authorID, key string) error
}

// A BlogReader knows how to fetch Blog data.
type BlogReader interface {
	Get(authorID, key string) (Blog, error)
	GetBlogsByAuthor(authorID string) ([]Blog, error)
}

// A BlogReadWriter knows how to fetch and persist Blog data.
type BlogReadWriter interface {
	BlogWriter
	BlogReader
}

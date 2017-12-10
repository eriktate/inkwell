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

// An Author represents someone who writes blogs within inkwell.
type Author struct {
	AuthorID string
	// TODO: Flesh out author more
}

// An AuthorReader knows how to fetch Author data.
type AuthorReader interface {
	Get(authorID string) (Author, error)
}

// An AuthorWriter knows how to persist Author data.
type AuthorWriter interface {
	Write(author Author) error
}

// An AuthorReadWriter knows how to fetch and persist Author data.
type AuthorReadWriter interface {
	AuthorReader
	AuthorWriter
}

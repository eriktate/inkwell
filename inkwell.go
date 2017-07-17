package inkwell

import (
	"fmt"
	"time"

	"github.com/eriktate/inkwell/html"
)

// A Blog represents all of the fields that make up an Inkwell blog.
type Blog struct {
	ID        string     `json:"id" dynamodbav:"blog_id"`
	Title     string     `json:"title" dynamodbav:"title""`
	AuthorID  string     `json:"authorId" dynamodbav:"authorID"`
	Content   string     `json:"content" dynamodbav:"-"`
	Comments  []Comment  `json:"comments" dynamodbav:"comments"`
	Published bool       `json:"published" dynamodbav:"published"`
	CreatedAt time.Time  `json:"createdAt" dynamodbav:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" dynamodbav:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" dynamodbav:"deleted_at"`
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
	Get(authorID, blogID string) (Blog, error)
	// GetByAuthor(authorID string) ([]Blog, error)
	Write(blog Blog) error
	Revise(authorID, blogID, content string) error
	Publish(authorID, blogID string) error
	Redact(authorID, blogID string) error
	// Delete(blogID string) error
}

type Post struct {
	ID   string
	Rows RowList
}

type Row struct {
	ID      string   `json:"id" toml:"id"`
	Columns []Column `json:"columns" toml:"columns"`
}

type RowList []Row

type Column struct {
	ID       string `json:"id" toml:"id"`
	Contents Renderer
}

type Renderer interface {
	Render() string
}

func (r Row) Render() string {
	contents := ""
	for _, col := range r.Columns {
		contents = fmt.Sprintf("%s%s", contents, RenderColumn(col))
	}

	return contents
}

func RenderColumn(col Column) string {
	idAttr := html.NewAttribute("id", col.ID)
	classAttr := html.NewAttribute("class", "col")
	return html.WrapDiv(col.Contents.Render(), idAttr, classAttr)
}

func RenderPost(post Post) string {
	idAttr := html.NewAttribute("id", post.ID)
	classAttr := html.NewAttribute("class", "post")
	return html.WrapDiv(post.Rows.Render(), idAttr, classAttr)
}

func (rl RowList) Render() string {
	contents := ""
	for _, row := range rl {
		contents = fmt.Sprintf("%s%s", contents, row.Render())
	}

	return contents
}

// Client wraps concrete implementations of all required inkwell services.
type Client interface {
	BlogService() BlogService
}

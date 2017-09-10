package inkwell

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/eriktate/inkwell/html"
)

// A Blog represents all of the fields that make up an Inkwell blog.
type Blog struct {
	ID        string     `json:"id" dynamodbav:"blog_id"`
	Title     string     `json:"title" dynamodbav:"title"`
	AuthorID  string     `json:"authorId" dynamodbav:"author_id"`
	Content   string     `json:"content" dynamodbav:"-"`
	Comments  []Comment  `json:"comments,omitempty" dynamodbav:"comments"`
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

// A Post represents the top level structure of a blog post.
type Post struct {
	ID    string  `json:"id" toml:"id"`
	Title string  `json:"title" toml:"title"`
	Rows  RowList `json:"rows" toml:"row"`
}

// A Row represents a single row of content in a blog post.
type Row struct {
	ID      string   `json:"id" toml:"id"`
	Columns []Column `json:"columns" toml:"column"`
}

// A RowList represents the list of Rows that ultimately make up a Blog post.
type RowList []Row

// A Column represents an individual column of content within a Row.
type Column struct {
	ID       string      `json:"id" toml:"id"`
	Contents ContentList `json:"contents" toml:"contents"`
}

// A Content represents something that can be rendered in a column. This can be individual
// elements or even another RowList.
type Content struct {
	ID         string          `json:"id" toml:"id"`
	Type       string          `json:"type" toml:"type"`
	RawElement json.RawMessage `json:"element"`
	Element    Renderer
}

// An ContentList aliases a Content slice so that we can implement the Renderer interface.
type ContentList []Content

// A Renderer represents anything capable of rendering itself.
type Renderer interface {
	Render() string
}

// Render implements the Renderer interface and converts a Row into raw HTML.
func (r Row) Render() string {
	contents := ""
	for _, col := range r.Columns {
		contents = fmt.Sprintf("%s%s", contents, RenderColumn(col))
	}

	classAttr := html.NewAttribute("class", "row")

	return html.WrapDiv(contents, classAttr)
}

// RenderColumn specifically prevents Columns from being considered a Renderer, but allows them to still be rendered.
func RenderColumn(col Column) string {
	idAttr := html.NewAttribute("id", col.ID)
	classAttr := html.NewAttribute("class", "col")
	return html.WrapDiv(col.Contents.Render(), idAttr, classAttr)
}

// RenderPost specifically prevents Posts from being considered a Renderer, but allows them to still be rendered.
func RenderPost(post Post) string {
	idAttr := html.NewAttribute("id", post.ID)
	classAttr := html.NewAttribute("class", "post")
	return html.WrapDiv(post.Rows.Render(), idAttr, classAttr)
}

// Render implements the Renderer interface and converts a RowList into raw HTML.
func (rl RowList) Render() string {
	var result string
	for _, row := range rl {
		result = fmt.Sprintf("%s%s", result, row.Render())
	}

	return result
}

// Render implements the Renderer interface and converts an ContentList into raw HTML.
func (cl ContentList) Render() string {
	var result string
	for _, content := range cl {
		result = fmt.Sprintf("%s%s", result, content.Render())
	}

	return result
}

// Render implements the Renderer interface and converts a Content struct into raw HTML.
func (c Content) Render() string {
	return c.Element.Render()
}

// BuildPost takes a byte slice of JSON and prepares a Post to be rendered.
func BuildPost(data []byte) (Post, error) {
	var post Post
	if err := json.Unmarshal(data, &post); err != nil {
		return post, err
	}

	for rowIDX, row := range post.Rows {
		for colIDX, col := range row.Columns {
			for conIDX, c := range col.Contents {
				switch c.Type {
				case "paragraph":
					p, err := BuildParagraph(c.RawElement)
					if err != nil {
						return post, err
					}
					// TODO: This is gross. Maybe use pointers instead.
					post.Rows[rowIDX].Columns[colIDX].Contents[conIDX].Element = p
				}
			}
		}
	}

	return post, nil
}

// Client wraps concrete implementations of all required inkwell services.
type Client interface {
	BlogService() BlogService
}

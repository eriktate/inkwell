package inkwell

import (
	"encoding/json"

	"github.com/eriktate/inkwell/html"
)

// Paragraph represents a block of text content in a blog.
type Paragraph struct {
	ID   string `json:"id" toml:"id"`
	Text string `json:"text" toml:"text"`
}

// Render returns the necessary DOm to represent a paragraph.
func (p Paragraph) Render() string {
	idAttr := html.NewAttribute("id", p.ID)
	// TODO: Render markdown here.
	return html.WrapParagraphUnsafe(p.Text, idAttr)
}

// Image represents static image content in a blog.
type Image struct {
	ID  string `json:"id" toml:"id"`
	Src string `json:"src" toml:"src"`
	URL string `json:"url" toml:"url"`
	Alt string `json:"alt" toml:"alt"`
}

// Render returns the necessary DOM to represent an Image.
func (i Image) Render() string {
	idAttr := html.NewAttribute("id", i.ID)
	altAttr := html.NewAttribute("alt", i.Alt)
	dom := html.MakeImage(i.Src)
	if len(i.URL) > 0 {
		return html.WrapAnchorUnsafe(dom, i.URL, idAttr, altAttr)
	}

	return dom
}

// BuildParagraph takes a bytee slicee of JSON and prepares a Paragraph to be rendered.
func BuildParagraph(data []byte) (Paragraph, error) {
	var p Paragraph
	err := json.Unmarshal(data, &p)
	return p, err
}

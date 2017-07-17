package inkwell

import "github.com/eriktate/inkwell/html"

// Paragrpah represents a block of text content in a blog.
type Paragraph struct {
	ID   string `json:"id" toml:"id"`
	Text string `json:"text" toml:"text"`
}

func (p Paragraph) Render() string {
	idAttr := NewAttribute("id", p.ID)
	// TODO: Render markdown here.
	return html.WrapParagraphUnsafe(p.Text, idAttr)
}

type Image struct {
	ID  string `json:"id" toml:"id"`
	Src string `json:"src" toml:"src"`
	URL string `json:"url" toml:"url"`
	Alt string `json:"alt" toml:"alt"`
}

func (i Image) Render() string {
	idAttr := html.NewAttribute("id", i.ID)
	altAttr := html.NewAttribute("alt", i.Alt)
	dom := html.MakeImage(src)
	if len(i.URL) > 0 {
		return html.WrapAnchorUnsafe(dom, i.URL, idAttr)
	}

	return dom
}

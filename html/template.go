package html

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

// Page represents the basic structure allowed when building out a blog.
type Page struct {
	Title   string
	Header  string
	Post    string
	SideBar string
	Footer  string
}

// GeneratePage returns a full HTML document given a base template and a
// page struct to fill it with.
func GeneratePage(base string, page Page) ([]byte, error) {
	// TODO: This feels pretty arbitrary. Should allocate more sensical space.
	buf := bytes.NewBuffer(make([]byte, 0, len(base)))
	tmpl, err := template.New("page").Parse(base)
	if err != nil {
		return nil, errors.Wrap(err, "could not build template.")
	}

	if err := tmpl.Execute(buf, page); err != nil {
		return nil, errors.Wrap(err, "failed to execute page template")
	}

	return buf.Bytes(), nil
}

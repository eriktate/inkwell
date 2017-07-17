package html

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

type Page struct {
	Header  []byte
	Post    []byte
	SideBar []byte
	Footer  []byte
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

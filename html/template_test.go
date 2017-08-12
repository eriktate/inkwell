package html_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/eriktate/inkwell/html"
)

func Test_GeneratePage(t *testing.T) {
	// SETUP
	base, err := ioutil.ReadFile("../template.html")
	if err != nil {
		t.Fatalf("Failed to read template file: %s", err)
	}

	page := html.Page{
		Title:   "Some page title",
		Header:  "<h1>Hello there! This is a header!</h1>",
		Post:    "<p>This is the actual post content!</p>",
		SideBar: "<aside>This is some sidebar thing</aside>",
		Footer:  "<div>This is where all the footer things go.</div>",
	}

	// RUN
	generated, err := html.GeneratePage(string(base), page)

	// ASSERT
	if err != nil {
		t.Fatalf("Failed to generate page: %s", err)
	}

	log.Printf("Actuallly generated a page:\n%s", string(generated))
}

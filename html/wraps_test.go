package html_test

import (
	"log"
	"testing"

	"github.com/eriktate/inkwell/html"
)

func Test_WrapDiv(t *testing.T) {
	// SETUP
	contents := "Hello, world!"
	id := html.NewAttribute("id", "12345")
	name := html.NewAttribute("name", "someDiv")

	// RUN
	dom := html.WrapDiv(contents, id, name)
	log.Printf("DOM: %s", dom)

	// ASSERT
	if dom != "<div id=\"12345\" name=\"someDiv\">Hello, world!</div>" {
		t.Fatal("WrapDiv did not produce the expected DOM.")
	}
}

func Test_WrapParagraph_SafeInput(t *testing.T) {
	// SETUP
	contents := "Hello, world!"
	id := html.NewAttribute("id", "12345")
	class := html.NewAttribute("class", "some-class")

	// RUN
	dom := html.WrapParagraph(contents, id, class)
	log.Printf("DOM: %s", dom)

	// ASSERT
	if dom != "<p id=\"12345\" class=\"some-class\">Hello, world!</p>" {
		t.Fatal("WrapParagraph did not produce the expected DOM.")
	}
}

func Test_WrapParagraph_UnsafeInput(t *testing.T) {
	// SETUP
	contents := `<script>alert("This should not be allowed.")</script>`
	id := html.NewAttribute("id", "12345")

	// RUN
	dom := html.WrapParagraph(contents, id)
	log.Printf("DOM: %s", dom)

	// ASSERT
	if dom != "<p id=\"12345\">&lt;script&gt;alert(&#34;This should not be allowed.&#34;)&lt;/script&gt;</p>" {
		t.Fatal("WrapParagraph did not produce the expected DOM.")
	}
}

func Test_WrapParagraphUnsafe(t *testing.T) {
	// SETUP
	contents := `<strong>the boldest of text</strong>`
	id := html.NewAttribute("id", "12345")

	// RUN
	dom := html.WrapParagraphUnsafe(contents, id)
	log.Printf("DOM: %s", dom)

	// ASSERT
	if dom != "<p id=\"12345\"><strong>the boldest of text</strong></p>" {
		t.Fatal("WrapParagraphUnsafe did not produce the expected DOM.")
	}
}

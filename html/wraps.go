package html

import (
	"fmt"
	"html"
)

type Attribute struct {
	Key string
	Val string
}

func NewAttribute(key, val string) Attribute {
	return Attribute{
		Key: html.EscapeString(key),
		Val: html.EscapeString(val),
	}
}

func WrapDiv(contents string, attrs ...Attribute) string {
	return fmt.Sprintf("<div%s>%s</div>", stringifyAttributes(attrs), contents)
}

// WrapParagraph will wrap the provided contents in p tags after first escaping.
func WrapParagraph(contents string, attrs ...Attribute) string {
	return fmt.Sprintf("<p%s>%s</p>", stringifyAttributes(attrs), html.EscapeString(contents))
}

// WrapParagraphUnafe will wrap the provided contents in p tags without escaping.
func WrapParagraphUnsafe(contents string, attrs ...Attribute) string {
	return fmt.Sprintf("<p%s>%s</p>", stringifyAttributes(attrs), contents)
}

func WrapHeading(contents string, weight int, attrs ...Attribute) string {
	return fmt.Sprintf("<h%d%s>%s</h%d", weight, stringifyAttributes(attrs), contents, weight)
}

func WrapAnchor(contents, href string, attrs ...Attribute) string {
	return fmt.Sprintf("<a href=\"%s\"%s>%s</a>", html.EscapeString(href), stringifyAttributes(attrs), html.EscapeString(contents))
}

func WrapAnchorUnsafe(contents, href string, attrs ...Attribute) string {
	return fmt.Sprintf("<a href=\"%s\"%s>%s</a>", html.EscapeString(href), stringifyAttributes(attrs), contents)
}

func MakeImage(src string, attrs ...Attribute) string {
	return fmt.Sprintf("<img src=\"%s\"%s />", src, stringifyAttributes(attrs))
}

func stringifyAttributes(attrs []Attribute) string {
	var result string

	for _, attr := range attrs {
		result = fmt.Sprintf("%s %s=\"%s\"", result, attr.Key, attr.Val)
	}

	return result
}

func EscapeString(s string) string {
	return html.EscapeString(s)
}

func UnescapeString(s string) string {
	return html.UnescapeString(s)
}

package html

import (
	"fmt"
	"html"
)

// An Attribute is an HTML key/val attribute.
type Attribute struct {
	Key string
	Val string
}

// NewAttribute returns a new Attribute given a key/val combination.
func NewAttribute(key, val string) Attribute {
	return Attribute{
		Key: html.EscapeString(key),
		Val: html.EscapeString(val),
	}
}

// WrapDiv wraps some content in div tags.
func WrapDiv(contents string, attrs ...Attribute) string {
	return fmt.Sprintf("<div%s>%s</div>", stringifyAttributes(attrs), contents)
}

// WrapParagraph will wrap the provided contents in p tags after first escaping.
func WrapParagraph(contents string, attrs ...Attribute) string {
	return fmt.Sprintf("<p%s>%s</p>", stringifyAttributes(attrs), html.EscapeString(contents))
}

// WrapParagraphUnsafe will wrap the provided contents in p tags without escaping.
func WrapParagraphUnsafe(contents string, attrs ...Attribute) string {
	return fmt.Sprintf("<p%s>%s</p>", stringifyAttributes(attrs), contents)
}

// WrapHeading will wrap the provided contents in h* tags given a weight from 1-6.
func WrapHeading(contents string, weight int, attrs ...Attribute) string {
	return fmt.Sprintf("<h%d%s>%s</h%d", weight, stringifyAttributes(attrs), contents, weight)
}

// WrapAnchor will wrap the provided href in a tags.
func WrapAnchor(contents, href string, attrs ...Attribute) string {
	return fmt.Sprintf("<a href=\"%s\"%s>%s</a>", html.EscapeString(href), stringifyAttributes(attrs), html.EscapeString(contents))
}

// WrapAnchorUnsafe will wrap the provided href in a tags without escaping.
func WrapAnchorUnsafe(contents, href string, attrs ...Attribute) string {
	return fmt.Sprintf("<a href=\"%s\"%s>%s</a>", html.EscapeString(href), stringifyAttributes(attrs), contents)
}

// MakeImage returns an image tag given a src and some attributes.
func MakeImage(src string, attrs ...Attribute) string {
	return fmt.Sprintf("<img src=\"%s\"%s />", src, stringifyAttributes(attrs))
}

// EscapeString HTML escapes the given string.
func EscapeString(s string) string {
	return html.EscapeString(s)
}

// UnescapeString reverts any HTML escaping present on the given string.
func UnescapeString(s string) string {
	return html.UnescapeString(s)
}

func stringifyAttributes(attrs []Attribute) string {
	var result string

	for _, attr := range attrs {
		result = fmt.Sprintf("%s %s=\"%s\"", result, attr.Key, attr.Val)
	}

	return result
}

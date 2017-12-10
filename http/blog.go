package http

import (
	"encoding/json"
	"net/http"

	"github.com/eriktate/inkwell"
)

type BlogHandler struct {
	svc inkwell.BlogWriter
}

// Handle handles requests hanging off of the '/blog' route. It needs an
// authorID to do anything.
func (h BlogHandler) Handle(authorID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var route string
		route, r.URL.Path = ShiftPath(r.URL.Path)

		switch route {
		case "":
			h.handlePost(authorID)
		default:
			// Route must be a blog key.
			h.handleBlogKey(authorID, route)
		}
		if route != "" {
			notFound(w, "Endpoint does not exist")
			return
		}

		if r.Method != "POST" {
			badRequest(w, "Method not supported for this route")
		}
	})
}

func (h BlogHandler) handleBlogKey(authorID, key string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var route string
		route, r.URL.Path = ShiftPath(r.URL.Path)

		switch route {
		case "":
			// TODO: Fetch a specific blog
		case "publish":
			// TODO: Publish the blog
		case "redact":
			// TODO: Redact the blog.
		case "content":
			// TODO: SetContent of blog
		case "title":
			// TODO: SetTitle of blog
		default:
			badRequest(w, "Endpoint does not exist")
		}
	})
}

func (h BlogHandler) handlePost(authorID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Grab blog struct
		var blog inkwell.Blog

		if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
			badRequest(w, "JSON payload could not be reconciled")
			return
		}

		// Make sure we get the authorID in the route. It's not explicitly
		// required client side.
		blog.AuthorID = authorID
		if err := h.svc.Write(blog); err != nil {
			serverError(w, "Something went wrong while write blog")
			return
		}

		noContent(w)
	})
}

package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/eriktate/inkwell"
	"github.com/sirupsen/logrus"
)

// BlogHandler handles requests hanging off of the 'blog' resource.
type BlogHandler struct {
	svc inkwell.BlogReadWriter
	log *logrus.Logger
}

func NewBlogHandler(svc inkwell.BlogReadWriter, log *logrus.Logger) BlogHandler {
	return BlogHandler{svc: svc, log: log}
}

// Handle handles requests hanging off of the '/blog' route. It needs an
// authorID to do anything.
func (h BlogHandler) Handle(authorID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var route string
		route, r.URL.Path = ShiftPath(r.URL.Path)

		switch route {
		case "":
			h.handlePost(authorID).ServeHTTP(w, r)
		default:
			// Route must be a blog key.
			h.handleBlogKey(authorID, route).ServeHTTP(w, r)
		}

		if route != "" {
			notFound(w, "Endpoint does not exist")
			return
		}

	})
}

func (h BlogHandler) handleBlogKey(authorID, key string) http.Handler {
	// Setup logger to automatically capture important data.
	log := h.log.WithFields(logrus.Fields{
		"authorID": authorID,
		"key":      key,
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var route string
		route, r.URL.Path = ShiftPath(r.URL.Path)

		// We can handle all of the method checking before the switch
		if route == "" {
			if r.Method != http.MethodGet {
				badRequest(w, "Method not supported")
				return
			}
		} else {
			if r.Method != http.MethodPost {
				badRequest(w, "Method not supported")
				return
			}
		}

		switch route {
		case "":
			blog, err := h.svc.Get(authorID, key)
			if err != nil {
				log.WithError(err).Error("Failed to Get blog")
				serverError(w, "Something went wrong while fetching blog")
				return
			}

			// TODO: Handle not founds?
			data, err := json.Marshal(&blog)
			if err != nil {
				log.WithError(err).Error("Failed to marshal blog")
				serverError(w, "Something went wrong when marshaling blog to JSON")
				return
			}

			ok(w, data)
			return
		case "publish":
			if err := h.svc.Publish(authorID, key); err != nil {
				log.WithError(err).Error("Failed to Publish blog")
				serverError(w, "Something went wrong while publishing blog")
				return
			}

			noContent(w)
			return
		case "redact":
			if err := h.svc.Redact(authorID, key); err != nil {
				log.WithError(err).Error("Failed to Redact blog")
				serverError(w, "Something went wrong while redacting blog")
				return
			}

			noContent(w)
			return
		case "content":
			content, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				log.WithError(err).Error("Failed to read request body")
				badRequest(w, "Could not read request body")
				return
			}

			if err := h.svc.SetContent(authorID, key, content); err != nil {
				serverError(w, "Something went wrong while setting blog content")
				return
			}

			noContent(w)
			return
		case "title":
			title, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				badRequest(w, "Could not read request body")
				return
			}

			if err := h.svc.SetTitle(authorID, key, string(title)); err != nil {
				log.WithError(err).Error("Failed to SetTitle")
				serverError(w, "Something went wrong while setting blog title")
				return
			}

			noContent(w)
			return
		default:
			badRequest(w, "Endpoint does not exist")
		}
	})
}

func (h BlogHandler) handlePost(authorID string) http.Handler {
	log := h.log.WithField("authorID", authorID)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			badRequest(w, "Method not supported for this route")
			return
		}

		defer r.Body.Close()

		// Grab blog struct
		var blog inkwell.Blog

		if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
			log.WithError(err).Error("Failed to decode body")
			badRequest(w, "JSON payload could not be reconciled")
			return
		}

		// Make sure we get the authorID in the route. It's not explicitly
		// required client side.
		blog.AuthorID = authorID
		if err := h.svc.Write(blog); err != nil {
			log.WithError(err).Error("Failed to Write blog")
			serverError(w, "Something went wrong while write blog")
			return
		}

		noContent(w)
	})
}

package http

import (
	"encoding/json"
	"net/http"

	"github.com/eriktate/inkwell"
)

// AuthorHandler handles requests to the /author
// base route.
type AuthorHandler struct {
	svc         inkwell.AuthorReadWriter
	blogHandler BlogHandler
}

// NewAuthorHandler builds a new AuthorHandler.
func NewAuthorHandler(svc inkwell.AuthorReadWriter, blogHandler BlogHandler) AuthorHandler {
	return AuthorHandler{svc: svc, blogHandler: blogHandler}
}

func (h AuthorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var route string
	route, r.URL.Path = ShiftPath(r.URL.Path)

	switch route {
	case "":
		if r.Method != "POST" {
			badRequest(w, "Method not supported on this endpoint")
			return
		}

		// TODO: Create new author
		noContent(w)
	default:
		// route must be an AuthorID.
		h.handleAuthorID(route).ServeHTTP(w, r)
	}
}

func (h AuthorHandler) handleAuthorID(id string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var route string
		route, r.URL.Path = ShiftPath(r.URL.Path)

		switch route {
		case "":
			h.HandleGet(id).ServeHTTP(w, r)
			return
		case "blog":
			h.blogHandler.Handle(id).ServeHTTP(w, r)
			return
		default:
			notFound(w, "Endpoint does not exist")
			return
		}
	})
}

// HandleGet handles requests to get specific Authors.
func (h AuthorHandler) HandleGet(id string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blog, err := h.svc.Get(id)
		if err != nil {
			serverError(w, "Could not retrieve author")
			return
		}

		data, err := json.Marshal(&blog)
		if err != nil {
			serverError(w, "Something went wrong while marshaling author to JSON")
			return
		}

		ok(w, data)
	})
}

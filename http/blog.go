package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/eriktate/inkwell"
	"github.com/eriktate/inkwell/html"
	"github.com/pkg/errors"
	"github.com/pressly/chi"
	log "github.com/sirupsen/logrus"
)

// BlogHandler is an HTTP handler for dealing with Blog related requests.
type BlogHandler struct {
	svc inkwell.BlogService
}

// NewBlogHandler returns a new HTTP BlogHandler given an inkwell BlogService.
func NewBlogHandler(svc inkwell.BlogService) BlogHandler {
	return BlogHandler{
		svc: svc,
	}
}

// Get handles request for retrieving a specific blog if it exists.
func (h BlogHandler) Get(w http.ResponseWriter, r *http.Request) {
	log.Println("Calling GetBlog")
	authorID := chi.URLParam(r, "authorID")
	blogID := chi.URLParam(r, "blogID")

	// validate our inputs
	if authorID == "" || blogID == "" {
		log.WithError(errors.New("empty author or blog ID")).Println("Request to GetBlog failed")
		logFailure(badRequest(w, "You must provide a valid authorID and blogID."))
		return
	}

	// prep for error logging
	logFields := map[string]interface{}{
		"authorID": authorID,
		"blogID":   blogID,
	}

	// retrieve blog from service
	blog, err := h.svc.Get(authorID, blogID)
	if err != nil {
		log.WithError(err).WithFields(logFields).Println("Failed to GetBlog")
		logFailure(serverError(w, "Something went wrong while retrieving blog"))
		return
	}

	// check if we actually found a blog
	if blog.ID == "" {
		log.WithFields(logFields).Println("Blog does not exist")
		logFailure(notFound(w, "The blog requested could not be found"))
		return
	}

	data, err := json.Marshal(&blog)
	if err != nil {
		log.WithError(err).WithFields(logFields).Println("Failed to marshal Blog")
	}

	logFailure(ok(w, data))
}

// Write handles requests for writing a blog entry.
func (h BlogHandler) Write(w http.ResponseWriter, r *http.Request) {
	var blog inkwell.Blog

	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		log.WithError(err).Println("Failed to unmarshal Blog for writing")
		logFailure(badRequest(w, "Something went wrong while decoding the request. Make sure your payload is formatted correctly"))
		return
	}
	// nolint
	defer r.Body.Close()

	// attempt to write the blog
	if err := h.svc.Write(blog); err != nil {
		log.WithError(err).WithField("blog", blog).Println("Failed to write Blog")
		logFailure(serverError(w, "Something went wrong while writing Blog"))
		return
	}

	logFailure(noContent(w))
}

// Revise handles requests for revising a blog entry.
func (h BlogHandler) Revise(w http.ResponseWriter, r *http.Request) {
	authorID := chi.URLParam(r, "authorID")
	blogID := chi.URLParam(r, "blogID")

	// validate our inputs
	if authorID == "" || blogID == "" {
		log.WithError(errors.New("empty author or blog ID")).Println("Request to GetBlog failed")
		logFailure(badRequest(w, "You must provide a valid authorID and blogID."))
		return
	}

	// prep for error logging
	logFields := map[string]interface{}{
		"authorID": authorID,
		"blogID":   blogID,
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).WithFields(logFields).Println("Failed to read request body during revision")
		logFailure(badRequest(w, "Something went wrong while parsing revision"))
		return
	}
	// nolint
	defer r.Body.Close()

	if err := h.svc.Revise(authorID, blogID, string(data)); err != nil {
		log.WithError(err).WithFields(logFields).Println("Failed to publish blog")
		logFailure(serverError(w, "Something went wrong while publishing blog"))
		return
	}

	logFailure(noContent(w))
}

// Publish handles requests for publishing a blog entry.
func (h BlogHandler) Publish(w http.ResponseWriter, r *http.Request) {
	authorID := chi.URLParam(r, "authorID")
	blogID := chi.URLParam(r, "blogID")

	// validate our inputs
	if authorID == "" || blogID == "" {
		log.WithError(errors.New("empty author or blog ID")).Println("Request to GetBlog failed")
		logFailure(badRequest(w, "You must provide a valid authorID and blogID."))
		return
	}

	// prep for error logging
	logFields := map[string]interface{}{
		"authorID": authorID,
		"blogID":   blogID,
	}

	if err := h.svc.Publish(authorID, blogID); err != nil {
		log.WithError(err).WithFields(logFields).Println("Failed to publish blog")
		logFailure(serverError(w, "Something went wrong while publishing blog"))
		return
	}

	logFailure(noContent(w))
}

// Redact handles requests for redacting a blog entry.
func (h BlogHandler) Redact(w http.ResponseWriter, r *http.Request) {
	authorID := chi.URLParam(r, "authorID")
	blogID := chi.URLParam(r, "blogID")

	// validate our inputs
	if authorID == "" || blogID == "" {
		log.WithError(errors.New("empty author or blog ID")).Println("Request to GetBlog failed")
		logFailure(badRequest(w, "You must provide a valid authorID and blogID."))
		return
	}

	// prep for error logging
	logFields := map[string]interface{}{
		"authorID": authorID,
		"blogID":   blogID,
	}

	if err := h.svc.Redact(authorID, blogID); err != nil {
		log.WithError(err).WithFields(logFields).Println("Failed to redact blog")
		logFailure(serverError(w, "Something went wrong while redacting blog"))
		return
	}

	logFailure(noContent(w))
}

// GetPage retrieves the HTMl page for a requested blog.
func (h BlogHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	authorID := chi.URLParam(r, "authorID")
	blogID := chi.URLParam(r, "blogID")

	// validate our inputs
	if authorID == "" || blogID == "" {
		log.WithError(errors.New("empty author or blog ID")).Println("Request to GetBlog failed")
		logFailure(badRequest(w, "You must provide a valid authorID and blogID."))
		return
	}

	// prep for error logging
	logFields := map[string]interface{}{
		"authorID": authorID,
		"blogID":   blogID,
	}

	blog, err := h.svc.Get(authorID, blogID)
	if err != nil {
		log.WithError(err).WithFields(logFields).Println("Failed to GetPage")
		logFailure(serverError(w, "Something went wrong while retrieving page"))
		return
	}

	post, err := inkwell.BuildPost([]byte(blog.Content))
	if err != nil {
		log.WithError(err).WithFields(logFields).Println("Failed to build blog post")
		logFailure(serverError(w, "Something went wrong while retrieving page"))
		return
	}

	page := html.Page{
		Title: blog.Title,
		Post:  inkwell.RenderPost(post),
	}

	data, err := html.GeneratePage(html.BaseTemplate, page)
	if err != nil {
		log.WithError(err).WithFields(logFields).Println("Failed to generate page")
		logFailure(serverError(w, "Something weng wrong while retrieiving page"))
		return
	}

	logFailure(okHTML(w, data))
}

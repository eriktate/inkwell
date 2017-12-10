package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/eriktate/inkwell"
	"github.com/eriktate/inkwell/http"
	"github.com/eriktate/inkwell/mock"
	log "github.com/sirupsen/logrus"
)

func Test_Handle_Post(t *testing.T) {
	// SETUP
	blogData := []byte(`{"key": "some-blog", "title": "Some Blog"}`)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/", bytes.NewReader(blogData))
	mockSvc := &mock.MockBlogReadWriter{}
	authorID := "test_author"
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusNoContent)
}

func Test_Handle_Post_BadPayload(t *testing.T) {
	// SETUP
	// some-blog is missing surrounding quotes
	blogData := []byte(`{"key": some-blog, "title": "Some Blog"}`)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/", bytes.NewReader(blogData))
	mockSvc := &mock.MockBlogReadWriter{}
	authorID := "test_author"
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusBadRequest)
}

func Test_Handle_PostWithGet(t *testing.T) {
	// SETUP
	blogData := []byte(`{"key": "some-blog", "title": "Some Blog"}`)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodGet, "/", bytes.NewReader(blogData))
	mockSvc := &mock.MockBlogReadWriter{}
	authorID := "test_author"
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusBadRequest)
}

func Test_Handle_PostWithFailure(t *testing.T) {
	// SETUP
	blogData := []byte(`{"key": "some-blog", "title": "Some Blog"}`)
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/", bytes.NewReader(blogData))
	// preload mockSvc to fail
	mockSvc := &mock.MockBlogReadWriter{Fail: true}
	authorID := "test_author"
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusInternalServerError)
}

func Test_HandleKey_Get(t *testing.T) {
	// SETUP
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodGet, "/some-blog-key", nil)
	mockSvc := &mock.MockBlogReadWriter{}
	authorID := "test_author"
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	var blog inkwell.Blog
	res := recorder.Result()
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&blog); err != nil {
		t.Fatalf("Data returned is not valid JSON: %s", err)
	}

	if blog.AuthorID != authorID {
		log.Printf("Blog: %+v", blog)
		t.Fatal("Data returned does not match what was requested")
	}

	assertStatusCode(t, res.StatusCode, stdhttp.StatusOK)
}

func Test_HandleKey_GetFail(t *testing.T) {
	// SETUP
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodGet, "/some-blog-key", nil)
	mockSvc := &mock.MockBlogReadWriter{Fail: true}
	authorID := "test_author"
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusInternalServerError)
}

func Test_HandleKey_Publish(t *testing.T) {
	// SETUP
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key/publish", nil)
	mockSvc := &mock.MockBlogReadWriter{}
	authorID := "test_author"
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusNoContent)
}

func Test_HandleKey_PublishFail(t *testing.T) {
	// SETUP
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key/publish", nil)
	mockSvc := &mock.MockBlogReadWriter{Fail: true}
	authorID := "test_author"
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusInternalServerError)
}

func Test_HandleKey_lRedact(t *testing.T) {
	// SETUP
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key/redact", nil)
	mockSvc := &mock.MockBlogReadWriter{}
	authorID := "test_author"
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusNoContent)
}

func Test_HandleKey_RedactFail(t *testing.T) {
	// SETUP
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key/redact", nil)
	mockSvc := &mock.MockBlogReadWriter{Fail: true}
	authorID := "test_author"
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusInternalServerError)
}

func Test_HandleKey_SetContent(t *testing.T) {
	// SETUP
	content := "This is some blog content"
	authorID := "test_author"

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key/content", bytes.NewReader([]byte(content)))
	mockSvc := &mock.MockBlogReadWriter{}
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusNoContent)
}

func Test_HandleKey_SetContentFail(t *testing.T) {
	// SETUP
	content := "This is some blog content"
	authorID := "test_author"

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key/content", bytes.NewReader([]byte(content)))
	mockSvc := &mock.MockBlogReadWriter{Fail: true}
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusInternalServerError)
}

func Test_HandleKey_SetContent_BadBody(t *testing.T) {
	// SETUP
	authorID := "test_author"

	recorder := httptest.NewRecorder()
	body := &bogusBody{}
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key/content", body)
	mockSvc := &mock.MockBlogReadWriter{}
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusBadRequest)
}

func Test_HandleKey_SetTitle(t *testing.T) {
	// SETUP
	title := "Test Title"
	authorID := "test_author"

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key/title", bytes.NewReader([]byte(title)))
	mockSvc := &mock.MockBlogReadWriter{}
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusNoContent)
}

func Test_HandleKey_SetTitleFail(t *testing.T) {
	// SETUP
	title := "Test Title"
	authorID := "test_author"

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key/title", bytes.NewReader([]byte(title)))
	mockSvc := &mock.MockBlogReadWriter{Fail: true}
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusInternalServerError)
}

func Test_HandleKey_SetTitle_BadBody(t *testing.T) {
	// SETUP
	authorID := "test_author"

	recorder := httptest.NewRecorder()
	body := &bogusBody{}
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key/title", body)
	mockSvc := &mock.MockBlogReadWriter{}
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusBadRequest)
}

func Test_HandleKey_BadRoute(t *testing.T) {
	// SETUP
	authorID := "test_author"

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key/asdf", nil)
	mockSvc := &mock.MockBlogReadWriter{}
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusBadRequest)
}

func Test_HandleKey_BadPost(t *testing.T) {
	// SETUP
	authorID := "test_author"

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodPost, "/some-blog-key", nil)
	mockSvc := &mock.MockBlogReadWriter{}
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusBadRequest)
}

func Test_HandleKey_BadGet(t *testing.T) {
	// SETUP
	authorID := "test_author"

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(stdhttp.MethodGet, "/some-blog-key/asdf", nil)
	mockSvc := &mock.MockBlogReadWriter{}
	handler := http.NewBlogHandler(mockSvc, log.New()).Handle(authorID)

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	assertStatusCode(t, res.StatusCode, stdhttp.StatusBadRequest)
}

func assertStatusCode(t *testing.T, code, target int) {
	if code != target {
		t.Fatalf("Expecting status code %d, but received: %d", target, code)
	}
}

// implements the Reader interface for forcing failures.
type bogusBody struct{}

func (r *bogusBody) Read(p []byte) (int, error) {
	return 0, errors.New("Read failed")
}

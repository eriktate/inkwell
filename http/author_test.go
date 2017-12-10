package http_test

import (
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/eriktate/inkwell/http"
	"github.com/eriktate/inkwell/mock"
)

func Test_AuthorServe_BadMethod(t *testing.T) {
	// SETUP
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	handler := http.NewAuthorHandler(nil, http.BlogHandler{})

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	if res.StatusCode != stdhttp.StatusBadRequest {
		t.Fatal("Wrong status code returned")
	}
}

func Test_AuthorServe_Write(t *testing.T) {
	// SETUP
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", nil)
	handler := http.NewAuthorHandler(nil, http.BlogHandler{})

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	if res.StatusCode != stdhttp.StatusNoContent {
		t.Fatal("Wrong status code returned")
	}
}

func Test_AuthorServe_Get(t *testing.T) {
	// SETUP
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/some_author_id", nil)
	mockSvc := &mock.MockAuthorReadWriter{}
	handler := http.NewAuthorHandler(mockSvc, http.BlogHandler{})

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	if res.StatusCode != stdhttp.StatusOK {
		t.Fatal("Wrong status code returned")
	}

	if mockSvc.GetCalled != 1 {
		t.Fatal("Get was not called the correct number of times")
	}
}

func Test_AuthorServe_AuthorID_BadRoute(t *testing.T) {
	// SETUP
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/some_author_id/asdf", nil)
	mockSvc := &mock.MockAuthorReadWriter{}
	handler := http.NewAuthorHandler(mockSvc, http.BlogHandler{})

	// RUN
	handler.ServeHTTP(recorder, req)

	// ASSERT
	res := recorder.Result()

	if res.StatusCode != stdhttp.StatusNotFound {
		t.Fatal("Wrong status code returned")
	}
}

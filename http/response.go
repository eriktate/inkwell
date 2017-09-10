package http

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

//2xx range responses
func ok(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	return err
}

func okHTML(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	return err
}

func created(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write(data)
	return err
}

func noContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	_, err := w.Write(nil)
	return err
}

// 4xx range responses
func badRequest(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(message))
	return err
}

func notFound(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte(message))
	return err
}

func unprocessableEntity(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusUnprocessableEntity)
	_, err := w.Write([]byte(message))
	return err
}

// 5xx range responses
func serverError(w http.ResponseWriter, message string) error {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(message))
	return err
}

// util functions
func logFailure(err error) {
	if err != nil {
		log.WithError(err).Println("Failed to write response")
	}
}

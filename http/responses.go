package http

import "net/http"

func badRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
}

func notFound(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(msg))
}

func serverError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}

func ok(w http.ResponseWriter, data []byte) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func noContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}

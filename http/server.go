package http

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Server provides an HTTP interface to working with inkwell data.
type Server struct {
	host string
	port int

	handler http.Handler
}

// NewServer returns a new instance of a Server struct.
func NewServer(host string, port int) Server {
	return Server{
		host: host,
		port: port,
	}
}

// Start begins listening over HTTP on the host and port combination given
// to the Server struct. It's a blocking operation.
func (s Server) Start(handler http.Handler) error {
	log.WithField("host", s.host).WithField("port", s.port).Println("Starting server")
	s.handler = handler

	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), s.handler)
}

// Address returns the full address this server is currently listening on.
func (s Server) Address() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

// HealthCheck returns a 200 response to check if the web server is alive.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	logFailure(ok(w, []byte("{\"healthy\":true}")))
}

package http

import (
	"fmt"
	"net/http"

	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

// Server is just an HTTP server configuration.
type Server struct {
	host string
	port uint

	handler http.Handler
	log     *logrus.Logger
}

// NewServer returns a new Server struct given a hostname and a port.
func NewServer(host string, port uint, log *logrus.Logger) Server {
	return Server{host: host, port: port, log: log}
}

// Start attempts to bind to the address configured in the Server
// struct and start processing requests.
func (s Server) Start() {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.log.Printf("Starting web server on: %s", addr)

	s.log.Fatalf("Web server failed: %s", http.ListenAndServe(addr, s.handler))
}

// ShiftPath peels a single route paramter off of the path. It returns
// that paramter (head) and the rest of the path (tail).
func ShiftPath(p string) (string, string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}

	return p[1:i], p[i:]
}

type AppHandler struct {
	authorHandler AuthorHandler
}

func (h AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var route string
	route, r.URL.Path = ShiftPath(r.URL.Path)

	switch route {
	case "author":
		h.authorHandler.ServeHTTP(w, r)
	}
}

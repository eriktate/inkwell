package http_test

import (
	"fmt"
	"testing"

	inkhttp "github.com/eriktate/inkwell/http"
)

func Test_NewServer(t *testing.T) {
	// SETUP
	host := "localhost"
	port := 8080

	// RUN
	server := inkhttp.NewServer(host, port)

	// ASSERT
	if server.Address() != fmt.Sprintf("%s:%d", host, port) {
		t.Fatal("Server's address is wrong!")
	}
}

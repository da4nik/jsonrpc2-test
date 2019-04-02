package http

import (
	"fmt"
	"net/http"

	"github.com/da4nik/jrpc2_try/internal/log"
)

// Opts is options struct for http server
type Opts struct {
}

// Server is http server struct
type Server struct {
	port   int
	server *http.Server
}

// NewHTTPServer returns server intance
func NewHTTPServer(port int) (*Server, error) {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	server := Server{
		port:   port,
		server: srv,
	}

	return &server, nil
}

// Start - starts http server
func (s *Server) Start() error {
	log.Infof("Starting http server on port %d...", s.port)
	return fmt.Errorf("not implemented")
}

// Stop - stops http server
func (s *Server) Stop() error {
	log.Infof("Stopping http server...")
	return fmt.Errorf("not implemented")
}

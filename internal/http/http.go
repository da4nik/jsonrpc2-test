package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/da4nik/jrpc2_try/internal/log"
)

// Server is http server struct
type Server struct {
	port   int
	server *http.Server
	mux    *http.ServeMux
}

// NewHTTPServer returns server intance
func NewHTTPServer(port int, handler http.Handler) (*Server, error) {
	mux := http.NewServeMux()
	mux.Handle("/rpc", withMiddlewares(handler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	server := Server{
		port:   port,
		server: srv,
		mux:    mux,
	}

	return &server, nil
}

// Start - starts http server
func (s *Server) Start() {
	log.Infof("Starting http server on port %d...", s.port)

	go s.server.ListenAndServe()
}

// Stop - stops http server
func (s *Server) Stop() {
	log.Infof("Stopping http server...")
	s.server.Shutdown(context.Background())
	log.Infof("HTTP server stopped.")
}

func (s *Server) listen() {
	s.server.ListenAndServe()
}

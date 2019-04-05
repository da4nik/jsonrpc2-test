package http

import (
	"net/http"

	"github.com/da4nik/jrpc2_try/internal/log"
)

func withMiddlewares(handler http.Handler) http.Handler {
	return withLogger(handler)
}

func withLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("[%s] %s", r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

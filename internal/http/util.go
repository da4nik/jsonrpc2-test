package http

import (
	"net/http"

	"github.com/da4nik/jrpc2_try/internal/log"
)

func withLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("[%s] %s", r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

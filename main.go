package main

import (
	"net/http"
	"log"
	"log/slog"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request received", "path", r.URL.Path, "method", r.Method)
		next.ServeHTTP(w, r)
	})
}

func main() {
	server := NewServer(nil)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /get", server.HandleGet)
	mux.HandleFunc("PUT /set", server.HandleSet)

	err := http.ListenAndServe(":4000", logMiddleware(mux))
	log.Fatal(err)
}

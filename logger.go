package main

import (
	"log"
	"net/http"
	"time"
)

func logger(handle http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writer := StatusWriter{w, 0, 0}
		handle.ServeHTTP(&writer, r)
		end := time.Now()

		var query string

		if r.URL.RawQuery != "" {
			query = "?" + r.URL.RawQuery
		}

		log.Printf("%s %s \"%s %s %s\" %d %d \"%s\" %v",
			r.Host,
			r.RemoteAddr,
			r.Method,
			r.URL.Path+query,
			r.Proto,
			writer.Status,
			writer.Length,
			r.Header.Get("User-Agent"),
			end.Sub(start))
	})
}

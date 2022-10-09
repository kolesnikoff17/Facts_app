package mw

import (
	"log"
	"net/http"
	"time"
)

// Logging is a part of middleware. Logs to stdout Method, URI and execution time of each request
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

// PanicRecovery is a part of middleware. If any panic occurs, it will recover it
// and log stack trace, URL, Method and Body of request, caused panic
func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, http.StatusText(500), 500)
				log.Println(rec)
				log.Printf("Request, caused panic:")
				log.Println(r.URL)
				log.Println(r.Method)
				log.Println(r.Body)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

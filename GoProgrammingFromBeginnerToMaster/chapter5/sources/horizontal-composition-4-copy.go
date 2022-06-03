package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func greetings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("Welcome: %s\n", r.RemoteAddr))
}

func auth(s string) error {
	if s != "123" {
		return errors.New("auth failed")
	}
	return nil
}

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), now)
		h.ServeHTTP(w, r)
	})
}

func authHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := auth(r.Header.Get("auth")); err != nil {
			http.Error(w, fmt.Sprintf("%s\n", err), http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	http.ListenAndServe(":8080", logHandler(authHandler(http.HandlerFunc(greetings))))
	sync.Pool
}

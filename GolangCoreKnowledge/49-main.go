package main

import (
	"log"
	"net/http"
	"net/http/pprof"
	"strings"
)

func showHttpProfile() {
	mux := http.NewServeMux()
	pathPrefix := "/d/pprof/"
	mux.HandleFunc(pathPrefix,
		func(w http.ResponseWriter, r *http.Request) {
			name := strings.TrimPrefix(r.URL.Path, pathPrefix)
			if name != "" {
				pprof.Handler(name).ServeHTTP(w, r)
				return
			}
			pprof.Index(w, r)
		})
	mux.HandleFunc(pathPrefix+"cmdline", pprof.Cmdline)
	mux.HandleFunc(pathPrefix+"profile", pprof.Profile)
	mux.HandleFunc(pathPrefix+"symbol", pprof.Symbol)
	mux.HandleFunc(pathPrefix+"trace", pprof.Trace)

	server := http.Server{
		Addr:    "localhost:8083",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func main() {
	showHttpProfile()
}

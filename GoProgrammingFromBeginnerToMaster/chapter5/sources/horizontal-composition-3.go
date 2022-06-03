package main

import (
	"fmt"
	"net/http"
)

func greetings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!\n")
}

func bar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprintf("/access bar from %s\n", r.RemoteAddr))
}

func main() {
	// http.HandleFunc("/", greetings)
	// http.HandleFunc("/bar", bar)
	// http.ListenAndServe(":8080", nil)
	http.ListenAndServe(":8080", http.HandlerFunc(greetings))
}

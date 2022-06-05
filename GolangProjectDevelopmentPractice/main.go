package main

import (
	"log"
	"net/http"
)

type S struct {
	a int
	b string
	c bool
}

func chapter12() {
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world\n"))
	})
	log.Println("Starting http server ...")
	log.Fatal(http.ListenAndServe(":50052", nil))
}

func pong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong\n"))
}

func main() {
	chapter12()
}

package main

import (
	"log"
	"net/http"
	// _ "net/http/pprof"
)

func main() {

	http.ListenAndServe(":8000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("!!! %+v", r.Header)
		w.Write([]byte("Hello, World!"))
	}))
}

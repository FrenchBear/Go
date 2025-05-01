// server1.go
// Learning Go, Th Gro Programming Language, chap 1.7
// Minimal web server
//
// When running, in a browser, open http://localhost:8000/once/upon/a/time
// response: URL.Path = "/once/upon/a/time"
//
// 2019-09-14	PV

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler) // Each request call handler
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the path component of requested URL
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

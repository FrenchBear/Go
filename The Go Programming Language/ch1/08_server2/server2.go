// server2.go
// Learning Go, Th Gro Programming Language, chap 1.7
// Minimal web echo and counter server
//
// When running, in a browser, open http://localhost:8000/once/upon/a/time
// response: URL.Path = "/once/upon/a/time"
// open http://localhost:8000/count
// response: Count 2
//
// 2019-09-26	PV

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count int

func main() {
	http.HandleFunc("/", handler) // Each request call handler
	http.HandleFunc("/count", counter)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the Path component of the requested URL.
func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// counter echoes the number of calls so far.
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	mu.Unlock()
}

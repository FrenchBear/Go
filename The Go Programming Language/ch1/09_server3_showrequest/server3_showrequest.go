// server3_showrequest.go
// Learning Go, Th Gro Programming Language, chap 1.7
// Minimal web server, echoes the HTTP request
//
// When running, in a browser, open http://localhost:8000/once/upon/a/time
// response:
// GET /once/upon/a/time HTTP/1.1
// Header["Dnt"] = ["1"]
// Header["User-Agent"] = ["Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36"]
// Header["Accept-Encoding"] = ["gzip, deflate, br"]
// Header["Accept-Language"] = ["en-US,en;q=0.9,fr;q=0.8,fr-FR;q=0.7"]
// Header["Connection"] = ["keep-alive"]
// Header["Upgrade-Insecure-Requests"] = ["1"]
// Header["Sec-Fetch-Mode"] = ["navigate"]
// Header["Accept"] = ["text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3"]
// Header["Sec-Fetch-Site"] = ["none"]
// Host = "localhost:8000"
// RemoteAddr = "127.0.0.1:55710"
//
// 2019-09-26	PV

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

// handler echoes the HTTP request.
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

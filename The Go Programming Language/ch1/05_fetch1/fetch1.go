// fetch1.go
// The Go Programming Language, chapter 1.5
// Simple get text from url(s) passed as parameter
//
// 2019-08-29	PV
// 2019-10-26	PV		Added test resp.StatusCode != http.StatusOK

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, urlarg := range os.Args[1:] {
		url := urlarg
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + urlarg
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch1: error %v\n", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fmt.Fprintf(os.Stderr, "fetch1: error accessing %s: %s", url, resp.Status)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch1: error reading %s: %v\n", url, err)
			continue
		}
		fmt.Printf("fetch %s -> HTTP %v:\n%s\n", url, resp.Status, body)
	}
}

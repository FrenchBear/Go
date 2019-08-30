// fetchall.go
// Learning Go, Th Gro Programming Language, chap 1.6
// Fetches multiple URLs in parallel
//
// Output example:
// C:\Development\GitHub\Go\The Go Programming Language\ch1\06_fetchall>06_fetchall.exe https://www.google.com https://golang.org http://gopl.io
// 0.77s    11571  https://www.google.com
// 0.84s    11062  https://golang.org
// 1.28s     4154  http://gopl.io
// 1.28s elapsed
//
// 2019-08-29	PV

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // Start a subroutine
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

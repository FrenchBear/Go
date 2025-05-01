// counttags.go
// The Go Programming Language, chapter 5.2
// Count HTML tags of a web page
//
// 2019-10-26	PV

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html" // Need to run 'go get golang.org/x/net/html'
)

const urlDefault = `https://xkcd.com/`

func main() {
	var url string
	if len(os.Args) == 2 {
		url = os.Args[1]
	} else if len(os.Args) == 1 {
		url = urlDefault
	} else {
		fmt.Println("Usage: counttags [url]")
		return
	}
	page, err := getPage(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "counttags: %v\n", err)
		os.Exit(1)
	}
	pageReader := bytes.NewReader(page)
	doc, err := html.Parse(pageReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "counttags: %v\n", err)
		os.Exit(1)
	}
	m := make(map[string]int)
	countTags(m, doc)

	fmt.Println("Tags count on page", url)
	for k, v := range m {
		fmt.Printf("%-10s %d\n", k, v)
	}
}

func getPage(url string) ([]byte, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "counttags: error %v\n", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("counttags: error accessing %s: %s", url, resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "counttags: error reading %s: %v\n", url, err)
		return nil, err
	}
	return body, nil
}

func countTags(m map[string]int, n *html.Node) {
	if n.Type == html.ElementNode {
		m[n.Data]++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		countTags(m, c)
	}
}

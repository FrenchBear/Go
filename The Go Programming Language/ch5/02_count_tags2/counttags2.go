// counttags2.go
// The Go Programming Language, chapter 5.2
// Count HTML tags of a web page, variant: uses a pageReader, and sorts output by tag and count
//
// 2019-10-26	PV

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"golang.org/x/net/html" // Need to run 'go get golang.org/x/net/html'
)

const urlDefault = `https://xkcd.com`

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
	pageReader, err := getPageReader(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "counttags: %v\n", err)
		os.Exit(1)
	}
	//defer func() { pageReader.Close() }()
	doc, err := html.Parse(pageReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "counttags: %v\n", err)
		os.Exit(1)
	}
	pageReader.Close()
	tagCount := make(map[string]int)
	countTags(tagCount, doc)

	fmt.Println("Tags count on page", url, " tag order")
	var tags []string
	for tag := range tagCount {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	for _, t := range tags {
		fmt.Printf("%-10s %d\n", t, tagCount[t])
	}

	fmt.Println("\nTags count on page", url, " count order")
	countTags := make(map[int][]string)
	counts := []int{}
	for tag, count := range tagCount {
		slice, ok := countTags[count]
		if !ok {
			slice = []string{}
			counts = append(counts, count)
		}
		countTags[count] = append(slice, tag)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	for _, count := range counts {
		for _, tag := range countTags[count] {
			fmt.Printf("%-10s %d\n", tag, count)
		}
	}
}

func getPageReader(url string) (io.ReadCloser, error) {
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
	return resp.Body, nil
}

func countTags(m map[string]int, n *html.Node) {
	if n.Type == html.ElementNode {
		m[n.Data]++
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		countTags(m, c)
	}
}

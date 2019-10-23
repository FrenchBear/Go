// wordfreq.go
// Exercise 4.9: Write a program wordfreq to report the frequency of each word in an input text file.
// Call input.Split(bufio.ScanWords) before the first call to Scan to break the input into words instead of lines.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Analyzing stdin")
		countWords(os.Stdin)
	} else if len(os.Args) == 2 {
		f, err := os.Open(os.Args[1])
		if err == nil {
			fmt.Printf("Analyzing file '%s'\n", os.Args[1])
			countWords(f)
		} else {
			fmt.Printf("Error opening file '%s': %v\n", os.Args[1], err)
		}
	} else {
		fmt.Println("Usage: wordfreq [file]")
	}
}

func countWords(f io.Reader) {
	input := bufio.NewScanner(f)
	input.Split(bufio.ScanWords)

	// First count words
	fmt.Println("Counting")
	wordsCount := map[string]int{}
	for input.Scan() {
		wordsCount[input.Text()]++
	}

	// Then group by freq=count
	fmt.Println("Goup words by freq")
	freqWords := map[int][]string{}
	for w, l := range wordsCount {
		list, ok := freqWords[l]
		if !ok {
			list = []string{}
		}
		list = append(list, w)
		freqWords[l] = list
	}

	// Sort frequencies by decreasing order
	fmt.Println("Sorting frequencies")
	sortFreq := []int{}
	for l := range freqWords {
		sortFreq = append(sortFreq, l)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortFreq)))

	// Print top 10
	fmt.Println("Top 10 most frequent words")
	n := 0
	for _, l := range sortFreq {
		for _, w := range freqWords[l] {
			fmt.Printf("%6d %s\n", l, w)
			n++
			if n > 10 {
				goto exit
			}
		}
	}
exit:
}

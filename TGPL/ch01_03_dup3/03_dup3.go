package main

// Learning go
// dup3: print lines appearing more than once from a list of provided files with filename and count per file
// exercise 1.4 in "The Go Programming Language"
//
// 2019-08-24	PV

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	counts2 := make(map[string][]string)
	files := os.Args[1:]
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "dup3: need file[s] as arguments")
		return
	}
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
		} else {
			for _, line := range strings.Split(string(data), "\n") {
				counts2[line] = append(counts2[line], file)
			}
		}
	}
	for line, arr := range counts2 {
		if len(arr) > 1 {
			m2 := make(map[string]int)
			for _, file := range arr {
				m2[file]++
			}
			fmt.Printf("%d\t", len(arr))
			for file, count := range m2 {
				fmt.Printf("%s:%d ", file, count)
			}
			fmt.Printf("\t%s\n", line)
		}
	}
}

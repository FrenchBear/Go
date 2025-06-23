// g25_files_2.go
// Learning go, System programming, files, reading a text file line by line
//
// 2025-06-23	PV		First version

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println("Go files #2")

	args := os.Args
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: g25_files_2 <textfile>")
		os.Exit(1)
	}

	err := lineByLine(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", args[1], err)
		os.Exit(2)
	}
}

func lineByLine(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// After making sure that you can open the given file for reading (os.Open()), you create a new reader using
	// bufio.NewReader().
	r := bufio.NewReader(f)

	for {
		// bufio.ReadString() returns two values: the string that was read and an error variable.
		line, err := r.ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			break
		}

		// The use of fmt.Print() instead of fmt.Println() for printing the input line shows that the newline character
		// is included in each input line.
		fmt.Print(line)
	}
	return nil
}

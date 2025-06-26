// g27_write_text.go
// Learning go, System programming, files, writing a text file
//
// 2025-06-26	PV		First version

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println("Go Writing a text file")
	fmt.Println()

	buffer := []byte("Data to write\n")
	f1, err := os.Create("C:/tmp/f1.txt")
	// os.Create() returns an *os.File value associated with the file path that is passed as
	// a parameter. Note that if the file already exists, os.Create() truncates it.
	if err != nil {
		fmt.Println("Cannot create file", err)
		return
	}
	defer f1.Close()
	fmt.Fprintf(f1, "%s", string(buffer))

	// The os.WriteString() method writes the contents of a string to a valid *os.File variable.
	f2, err := os.Create("C:/tmp/f2.txt")
	if err != nil {
		fmt.Println("Cannot create file", err)
		return
	}
	defer f2.Close()
	n, err := f2.WriteString(string(buffer))
	fmt.Printf("wrote %d bytes\n", n)

	// Here we create a temporary file on our own. Later on in this chapter you are going to
	// learn about using os.CreateTemp() for creating temporary files.
	f3, err := os.Create("C:/tmp/f3.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	w := bufio.NewWriter(f3)
	// This function returns a bufio.Writer, which satisfies the io.Writer interface.
	n, err = w.WriteString(string(buffer))
	fmt.Printf("wrote %d bytes\n", n)
	w.Flush()
	// Don't close it?

	f4_name := "C:/tmp/f4.txt"
	f4, err := os.Create(f4_name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		fmt.Println("f4 1st version is closed")
		f4.Close()
	}()

	for i := 0; i < 5; i++ {
		//n, err = io.WriteString(f4, string(buffer))		// Better to call directly io.Writer.Write that convert a byte[] to string and call io.Write
		n, err = io.Writer.Write(f4, buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("wrote %d bytes\n", n)
	}
	// Append to a file
	fmt.Println("About to redefine f4")
	f4, err = os.OpenFile(f4_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fmt.Println("f4 reopened")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		fmt.Println("f4 2nd version is closed")
		f4.Close()
	}()
	
	// Write() needs a byte slice
	n, err = f4.Write([]byte("Put some more data at the end.\n"))

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("wrote %d bytes\n", n)
}

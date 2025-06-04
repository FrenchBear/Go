// g04_error
// Learning go, Update of g02_module, but Hell function can now return an error
// need to run "go mod edit -replace example.com/greetings=./greetings"
// https://go.dev/doc/tutorial/handle-errors
//
// 2025-06-04	PV		First version

package main

import (
	"fmt"
	"log"

	"example.com/greetings"
)

func main() {
    // Set properties of the predefined Logger, including
    // the log entry prefix and a flag to disable printing
    // the time, source file, and line number.
    log.SetPrefix("greetings: ")
    log.SetFlags(0)
	
    // Request a greeting message.
    message, err := greetings.Hello("")
    // If an error was returned, print it to the console and exit the program.
    if err != nil {
        log.Fatal(err)
    }

    // If no error was returned, print the returned message to the console.
    fmt.Println(message)
}

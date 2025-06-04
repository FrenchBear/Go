// g02_module
// Learning go, Example of code calling a function in a separate module
// This is test module calling module greetings in a subfolder
// need to run "go mod edit -replace example.com/greetings=./greetings"
// https://go.dev/doc/tutorial/create-module
//
// 2025-06-04	PV		First version

package main

import (
	"fmt"

	"example.com/greetings"
)

func main() {
	// Get a greeting message and print it.
	message := greetings.Hello("Gladys")
	fmt.Println(message)
}

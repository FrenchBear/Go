// g01_package
// Learning go, basic example of package referencing an external library
// https://go.dev/doc/tutorial/getting-started
// Need to run "go mod tidy" to update go.mod
//
// 2025-06-04	PV		First version

package main

import (
	"fmt"

	"rsc.io/quote"
)

func main() {
	fmt.Println("Hello, world")
	fmt.Println(quote.Go())
}

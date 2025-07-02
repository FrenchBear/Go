// mymarkup_main.go
// Main package (simple testing) for package MyMarkup
//
// go mod edit -replace github.com/PieVio/MyMarkup=C:\Development\GitHub\Go\Packages\MyMarkup
// go mod tidy
//
// 2025-06-23	PV		First version

package main

import (
	"fmt"

	MyMarkup "github.com/PieVio/MyMarkup"
)

func main() {
	fmt.Printf("MyMarkup lib version: %s\n\n", MyMarkup.Version())

	MyMarkup.RenderMarkup("⌊Hello⌋, ⟪world⟫⦃!⦄")
}


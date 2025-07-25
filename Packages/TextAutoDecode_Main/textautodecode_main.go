// textautodecode_main.go
// Main package (simple testing) for package TextAutoDecode
//
// go mod edit -replace github.com/PieVio/TextAutoDecode=C:\Development\GitHub\Go\Packages\TextAutoDecode
// go mod tidy
//
// 2025-07-05	PV		First version

package main

import (
	"fmt"

	TextAutoDecode "github.com/PieVio/TextAutoDecode"
)

func main() {
	fmt.Printf("TextAutoDecode lib version: %s\n\n", TextAutoDecode.Version())

	tadRes, err := TextAutoDecode.ReadTextFile(`C:\DocumentsOD\Doc tech\Unicode\Marque d'ordre des octets - Wikipédia.website`)
	fmt.Println("tadRes:", tadRes)
	fmt.Println("err: ", err)
}


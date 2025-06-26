// textautodecode_tests.go
// Tests for package TestAutoDecode
//
// go mod tidy
// go mod edit -replace github.com/PieVio/TextAutoDecode=C:\Development\GitHub\Go\Packages\TextAutoDecode
//
// 2025-06-23	PV		First version

package main

import (
	"fmt"
	"strings"

	textautodecode "github.com/PieVio/TextAutoDecode"
)

func main() {
	fmt.Println("Tests of package TextAutoDecode ", textautodecode.Version())
	fmt.Println()

	test(`C:\DocumentsOD\Doc tech\Encodings\inexistent`)
	test(`.\textautodecode_tests.exe`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-empty.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-ascii.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf8bom.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16lebom.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16bebom.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf8.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16le.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16be.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-1252.txt`)
}

func test(filename string) {
	tad, err := textautodecode.ReadTextFile(filename)
	if err != nil {
		fmt.Printf("%-65.65s Err: %v\n", filename, err)
		return
	}

	fmt.Printf("%-65.65s %s\n", filename, tad.Encoding.ToString())
	if strings.Contains(filename, "empty") || strings.Contains(filename, ".exe") {
		return
	}

	var beginning string
	if strings.Contains(filename, "ascii") {
		beginning = "juliette sophie brigitte geraldine"
	} else {
		beginning = "juliette sophie brigitte géraldine"
	}

	if !strings.HasPrefix(tad.Text, beginning) {
		l := min(len(tad.Text), 80)
		fmt.Println("Wrong prefix:", "«"+tad.Text[:l]+"»")
	}
}

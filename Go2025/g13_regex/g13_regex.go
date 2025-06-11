// g13_regex.go
// Learning go, Regular expressions
//
// 2025-06-11	PV		First version

package main

import (
	"fmt"
	"regexp"
)

func matchNameSur(s string) bool {
	t := []byte(s)
	re := regexp.MustCompile(`^[A-Z][a-z]*$`)
	return re.Match(t)
}

func matchInt(s string) bool {
t := []byte(s)
re := regexp.MustCompile(`^[-+]?\d+$`)
return re.Match(t)
}

func main() {
	fmt.Println("Regex in Go")
	fmt.Println()

	fmt.Println(matchNameSur("Pierre"))		// true
	fmt.Println(matchNameSur("Jean-Paul"))	// false

	fmt.Println(matchInt("-355"))		// true
	fmt.Println(matchInt("12 345"))		// false

}

// loose_ends.go
// Learning Go, Th Gro Programming Language, chap 1.8
//
// 2019-09-26	PV

package main

import "fmt"

// example of structure
type point struct {
	X, Y int
}

var p point


func main() {
	heads := 0
	tails := 0

	// usual switch example
	switch coinflip() {
	case "heads":
		heads++
	case "tails":
		tails++
	default:
		fmt.Println("landed on edge!")
	}
}

func coinflip() string {
	return "heads"
}

func signum(x int) int {
	// tagless switch example
	switch {
	case x > 0:
		return +1
	default:
		return 0
	case x < 0:
		return -1
	}
}

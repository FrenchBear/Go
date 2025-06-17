// mathelem.go
// Learning go, my first package
//
// 2025-06-17	PV		First version

package MathElem

import (
	"fmt"
)

func init() {
	fmt.Println("Running first init() of MathElem package")
}

// A packahe can contain multiple init() functions, they are executed in order of declaration (before main() for package main)
func init() {
	fmt.Println("Running second init() of MathElem package")
}

func Square(x float64) float64 {
	return x*x
}

func Cube(x float64) float64 {
	return x*x*x
}

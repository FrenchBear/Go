// g19_functions.go
// Learning go, simple code related to functions
//
// 2025-06-16	PV		First version

package main

import (
	"fmt"
)

func main() {
	// An anonymous function
	square := func(x float64) float64 {
		return x * x
	}
	fmt.Println("25²:", square(25))

	res := doublecall(func(x int) int { return +1 }, 5)
	fmt.Println("res:", res)

	// Parentheses are required for return type of function definition, but they're forbidden here...
	mi, ma := sortTwo(5, 3)
	fmt.Println("mi:", mi, "  ma:", ma)

	mi, ma = minMax(5, 3)
	fmt.Println("mi:", mi, "  ma:", ma)

	f1 := funRet(1)
	fmt.Println("f1(3):", f1(3))
	f2 := funRet(-1)
	fmt.Println("f2(3):", f2(3))

	s:=addFloats("Summing floats", 1.8, -2.5, 3.1416, 1.1414)
	fmt.Println("s:", s)

	everything("pomme", 42, false, 3.1416)
}

// Parameter of type function
func doublecall(f func(x int) int, i int) int {
	// defer executes stuff when function terminates
	defer func() {
		fmt.Println("defer called lambda")
	}() 		// () since defer needs a function call

	// Last defer statement is executed first
	for i:=3 ; i>0 ; i-- {
		defer fmt.Println("i:", i)
	}

	return f(f(i))
}

// Function returning two values, () are mandatory in the return type
// Sorting from smaller to bigger value
func sortTwo(x, y int) (int, int) {
	if x > y {
		return y, x
	}
	return x, y
}

// Return values can be named
func minMax(x, y int) (min, max int) {
	if x > y {
		min = y
		max = x
		return min, max
	}
	min = x
	max = y
	return // Return the values stored in min and max
}

// Functions can return functions
func funRet(i int) func(int) int {
	if i < 0 {
		return func(k int) int {
			k = -k
			return k + k
		}
	}
	return func(k int) int {
		return k * k
	}
}

// Variadoc function
func addFloats(message string, s ...float64) float64 {
	fmt.Println(message)
	sum := float64(0)
	for _, a := range s {
		sum = sum + a
	}
	return sum
}

// Another variadic function
func everything(input ...interface{}) {
	fmt.Println(input...)
}

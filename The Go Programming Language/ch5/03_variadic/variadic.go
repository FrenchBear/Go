// variadic.go
// The Go Programming Language, chapter 5.7
// Variadic functions examples
//
// 2019-12-06	PV

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(sum())           //  "0"
	fmt.Println(sum(3))          //  "3"
	fmt.Println(sum(1, 2, 3, 4)) //  "10"

	values := []int{1, 2, 3, 4}
	fmt.Println(sum(values...)) // "10"

	// Signature is different
	fmt.Printf("%T\n", f) // "func(...int)"
	fmt.Printf("%T\n", g) // "func([]int)"

	errorf(12, "Test of error %d: %s", 5, "Nil pointer dereferenced")

	fmt.Println(joinString("[", ",", "]", "Once", "Upon", "a", "time"))
}

func f(...int) {}
func g([]int)  {}

func sum(vals ...int) int {
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}

func errorf(linenum int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Line %d: ", linenum)
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
}

func joinString(begin, mid, end string, str ...string) string {
	var res string
	for _, s := range str {
		if len(res) == 0 {
			res = begin
		} else {
			res += mid
		}
		res += s
	}
	return res + end
}

// misc.go
// Just a simple example while learning go
// 2019-09-29	PV

package main

import "fmt"

func main() {
	r := getR()
	fmt.Println("r =", r)

	// Print a pointer
	p := &r
	fmt.Println(p) // 0xc000012098

	// Difference of address when it leaks or not
	pv := loc1()
	loc2()
	*pv = 0
	v3 := new(int)
	fmt.Println("v3:", &v3)

	// Function type
	fmt.Printf("%T\n", sum)

	// Test half
	h1a, h1b := half(1)
	h2a, h2b := half(2)
	fmt.Println("half(1):", h1a, h1b)
	fmt.Println("half(2):", h2a, h2b)

	// Test max
	fmt.Println("max:", max(3.14, -2.78, 1.73, 1.41, 6.28, 0, 3.33))

	// Test Fibonacci generator
	fiboGen := makeFiboGen()
	for i := 1; i < 10; i++ {
		fmt.Printf("%d ", fiboGen())
	}
	fmt.Println()

	d := dog{"Cubitus"}
	c := cat{"Gros minet"}
	b := bear{"Baloo"}
	crieMeute(&d, &c, b) // Can't use a pointer receiver with an object
	crieMeute(&b)        // But can use an object receiver with a pointer
}

func crieMeute(animals ...animal) {
	for _, a := range animals {
		a.cri()
	}
}

// sum is a function which takes a slice of numbers
// and adds them together. What would its function
// signature look like in Go?
func sum(numbers ...int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

// Write a function which takes an integer and
// halves it and returns true if it was even or false
// if it was odd. For example half(1) should return
// (0, false) and half(2) should return (1, true).
func half(n int) (int, bool) {
	n2 := n >> 1
	return n2, n%2 == 0
}

// Write a function with one variadic parameter
// that finds the greatest number in a list of numbers.
func max(numbers ...float64) float64 {
	m := -1e100
	for _, f := range numbers {
		if f > m {
			m = f
		}
	}
	return m
}

func getR() (r int) {
	// Just to make sure that defer is executed last, and overrides the 1 in "return 1"
	defer func() {
		r = 3
	}()
	r = 0
	return 1
}

func loc1() *int {
	v1 := 0
	fmt.Println("&v1:", &v1)
	return &v1
}

func loc2() {
	v2 := 0
	fmt.Println("&v2:", &v2)
}

// Fibonacci generator
func makeFiboGen() func() uint {
	last := 1
	x1 := uint(1)
	x2 := uint(1)
	return func() uint {
		if last <= 2 {
			last++
			return 1
		}
		x1, x2 = uint(x1+x2), x1
		return x1
	}
}

type cat struct {
	name string
}

type dog struct {
	name string
}

type bear struct {
	name string
}

// Method with a pointer reveiver
func (d *dog) cri() {
	fmt.Println(d.name, "Ouah!")
}

// Simple function
func cri() {
	fmt.Println(c.name, "Miaou!")
}

// Method with an object receiver
func (b bear) cri() {
	fmt.Println(b.name, "Grrr!")
}

// An interface is just a list of functions
// How it's implement does not matter as long as it exists
// Both dog, bear and cat implement this interface
type animal interface {
	cri()
}

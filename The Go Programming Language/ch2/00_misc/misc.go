// misc.go
// Just a simple example while learning go
// 2019-09-29	PV

package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	r := getR()
	fmt.Println("r =", r)

	// Print a pointer
	pr := &r
	fmt.Println(pr) // 0xc000012098

	// Difference of address when it leaks or not
	pv := loc1()
	loc2()
	*pv = 0
	v3 := new(int)
	fmt.Println("v3:", &v3)

	// Function type
	fmt.Printf("%T\n", sum) // func(...int) int

	// Test half
	h1a, h1b := half(1)
	h2a, h2b := half(2)
	fmt.Println("half(1):", h1a, h1b) // 0 false
	fmt.Println("half(2):", h2a, h2b) // 1 true

	// Test max
	fmt.Println("max:", maxFloat64(3.14, -2.78, 1.73, 1.41, 6.28, 0, 3.33)) // 6.28
	fmt.Println("max (no arg):", maxFloat64())                              // NaN

	// Test Fibonacci generator
	fiboGen := makeFiboGen()
	var f20 uint64
	for i := 1; i <= 20; i++ {
		f20 = fiboGen()
		fmt.Printf("%d ", f20)
	}
	fmt.Println()
	fmt.Println("fiboGen(20):", f20)

	// Test Fibonacci memoizer
	f20 = fiboMem(20)
	fmt.Println("fiboMem(20):", f20)

	// Test Fibonacci direct
	f20 = fiboDir(20)
	fmt.Println("fiboDir(20):", f20)

	// Tests of receivers and interfaces
	d := dog{"Cubitus"}
	c := cat{"Gros minet"}
	b := bear{"Baloo"}
	k := duck{"Donald"}
	//p := puppy{dog{"Dago"}}
	var p puppy
	p.name = "Dago"
	crieMeute(&d, &c, b, &p) // Can't use a pointer receiver with an object
	crieMeute(&b)            // But can use an object receiver with a pointer
	cri(&k)
	// Test meute which is an array of animal interface that itself implements animal interface
	var m meute
	m.animal = []animal{&d, &c, &b}
	crieMeute(&p, &m)

	dm := dog{"Medor"}
	ta := []animal{&d, &dm}
	crieMeute(ta...) // Expansion of array during call

	// Strings
	fmt.Println(
		strings.Contains("Aé♫山𝄞🐗", "♫山𝄞"),               // true
		strings.Count("Aé♫山𝄞🐗Aé♫山𝄞🐗", "𝄞"),              // 2
		strings.HasPrefix("♫山𝄞🐗", "♫山"),                 // true
		strings.HasSuffix("♫山𝄞🐗", "𝄞🐗"),                 // true
		strings.Index("♫山𝄞🐗", "𝄞"),                      // 6
		strings.Join([]string{"♫", "山", "𝄞", "🐗"}, "●"), // ♫●山●𝄞●🐗
		strings.Repeat("♫", 5),                          // ♫♫♫♫♫
		strings.Replace("♫山♫🐗♫", "♫", "♪", 2),           // ♪山♪🐗♫
		strings.Split("A●é●♫●山●𝄞●🐗", "●"),               // [A é ♫ 山 𝄞 🐗]
	)

	// new, make and delete
	pi := new(int) // Returns int *
	*pi = 2

	pj := make([]int, 5, 10) // Create array of 10 ints, sclice it to first 5 ints
	pj[0] = 3

	col := make(map[string]string)
	col["red"] = "rouge"
	delete(col, "red")
}

// Variadic function
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
func maxFloat64(numbers ...float64) float64 {
	if len(numbers) == 0 {
		return math.NaN() // It's a function, not a constant
	}
	m := -math.MaxFloat64
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
// Returns a function that captures its state in local variables
func makeFiboGen() func() uint64 {
	last := 1
	x1, x2 := uint64(1), uint64(1)
	return func() uint64 {
		// The first two values are not computed but returned directly
		if last <= 2 {
			last++
			return 1
		}
		x1, x2 = x1+x2, x1
		return x1
	}
}

// Recursive computation of Fibonacci sequence using a memoizer
var fiboMemoizer = map[uint64]uint64{1: 1, 2: 1}

func fiboMem(n uint64) uint64 {
	if s, ok := fiboMemoizer[n]; ok {
		return s
	}
	s := fiboMem(n-1) + fiboMem(n-2)
	fiboMemoizer[n] = s
	return s
}

func fiboDir(n uint) uint64 {
	x, y := uint64(1), uint64(1)
	for n > 2 {
		x, y = x+y, x
		n--
	}
	return x
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

type duck struct {
	name string
}

// Method with a pointer receiver
func (d *dog) cri() {
	fmt.Println(d.name, "Ouah!")
}

// Simple function
func (c *cat) cri() {
	fmt.Println(c.name, "Miaou!")
}

// Method with an object receiver
// Beware, method gets a copy of a the receiver and can't update it
func (b bear) cri() {
	fmt.Println(b.name, "Grrr!")
}

func cri(d *duck) {
	fmt.Println(d.name, "Quack!")
}

// An interface is just a list of functions
// How it's implement does not matter as long as it exists
// Both dog, bear and cat implement this interface
// But it must be implemented through a method, duck does not implement this interface
type animal interface {
	cri()
}

// Can define an array of interfaces
type meute struct {
	animal []animal
}

// That itself implement the interface
func (m *meute) cri() {
	for _, a := range m.animal {
		a.cri()
	}
}

// Inheritance
type puppy struct {
	dog
}

func (p *puppy) cri() {
	fmt.Println(p.name, "Wif!")
}

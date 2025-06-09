// g09_data_types.go
// Learning go, data types
//
// 2025-06-06	PV		First version

package main

import (
	"fmt"
	"sort"
	"strconv"
	s "strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type Digit int
type Power2 int

const PI = 3.1415926
const BIG = 123_456_789
const (
	C1 = "C1C1C1"
	C2 = "C2C2C2"
	C3 = "C3C3C3"
)

func main() {
	numeric()
	conversion_int_to_string()
	string_character_rune()
	test_unicode()
	test_strings()
	test_time_date()
	test_iota()
	test_arrays_slices()
	test_pointers()
}

func numeric() {
	fmt.Println("--- numeric")

	c1 := 12 + 1i
	c2 := complex(5, 7)
	fmt.Printf("Type of c1: %T\n", c1)
	fmt.Printf("Type of c2: %T\n", c2)
	var c3 complex64 = complex64(c1 + c2)
	fmt.Println("c3:", c3)
	fmt.Printf("Type of c3: %T\n", c3)
	cZero := c3 - c3
	fmt.Println("cZero:", cZero)

	x := 12
	k := 5
	fmt.Println(x)
	fmt.Printf("Type of x: %T\n", x)

	div := x / k //	int/int -> int
	fmt.Println("div", div)

	var m, n float64
	m = 1.223
	fmt.Println("m, n:", m, n)

	y := 4 / 2.3
	fmt.Println("y:", y)

	divFloat := float64(x) / float64(k)
	fmt.Println("divFloat", divFloat)
	fmt.Printf("Type of divFloat: %T\n\n", divFloat)
}

func conversion_int_to_string() {
	fmt.Println("--- conversions_int_to_string")

	n := 100
	_ = strconv.Itoa(n)
	_ = strconv.FormatInt(int64(n), 10)
	r := string(n) // Generates 'd' rune
	fmt.Println("r=", r)
	fmt.Println()
}

func string_character_rune() {
	fmt.Println("--- string_character_rune")

	// Go doesn't have char type, it uses bytes and rune
	byte_slice := []byte("Où ça? Là!") // UTF-8 encoding
	s := string(byte_slice)
	fmt.Printf("Byte slice: %v\n", byte_slice)                         // [79 195 185 32 195 167 97 63 32 76 195 160 33]
	fmt.Printf("String: %v\n", s)                                      // Où ça? Là!
	fmt.Printf("Bytes count: %v\n", len(byte_slice))                   // 13 also = len(s)
	fmt.Printf("Runes (chars) count: %v\n", utf8.RuneCountInString(s)) // 10

	r := '€' // A rune
	fmt.Printf("r: %d = %c\n", r, r)

	for ix, ru := range s {
		fmt.Printf("[%d] %c\n", ix, ru)
	}
	fmt.Println()
}

func test_unicode() {
	fmt.Println("--- test_unicode")

	sL := "\x99\x00ab\x50\x00\x23\x50\x29\x9c"
	for i := 0; i < len(sL); i++ {
		if unicode.IsPrint(rune(sL[i])) {
			fmt.Printf("%c\n", sL[i])
		} else {
			fmt.Println("Not printable!")
		}
	}
	fmt.Println()
}

func test_strings() {
	fmt.Println("--- test_string")
	var f = fmt.Printf

	// Case insensitive string comparison
	f("EqualFold: %v\n", s.EqualFold("Mihalis", "MIHAlis"))
	f("EqualFold: %v\n", s.EqualFold("Mihalis", "MIHAli"))

	// Instr
	f("Index: %v\n", s.Index("Mihalis", "ha"))
	f("Index: %v\n", s.Index("Mihalis", "Ha"))

	// StartsWith, EndsWith
	f("Prefix: %v\n", s.HasPrefix("Mihalis", "Mi"))
	f("Prefix: %v\n", s.HasPrefix("Mihalis", "mi"))
	f("Suffix: %v\n", s.HasSuffix("Mihalis", "is"))
	f("Suffix: %v\n", s.HasSuffix("Mihalis", "IS"))

	// splits the given string around one or more
	// white space characters as de昀椀ned by the unicode.IsSpace() function and returns
	// a slice of substrings found in the input string
	t := s.Fields("This is a string!")
	f("Fields: %v\n", len(t))
	t = s.Fields("ThisIs a\tstring!")
	f("Fields: %v\n", len(t))

	// strings.Split() kreaks on any separator string
	for _, segment := range s.Split("Je n'aime pas le classique B1 - 09 - Franz Liszt - Rêve d'amour.mp3", " - ") {
		f("Segment: %s\n", segment)
	}
	// No separator breaks string character by character
	for _, car := range s.Split("Où ça?", "") {
		f("Car: %s\n", car)
	}

	// Replace, -1 indicates no limit count
	f("%s\n", s.Replace("Bonjour", "jour", "soir", -1))
	f("%s\n", s.Replace("abcd efg", "", "_", 4))

	// SplitAfter includes separator in output
	f("SplitAfter: %s\n", s.SplitAfter("123++432++", "++"))

	// TrimFunc returns a slice of the string s with all leading and trailing Unicode code points c satisfying f(c) removed.
	trimFunction := func(c rune) bool {
		fmt.Println("c:", c, " ", string(c), " IsLetter:", unicode.IsLetter(c))
		return !unicode.IsLetter(c)
	}
	f("TrimFunc: %s\n", s.TrimFunc("123 abc ABC \t .", trimFunction))

	fmt.Println()
}

func test_time_date() {
	fmt.Println("--- test_time_date")
	// time.Time data type represents an instant in time with nanosecond precision. Each time.Time value is associated
	// with a location (time zone).
	// See https://pkg.go.dev/time

	start := time.Now()

	dateString := "26 February 1965"
	d, err := time.Parse("02 January 2006", dateString)
	if err == nil {
		fmt.Println("Full:", d)
		fmt.Println("Date:", d.Day(), "/", int(d.Month()), "/", d.Year())
	}

	dateString = "07/06/2025 13:02"
	d, err = time.Parse("02/01/2006 15:04", dateString)
	if err == nil {
		fmt.Println("Full:", d)
		fmt.Println("Date:", d.Day(), d.Month(), d.Year())
		fmt.Println("Time:", d.Hour(), d.Minute())
	}

	dateString = "17:25"
	d, err = time.Parse("15:04", dateString)
	if err == nil {
		fmt.Println("Full:", d)
		fmt.Println("Time:", d.Hour(), d.Minute())
	}

	// Back and forth conversion to Unix Epoch time
	t := time.Now().Unix()
	fmt.Println("Epoch time:", t)
	// Convert Epoch time to time.Time
	d = time.Unix(t, 0)
	fmt.Println("Date:", d.Day(), d.Month(), d.Year())
	fmt.Printf("Time: %d:%02d\n", d.Hour(), d.Minute())

	duration := time.Since(start)
	fmt.Println("Execution time:", duration)

	fmt.Println()
}

func test_iota() {
	fmt.Println("--- test_iota")
	const (
		Zero Digit = iota
		One
		Two
		Three
		Four
	)

	fmt.Println(One)
	fmt.Println(Two)

	const (
		p2_0 Power2 = 1 << iota
		_
		p2_2
		_
		p2_4
		_
		p2_6
	)

	fmt.Println("2^0:", p2_0)
	fmt.Println("2^2:", p2_2)
	fmt.Println("2^4:", p2_4)
	fmt.Println("2^6:", p2_6)

	fmt.Println()
}

func test_arrays_slices() {
	fmt.Println("--- test_array")

	// Arrays
	t1 := [4]string{"Once", "upon", "a", "time"}
	t2 := [...]string{"Once", "upon", "a", "time"}
	fmt.Println("t1:", t1)
	fmt.Println("t2:", t2)

	// Slice
	s1 := []string{"Once", "upon", "a", "time"}
	s2 := make([]float64, 3) // Initialized at 0.0 for 3 elements
	fmt.Println("s1:", s1)
	fmt.Println("s2:", s2)

	// 2D slice
	s3 := make([][]int, 2) // Create with two rows
	s3[0] = make([]int, 4) // Fill first row
	s3[1] = make([]int, 2) // Fill second row with a different length
	fmt.Println("s3:", s3)
	// Create and initialize a 2D slice
	s4 := [][]int{{1, 2, 3}, {4, 5, 6}}
	fmt.Println("s4:", s4, '\n')

	// Create an empty slice
	mySlice := []float64{}
	// Both length and capacity are 0 because aSlice is empty
	fmt.Println(mySlice, len(mySlice), cap(mySlice))
	// Add elements to a slice
	mySlice = append(mySlice, 1234.56)
	mySlice = append(mySlice, -34.0)
	fmt.Println(mySlice, "with length", len(mySlice), "with capacity", cap(mySlice), '\n')

	// A slice with length 4
	t := make([]int, 4)
	t[0] = -1
	t[1] = -2
	t[2] = -3
	t[3] = -4
	fmt.Println("t:", ", len:", len(t), ", cap:", cap(t))
	// If append doesn't have enough capacity, it'll double the capacity, so here capacity goes from 4 to 8
	t = append(t, -5)
	fmt.Println("t:", ", len:", len(t), ", cap:", cap(t))
	fmt.Println()

	// A 2D slice
	// You can have as many dimensions as needed
	twoD := [][]int{{1, 2, 3}, {4, 5, 6}}
	// Visiting all elements of a 2D slice
	// with a double for loop
	for _, i := range twoD {
		for _, k := range i {
			fmt.Print(k, " ")
		}
		fmt.Println()
	}
	fmt.Println()

	make2D := make([][]int, 2)
	fmt.Println(make2D)
	make2D[0] = []int{1, 2, 3, 4}
	make2D[1] = []int{-1, -2, -3, -4}
	fmt.Println(make2D)
	fmt.Println()

	// Append and expand
	// Same length and capacity
	aSlice := make([]int, 4, 4)
	fmt.Println(aSlice)
	// This time the capacity of slice aSlice is the same as its length, not because Go decided to do so but because we specified it.
	// Add an element
	aSlice = append(aSlice, 5)
	// When you add a new element to slice aSlice, its capacity is doubled and becomes 8.
	fmt.Println(aSlice)
	// The capacity is doubled
	fmt.Println("len:", len(aSlice), "cap:", cap(aSlice))
	// Now add four elements
	// The ... operator expands []int{-1, -2, -3, -4} into multiple arguments and append() appends each argument one by one to aSlice.
	aSlice = append(aSlice, []int{-1, -2, -3, -4}...)
	fmt.Println(aSlice)
	// The capacity is doubled
	fmt.Println("len:", len(aSlice), "cap:", cap(aSlice))
	fmt.Println()

	cSlice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(cSlice)
	fmt.Println(cSlice[0:5]) // First 5 elements
	fmt.Println(cSlice[:5])  // First 5 elements

	l := len(cSlice)
	fmt.Println(cSlice[l-2 : l]) // Last 2 elements
	fmt.Println(cSlice[l-2:])    // Last 2 elements

	t5 := cSlice[0:5:10]          // First 5 elements
	fmt.Println(len(t5), cap(t5)) // 5 10
	// Elements at indexes 2,3,4
	// Capacity will be 10-2
	t5 = cSlice[2:5:10]
	fmt.Println(len(t), cap(t5)) // 5 8
	// Initially, the capacity of t will be 10-0, which is 10. In the second case, the capacity of t will be 10-2.

	// Elements at indexes 0,1,2,3,4
	// New capacity will be 6-0
	t5 = cSlice[:5:6]
	fmt.Println(len(t), cap(t5))
	fmt.Println()

	t5 = cSlice[0:5:10] // First 5 elements
	t5 = append(t5, 99)
	// t5[9] = 1000		// Not possible: index out of range [9] with length 6
	fmt.Println(t5)
	fmt.Println(cSlice)
	fmt.Println()

	// Deleting an element of a slice
	dSlice := []string{"Il", "était", "un", "petit", "navire"}
	eSlice := append(dSlice[:3], dSlice[4:]...) // The ... operator expands aSlice[i+1:] so that its elements can be appended one by one
	fmt.Println(dSlice)                         // dSlice has been modified!
	fmt.Println(eSlice)
	fmt.Println()

	// copy() function for copying an existing array to a slice or an existing slice to another slice. However, the use
	// of copy() can be tricky because the destination slice is not auto-expanded if the source slice is bigger than the
	// destination slice. Additionally, if the destination slice is bigger than the source slice, then copy() does not
	// empty the elements from the destination slice that did not get copied

	// input1       input2        copy(input1, input2)     copy(input2, input1)
	// 1 2 3 4 5    0 0 0 0       0 0 0 0 5                1 2 3 4
	// 1 2 3 4 5    0 0           0 0 3 4 5                1 2
	// 1 2          0 0           0 0                      1 2

	// Sorting slices
	sInts := []int{1, 0, 2, -3, 4, -20}
	sFloats := []float64{1.0, 0.2, 0.22, -3, 4.1, -0.1}
	sStrings := []string{"aa", "a", "A", "Aa", "aab", "AAa"}
	fmt.Println("sInts original:", sInts)
	sort.Ints(sInts)
	fmt.Println("sInts:", sInts)
	sort.Sort(sort.Reverse(sort.IntSlice(sInts)))
	fmt.Println("Reverse:", sInts)
	// As sort.Interface knows how to sort integers, it is trivial to sort them in reverse
	// order. Sorting in reverse order is as simple as calling the sort.Reverse() function.
	fmt.Println("sFloats original:", sFloats)
	sort.Float64s(sFloats)
	fmt.Println("sFloats:", sFloats)
	sort.Sort(sort.Reverse(sort.Float64Slice(sFloats)))
	fmt.Println("Reverse:", sFloats)
	fmt.Println("sStrings original:", sStrings)
	sort.Strings(sStrings)
	fmt.Println("sStrings:", sStrings)
	sort.Sort(sort.Reverse(sort.StringSlice(sStrings)))
	fmt.Println("Reverse:", sStrings)
	fmt.Println()
}

func test_pointers() {
	fmt.Println("--- test_pointers")
	// Go support pointers, but nor pointer arithmetic

	var f float64 = 12.123
	fmt.Println("Memory address of f:", &f)
	// Pointer to f
	fP := &f
	fmt.Println("Memory address of f:", fP)
	fmt.Println("Value of f:", *fP)

	// The value of f changes
	processPointer(fP)
	fmt.Printf("Value of f: %.2f\n", f)

	// The value of f does not change
	x := returnPointer(f)
	fmt.Printf("Value of x: %.2f\n", *x)

	// The value of f does not change
	xx := bothPointers(fP)
	fmt.Printf("Value of xx: %.2f\n", *xx)
	fmt.Println()

	// Check for empty structure
	var k *aStructure
	// The k variable is a pointer to an aStructure structure. As k points to nowhere, Go
	// makes it point to nil, which is the zero value for pointers.

	// This is nil because currently k points to nowhere
	fmt.Println(k)

	// Therefore you are allowed to do this:
	if k == nil {
		k = new(aStructure)
	}

	// As k is nil, we are allowed to assign it to an empty aStructure value with
	// new(aStructure) without losing any data. Now, k is no longer nil but both fields of
	// aStructure have the zero values of their data types.
	fmt.Printf("%+v\n", k)
	if k != nil {
		fmt.Println("k is not nil!")
	}
	fmt.Println()
}

type aStructure struct {
	field1 complex128
	field2 int
}

func processPointer(x *float64) {
	*x = *x * *x
}

// This is a function that gets a pointer to a float64 variable as input. As we are using
// a pointer, all changes to the function parameter inside the function are persistent.
func returnPointer(x float64) *float64 {
	temp := 2 * x
	return &temp
}

func bothPointers(x *float64) *float64 {
	temp := 2 * *x
	return &temp
}

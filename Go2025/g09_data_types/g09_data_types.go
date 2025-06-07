// g09_data_types.go
// Learning go, data types
//
// 2025-06-06	PV		First version

package main

import (
	"fmt"
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

	byte_slice := []byte("Où ça? Là!")
	s := string(byte_slice)
	fmt.Printf("Byte slice: %v\n", byte_slice)                              // [79 195 185 32 195 167 97 63 32 76 195 160 33]
	fmt.Printf("String: %v\n", s)                                           // Où ça? Là!
	fmt.Printf("Bytes count: %v\n", len(byte_slice))                        // 13 also = len(s)
	fmt.Printf("Characters (runes) count: %v\n", utf8.RuneCountInString(s)) // 10

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
	s3 := make([][]int, 2)	// Create with two rows
	s3[0] = make([]int, 4)	// Fill first row
	s3[1] = make([]int, 2)	// Fill second row with a different length
	fmt.Println("s3:", s3)
	// Create and initialize a 2D slice
	s4 := [][]int{{1, 2, 3}, {4, 5, 6}}
	fmt.Println("s4:", s4)

	fmt.Println()
}

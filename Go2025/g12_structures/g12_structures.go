// g12_structures.go
// Learning go, struct
//
// 2025-06-10	PV		First version

package main

import (
	"fmt"
	"strconv"
)

type Entry struct {
	Name    string
	Surname string
	Year    int
}

func zeroS() Entry {
	return Entry{}
}

// Initialized by the user
func initS(N, S string, Y int) Entry {
	if Y < 2000 {
		return Entry{Name: N, Surname: S, Year: 2000}
	}
	return Entry{Name: N, Surname: S, Year: Y}
}

// Initialized by Go - returns pointer
func zeroPtoS() *Entry {
	t := &Entry{}
	return t
}

func main() {
	fmt.Println("Structures in Go")
	fmt.Println()

	e1 := Entry{"Pierre", "Violent", 1965}
	fmt.Println(e1)

	e2 := zeroS()
	fmt.Println(e2)

	// Allocates a new Entry in memory, fields initialized to zero, and returns pointer
	p3 := new(Entry)
	fmt.Println(*p3)

	e4 := initS("Pierre", "Violent", 1965)
	fmt.Println(e4)

	p5 := zeroPtoS()
	fmt.Println(p5)		// p5, not *p5 for a change

	S := []Entry{}
	for i:=0;i<10;i++ {
		text := "text " +strconv.Itoa(i)
		temp := Entry{Name: text, Surname: "", Year: 2000+i}
		S = append(S, temp)
	}
	fmt.Println(S)
}

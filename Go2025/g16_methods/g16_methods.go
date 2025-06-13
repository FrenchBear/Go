// g16_methods.go
// Learning go, Mothods
//
// 2025-06-13	PV		First version

package main

import (
	"fmt"
	"os"
	"strconv"
)

type ar2x2 [2][2]int

// Traditional Add() function
func Add(a, b ar2x2) ar2x2 {
	c := ar2x2{}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			c[i][j] = a[i][j] + b[i][j]
		}
	}
	return c
}

// Type method Add()
// The ar2x2 variable that called the Add() method is going to be modified and hold the result—this is the reason for
// using a pointer when defining the type method.
func (a *ar2x2) Add(b ar2x2) {
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			a[i][j] = a[i][j] + b[i][j]
		}
	}
}

// Type method Subtract()
func (a *ar2x2) Subtract(b ar2x2) {
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			a[i][j] = a[i][j] - b[i][j]
		}
	}
}

// Type method Multiply()
func (a *ar2x2) Multiply(b ar2x2) {
	a[0][0] = a[0][0]*b[0][0] + a[0][1]*b[1][0]
	a[1][0] = a[1][0]*b[0][0] + a[1][1]*b[1][0]
	a[0][1] = a[0][0]*b[0][1] + a[0][1]*b[1][1]
	a[1][1] = a[1][0]*b[0][1] + a[1][1]*b[1][1]
}

func main() {
	k := [8]int{}
	// Test mode, fixed arguments
	if len(os.Args) == 1 {
		k[0] = 1
		k[1] = 2
		k[2] = 0
		k[3] = 0
		k[4] = 2
		k[5] = 1
		k[6] = 1
		k[7] = 1
	} else {
		if len(os.Args) != 9 {
			fmt.Println("Need 8 integers")
			return
		}
		for index, i := range os.Args[1:] {
			v, err := strconv.Atoi(i)
			if err != nil {
				fmt.Println(err)
				return
			}
			k[index] = v
		}
	}

	a := ar2x2{{k[0], k[1]}, {k[2], k[3]}}
	b := ar2x2{{k[4], k[5]}, {k[6], k[7]}}

	// The main() function gets the input and creates two 2x2 matrices. After that, it
	// performs the desired calculations with these two matrices.
	fmt.Println("Traditional a+b:", Add(a, b))
	a.Add(b)
	fmt.Println("a+b:", a)
	a.Subtract(a)
	fmt.Println("a-a:", a)
	a = ar2x2{{k[0], k[1]}, {k[2], k[3]}}
	a.Multiply(b)
	fmt.Println("a*b:", a)
	a = ar2x2{{k[0], k[1]}, {k[2], k[3]}}
	b.Multiply(a)
	fmt.Println("b*a:", b)
}

// g41_pipe_primes.go
// Learning go, Concurrency, Play with channels, primes sieve using channels
// Totally inefficient but fun, and really simple to write in Go thanks to channels
//
// 2025-07-08	PV		First version

package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("Go Concurrency, Primes sieve")

	n := 100_000
	r := int(math.Ceil(math.Sqrt(float64(n))))

	source := make(chan int)
	go generate(source, n)

	np := 0
	for {
		p, ok := <-source
		if !ok {
			break
		}
		//fmt.Print(p, " ")
		np++

		if p<=r {
			nextsource := make(chan int)
			go filter(source, nextsource, p)
			source = nextsource
		}
	}

	fmt.Printf("\n2..%d: %d primes\n", n, np)
}

func filter(source <-chan int, nextsource chan<- int, p int) {
	for n := range source {
		if n%p != 0 {
			nextsource <- n
		}
	}
	close(nextsource)
}

func generate(source chan int, n int) {
	for i := 2; i <= n; i++ {
		source <- i
	}
	close(source)
}

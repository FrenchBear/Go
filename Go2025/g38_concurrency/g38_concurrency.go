// g38_concurrency.go
// Learning go, Concurrency
//
// 2025-07-06	PV		First version

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("Go Concurrency")

	fmt.Print("You are using ", runtime.Compiler, " ")
	fmt.Println("on a", runtime.GOARCH, "machine")
	fmt.Println("Using Go version", runtime.Version())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Println()

	count := 15

	fmt.Printf("Going to create %d goroutines.\n", count)
	for i := 0; i < count; i++ {
		go func(x int) {
			fmt.Printf("%d ", x)
		}(i)
	}
	// Just wait to be sure all goroutines terminate (end of program doesn't wait)
	time.Sleep(time.Second)
	fmt.Println()

	// Here we use a WaitGroup to wait for proper goroutines terination
	var waitGroup sync.WaitGroup
	fmt.Printf("Going to create %d goroutines.\n", count)
	for i := 0; i < count; i++ {
		waitGroup.Add(1)
		// We call Add(1) just before we create the goroutine in order to avoid race conditions.
		go func(x int) {
			defer waitGroup.Done()
			// The Done() call is going to be executed just before the anonymous function returns because of the defer keyword.
			fmt.Printf("%d ", x)
		}(i)
	}
	waitGroup.Wait()
	fmt.Println()

	fmt.Println("\nExiting...")
}

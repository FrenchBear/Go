// g43_mutex.go
// Learning go, Concurrent programming, Mutexes
//
// 2025-07-09	PV		First version

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Printf("Go Mutexes\n\n")

	wg.Add(2)
	go change()
	go read()
	wg.Wait()
	fmt.Println("main ends")
}

var m sync.Mutex
var v1 int
var wg sync.WaitGroup

func change() {
	fmt.Println("change starts")
	for i := 0; i < 10; i++ {
		m.Lock()
		fmt.Println("change lock")
		time.Sleep(time.Millisecond * time.Duration(347))
		v1++
		fmt.Println("change unlock")
		m.Unlock()
		time.Sleep(time.Millisecond * time.Duration(428))
	}
	fmt.Println("change ends")
	wg.Done()
}

func read() {
	fmt.Println("read starts")
	for i := 0; i < 10; i++ {
		res := m.TryLock()
		if res {
			fmt.Println("read lock immediate")
		} else {
			fmt.Println("read lock delayed")
			m.Lock()
			fmt.Println("read lock")
		}
		fmt.Println("v1:", v1)
		fmt.Println("read unlock")
		m.Unlock()
		time.Sleep(time.Millisecond * time.Duration(682))
	}
	fmt.Println("read ends")
	wg.Done()
}

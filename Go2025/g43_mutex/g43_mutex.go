// g43_mutex.go
// Learning go, Concurrent programming, Mutexes, atomic, semaphores
//
// 2025-07-09	PV		First version

// This program tests "simple" mutex, but there is also sync.RWMutex that allows
// multiple readers and a single writer (writer can't lock until at least a reader has the lock)
// sync.RWMutex use .Lock() and .Unlock() for the writer, and .RLock() and .RUnlock() for readers

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/semaphore"
)

func main() {
	fmt.Println("Go Mutexes, Atomic, Monitor, Semaphores")

	testMutex()
	textAtomic()
	testMonitor()
	testSemaphore()
}

// ======================================================================================================

func testMutex() {
	fmt.Printf("\nTest Mutex\n\n")

	fmt.Println("textMutex starts")
	wg.Add(2)
	go change()
	go read()
	wg.Wait()
	fmt.Println("testMutex ends")
	fmt.Println()
}

var m sync.Mutex
var v1 int
var wg sync.WaitGroup

func change() {
	fmt.Println("change starts")
	for i := 0; i < 10; i++ {
		m.Lock()
		fmt.Println("change lock")
		time.Sleep(time.Millisecond * 437)
		v1++
		fmt.Println("change unlock")
		m.Unlock()
		time.Sleep(time.Millisecond * 428)
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
		time.Sleep(time.Millisecond * 982)
	}
	fmt.Println("read ends")
	wg.Done()
}

// ======================================================================================================

type atomCounter struct {
	val int64
}

func (c *atomCounter) increment() {
	atomic.AddInt64(&c.val, 1)
}

func (c *atomCounter) decrement() {
	atomic.AddInt64(&c.val, -1)
}

func (c *atomCounter) value() int64 {
	return atomic.LoadInt64(&c.val)
}

func textAtomic() {
	fmt.Printf("\nTest Atomic\n\n")

	fmt.Println("testAtomic starts")

	var waitGroup sync.WaitGroup
	counter := atomCounter{}
	for i := 0; i < 10000; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			counter.increment()
		}()
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			counter.decrement()
		}()
	}
	waitGroup.Wait()
	fmt.Println("Final value:", counter.value()) // Could use counter.val since it's single access now

	fmt.Println("testAtomic ends")
	fmt.Println()
}

// ======================================================================================================

// Another way to share memory is to access it through a monitor, that will only manage one request at  a time,
// thus ensuring safe sharing. Actual data is stored in monitor(), here a single integer.

func testMonitor() {
	fmt.Printf("\nTest Monitor\n\n")

	fmt.Println("testMonitor starts")

	// Starts goroutine first
	go monitor()

	var waitGroup sync.WaitGroup
	for i := 0; i < 10; i++ {
		waitGroup.Add(1)
		// Note that values won't be stored in sequence because of random wait before calling set(i)
		go func() {
			defer waitGroup.Done()
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			set(i)
		}()
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			fmt.Println("i:", get())
		}()
	}
	waitGroup.Wait()
	fmt.Println("Final value:", get())

	fmt.Println("testMonitor ends")
	fmt.Println()
}

var readValue = make(chan int)
var writeValue = make(chan int)

func set(newValue int) {
	writeValue <- newValue
}

func get() int {
	return <-readValue
}

func monitor() {
	var value int
	for {
		select {
		case newValue := <-writeValue:
			value = newValue
			fmt.Println("monitor set value to", value)
		case readValue <- value:
		}
	}
}

// ======================================================================================================

// Semaphores can have weights that limit the number of threads or goroutines that can have access to a resource
// The process is supported via the Acquire() and Release() methods, which are defied as follows:
// - func (s *Weighted) Acquire(ctx context.Context, n int64) error
// - func (s *Weighted) Release(n int64)
// The second parameter of Acquire() defines the weight of the semaphore

func testSemaphore() {
	fmt.Printf("\nTest Semaphore\n\n")

	var Workers = 4
	var sem = semaphore.NewWeighted(int64(Workers))

	nJobs := 10
	// Where to store the results
	var results = make([]int, nJobs)

	// Needed by Acquire()
	ctx := context.TODO()
	for i := range results {
		err := sem.Acquire(ctx, 1)
		if err != nil {
			fmt.Println("Cannot acquire semaphore:", err)
			break
		}

		go func(i int) {
			defer sem.Release(1)
			temp := worker(i)
			results[i] = temp
		}(i)
	}

	// This is a clever trick: we acquire all of the tokens so that the sem.Acquire() call
	// blocks until all workers/goroutines have finished. This is similar in functionality to a
	// Wait() call.
	err := sem.Acquire(ctx, int64(Workers))
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range results {
		fmt.Println(k, "->", v)
	}
	fmt.Println()
}

func worker(n int) int {
	square := n * n
	time.Sleep(time.Second)
	return square
}

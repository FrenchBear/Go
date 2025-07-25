// g39_channels.go
// Learning go, Concurrency, Channels
//
// 2025-07-07	PV		First version

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Go Concurrency, Channels")

	example1()
	example2()
	example3()
	example4()
	example5()
	example6()
}

func example1() {
	// Buffered channel with a size of 1
	c := make(chan int, 1)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func(c chan int) {
		defer waitGroup.Done()
		writeToChannel(c, 10) // Will close the channel after writing
		fmt.Println("Exit.")
	}(c)
	fmt.Println("Read:", <-c)

	_, ok := <-c
	if ok {
		fmt.Println("Channel is open!")
	} else {
		fmt.Println("Channel is closed!")
	}

	waitGroup.Wait()
	fmt.Println()
}

func example2() {
	var ch chan bool = make(chan bool) // Unbuffered channel, writing blocks until there's someone to read the value
	for i := 0; i < 5; i++ {
		go printer(ch) // No synchronization
	}

	// Range on channels
	// IMPORTANT: As the channel c is not closed, the range loop does not exit on its own.
	n := 0
	for i := range ch {
		// The range keyword works with channels. However, a range loop on a channel only
		// exits when the channel is closed or using the break keyword.
		fmt.Println(i)
		if i == true {
			n++
		}
		if n > 4 {
			// We close the ch channel when a condition is met and exit the for loop using break.
			fmt.Println("n:", n, " close channel and break the loop")
			close(ch)
			break
		}
	}
	// If the channel was closed before reading 5 booleans, the remaining waiting printer() calls would crash

	x, ok := <-ch
	if ok {
		fmt.Println("Channel is open!")
		fmt.Println(x)
	} else {
		fmt.Println("Channel is closed!")
	}

	// Channel is closed, but reading from it still receive default value from the channel without error...
	// Note that writing on a closed channel cause program to panic
	for i := 0; i < 5; i++ {
		fmt.Println(<-ch)
	}
	fmt.Println()
}

func example3() {
	// Multiple writers
	ci := make(chan int, 1)
	for i := 0; i < 5; i++ {
		go func(i int) {
			ci <- i
		}(i)
	}

	// Multiple readers
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()
			val := <-ci
			fmt.Printf("%d: %d\n", i, val)
		}(i)
	}
	wg.Wait()
	fmt.Println()

	cf := make(chan float64, 10)
	// Note that creating a channel with less than 3 slots will crash, since write can't success and there
	// is noone to read this channel
	cf <- 3.1416
	cf <- 1.4142
	cf <- 2.7181
	close(cf)

	// Reading after channel is closed is Ok and will retrieve values
	fmt.Println(<-cf)
	fmt.Println(<-cf)
	fmt.Println(<-cf)
	fmt.Println(<-cf) // And 0 past the end
	fmt.Println()
}

func example4() {
	// select allows to listen to multiple channels, all channels are examines simultaneously
	// If multiple channels are ready, then select makes a random selection

	// Create two channels
	messageCh := make(chan string)
	timeoutCh := time.After(2 * time.Second) // This channel will receive a value after 2 seconds

	// A goroutine to send a message to messageCh after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		messageCh <- "Hello, from the other side!"
	}()

	// The select statement will wait for one of the cases to be ready
	select {
	case msg := <-messageCh:
		fmt.Println("Received message:", msg)
	case <-timeoutCh:
		fmt.Println("Timeout: Waited for 2 seconds.")
	default:
		fmt.Println("No message received yet. The default case is running.")
		// We can add a small delay here to prevent busy-waiting
		time.Sleep(500 * time.Millisecond)
	}

	// To demonstrate the timeout, we'll have another select that waits longer
	fmt.Println("\n--- Demonstrating Timeout ---")

	// A new channel for this example
	anotherMessageCh := make(chan string)

	// This goroutine will send a message after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		anotherMessageCh <- "This message will be late."
	}()

	// This select will wait for a message or timeout after 2 seconds
	select {
	case msg := <-anotherMessageCh:
		fmt.Println("Received message:", msg)
	case <-time.After(2 * time.Second):
		fmt.Println("Timeout: The operation took too long!")
	}
	fmt.Println()
}

// Write to channel, and immediately close it
func writeToChannel(c chan int, x int) {
	c <- x
	close(c)
}

// The channel parameter specifies the direction, here writing only
// This is optional, func printer(ch chan bool) also works, but it may detect errors with channel misuse
func printer(ch chan<- bool) {
	ch <- true
}

// One channel parameter for reading, one for writing
// func f2(input <-chan int, output chan<- int) {
// 	x := <-input
// 	fmt.Println("Read (f2):", x)
// 	output <- x
// }

// Select can be used when queueing values in a channel that can block if full
// and when dequeueing data when channel can be empty
func example5() {
	numbers := make(chan int, 5)

	counter := 10
	for i := 0; i < counter; i++ {
		select {
		case numbers <- i:
			fmt.Println("Queued", i)
		default:
			fmt.Println("No space for", i)
		}

	}
	fmt.Println()
	for exhausted := false; !exhausted; {
		select {
		case num := <-numbers:
			fmt.Println("Processing", num)
		default:
			fmt.Println("Nothing left to read!")
			exhausted = true
		}
	}
	fmt.Println()
}

// Specifying the orger of execution
var wg sync.WaitGroup

func A(a, b chan struct{}) {
	<-a
	fmt.Println("A()!")
	time.Sleep(time.Second)
	close(b)
}

//Function A() is going to be blocked until channel a, which is passed as a parameter, is closed.
// Just before it ends, it closes channel b, which is passed as a parameter. This is going to unblock the next goroutine, which is going to be function B().
func B(a, b chan struct{}) {
	<-a
	fmt.Println("B()!")
	time.Sleep(3 * time.Second)
	close(b)
}

func C(a, b chan struct{}) {
	<-a
	fmt.Println("C()!")
	close(b)
}

func D(a chan struct{}) {
	<-a
	fmt.Println("D()!")
	wg.Done()
}

func example6() {
	// We need to have as many channels as the number of functions we want to execute as goroutines.
	x := make(chan struct{})
	y := make(chan struct{})
	z := make(chan struct{})
	w := make(chan struct{})

	// This proves that the order of execution dictated by the Go code does not matter as D() is going to be executed last.
	wg.Add(1)
	go func() {
		D(w)
	}()

	wg.Add(1)
	go func() {
		D(w)
	}()
	go A(x, y)
	wg.Add(1)
	go func() {
		D(w)
	}()
	// Although we run C() before B(), C() is going to finish after B() has finished.
	go C(z, w)
	go B(y, z)

	wg.Add(1)
	go func() {
		D(w)
	}()
}

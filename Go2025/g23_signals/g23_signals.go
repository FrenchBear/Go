// g23_signals.go
// Learning go, System programming, signals
//
// Windows adaptation:
// SIGINT = Ctrl+C
// SIGTERM = taskkill /pid <pid>, taskkill /im <imagename.exe>			Doesn't work (ERROR: The process with PID 38844 could not be terminated. Reason: This process can only be terminated forcefully (with /F option).)
// SIGKILL = taskkill /f /pid <pid>, taskkill /f /im <imagename.exe>    Always kill process, no signal caught
//
// 2025-06-22	PV		First version

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleSignal(sig os.Signal) {
	fmt.Println("handleSignal() Caught:", sig)
}

func main() {
	fmt.Printf("Process ID: %d\n", os.Getpid())

	// We create a channel with data of type os.Signal because all channels must have a type.
	sigs := make(chan os.Signal, 1)

	// The previous statement means handle all signals that can be handled.
	signal.Notify(sigs)
	start := time.Now()
	go func() {
		for {
			// Wait until you read data (<-) from the sigs channel and store it in the sig variable.
			sig := <-sigs
			
			// Depending on the read value, act accordingly. This is how you differentiate between signals.
			switch sig {

			case syscall.SIGINT:
				duration := time.Since(start)
				fmt.Println("Execution time:", duration)
			
				// For the handling of syscall.SIGINT, we calculate the time that has passed since the beginning of the
				// program execution, and print it on screen.
			case syscall.SIGTERM:
				handleSignal(sig)
				// The code of the syscall.SIGTERM case calls the handleSignal() function—it is up to the developer to decide on the
				// details of the implementation.
				// Do not use return here because the goroutine exits but the time.Sleep() will continue to work!
				os.Exit(0)
			
			default:
				fmt.Println("Caught:", sig)
			}
		}
	}()

	for {
		fmt.Print("+")
		time.Sleep(10 * time.Second)
	}
}

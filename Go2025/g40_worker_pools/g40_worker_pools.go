// g40_worker_pools.go
// Learning go, Concurrency, Worker Pools
//
// 2025-07-08	PV		First version

package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// The Client structure is used for keeping track of the requests that the program is going to process.
type Client struct {
	id      int
	integer int
}

type Result struct {
	job    Client
	square int
	workerID int
}

var size = runtime.GOMAXPROCS(0)
var clients = make(chan Client, size)
var data = make(chan Result, size)

func worker(workerID int, wg *sync.WaitGroup) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for c := range clients {
		square := c.integer * c.integer
		output := Result{c, square, workerID}
		data <- output

		n := rnd.Intn(1000)
		time.Sleep(time.Millisecond * time.Duration(n))
	}
	wg.Done()
}

func create(n int) {
	for i := 0; i < n; i++ {
		c := Client{i, i}
		clients <- c
	}
	close(clients)
}

func main() {
	fmt.Println("Go Concurrency, Worker Pools")

	if len(os.Args) != 3 {
		fmt.Println("Need #jobs and #workers!")
		return
	}
	nJobs, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	nWorkers, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	// The create() call mimics the client requests that you are going to process.
	go create(nJobs)

	// The finished channel is used for blocking the program and, therefore, needs no particular data type.
	finished := make(chan interface{})

	go func() {
		for d := range data {
			fmt.Printf("Client ID:%3d, Worker Id:%3d   %d² = %d\n", d.job.id, d.workerID, d.job.integer, d.square)
		}

		// The finished <- true statement is used for unblocking the program as soon as
		// the for range loop ends. The for range loop ends when the data channel is closed,
		// which happens after wg.Wait(), which means after all workers have finished.
		finished <- true
	}()

	var wg sync.WaitGroup

	// The purpose of thethis loop is to generate the required number of worker() goroutines to process all requests.
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait()
	close(data)
	fmt.Printf("Finished: %v\n", <-finished)
}

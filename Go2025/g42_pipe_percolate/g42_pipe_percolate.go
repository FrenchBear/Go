// g42_pipe_percolate.go
// Learning go, Concurrency, Play with channels, Percolation using channels
//
// 2025-07-08	PV		First version

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Cell struct {
	IsFilled bool
	North    <-chan bool
	West     <-chan bool
	South    chan bool
	East     chan bool
}

func main() {
	fmt.Println("Go Concurrency, Percolation")
	fmt.Println()

	const SIZE = 10
	const THRESHOLD = 0.4
	source := make(chan bool, SIZE)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create grid of channels
	grid := make([][]Cell, SIZE)
	for r := 0; r < SIZE; r++ {
		grid[r] = make([]Cell, SIZE)
		for c := 0; c < SIZE; c++ {
			var n <- chan bool
			if r == 0 {
				n = source
			} else {
				n = grid[r-1][c].South
			}

			var w <- chan bool
			if c == 0 {
				w = nil
			} else {
				w = grid[r][c-1].East
			}

			var e chan bool
			if c == SIZE-1 {
				e = nil
			} else {
				e = make(chan bool)
			}

			var s chan bool
			s = make(chan bool)

			grid[r][c] = Cell{
				North: n,
				West:  w,
				East:  e,
				South: s,
			}
		}
	}

	// Fill some cells according to density (global fill density algorithm)
	nf := 0
	for {
		r := rnd.Intn(SIZE)
		c := rnd.Intn(SIZE)
		if !grid[r][c].IsFilled {
			grid[r][c].IsFilled = true
			nf++
			if float64(nf)/float64(SIZE*SIZE) >= THRESHOLD {
				break
			}
		}
	}

	// Print grid
	for r := 0; r < SIZE; r++ {
		for c := 0; c < SIZE; c++ {

			if grid[r][c].IsFilled {
				fmt.Printf("X ")
			} else {
				fmt.Printf(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println()

	// Setup percolation grid
	for r := 0; r < SIZE; r++ {
		for c := 0; c < SIZE; c++ {
			go percolate(grid[r][c])
		}
	}

	// Pour water
	for i := 0; i < SIZE; i++ {
		source <- true
	}

	// Look at what get through
	perco := false
	for c := 0; c < SIZE; c++ {
		res_c := <-grid[SIZE-1][c].South
		//fmt.Printf("%d: %v\n", c, res_c)
		if res_c {
			perco = true
		}
	}
	//fmt.Println()

	if perco {
		fmt.Printf("Density %.3f, Percolation\n", THRESHOLD)
	} else {
		fmt.Printf("Density %.3f, No percolation\n", THRESHOLD)

	}
}

func percolate(cell Cell) {
	n := false
	if cell.North != nil {
		n = <-cell.North
	}
	w := false
	if cell.West != nil {
		w = <-cell.West
	}

	res := false
	if !cell.IsFilled {
		res = n || w
	}

	// if !cell.IsFilled {
	// 	fmt.Printf("percolate[%v,%v] -> %v\n", r, c, res)
	// } else {
	// 	res = false
	// 	fmt.Printf("percolate[%v,%v], filled cell -> %v\n", r, c, res)
	// }

	if cell.East != nil {
		cell.East <- res
		close(cell.East)
	}
	if cell.South != nil {
		cell.South <- res
		close(cell.South)
	}
}

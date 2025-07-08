// g42_pipe_percolate.go
// Learning go, Concurrency, Play with channels, Percolation using channels
//
// 2025-07-08	PV		First version

package main

import (
	"fmt"
)

type Cell struct {
	North chan bool
	West  chan bool
	South chan bool
	East  chan bool
}

func main() {
	fmt.Println("Go Concurrency, Percolation")

	const SIZE = 4
	source := make(chan bool, SIZE)

	grid := make([][]Cell, SIZE)
	for r := 0; r < SIZE; r++ {
		grid[r] = make([]Cell, SIZE)
		for c := 0; c < SIZE; c++ {

			var n chan bool
			if r == 0 {
				n = source
			} else {
				n = grid[r-1][c].South
			}

			var w chan bool
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

	// Setup percolation grid
	for r := 0; r < SIZE; r++ {
		for c := 0; c < SIZE; c++ {
			go percolate(r, c, grid[r][c])
		}
	}

	// Pour water
	for i := 0; i < SIZE; i++ {
		source <- true
	}


	// Look at what get through
	for c := 0; c < SIZE; c++ {
		res_c := <-grid[SIZE-1][c].South
		fmt.Printf("%d: %v\n", c, res_c)
	}
}

func percolate(r, c int, cell Cell) {
	n := false
	if cell.North!=nil {
		n = <-cell.North
	}
	e:=false
	if cell.East!=nil {
		e = <-cell.East
	}
	res := n || e

	fmt.Printf("percolate[%v,%v] -> %v\n", r, c, res)


	if cell.West!=nil {
		cell.West<-res
		close(cell.West)
	}
	if cell.South!=nil {
		cell.South<-res
		close(cell.South)
	}
}

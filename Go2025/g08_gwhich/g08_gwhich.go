// g08_gwhich
// Learning go, which utility
// https://go.dev/doc/tutorial/getting-started
//
// 2025-06-04	PV		First version

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide an argument!")
		return
	}
	file := arguments[1]
	path := os.Getenv("PATH")
	fmt.Println(path)
	pathSplit := filepath.SplitList(path)
	for _, directory := range pathSplit {
		fullPath := filepath.Join(directory, file)
		fmt.Println(fullPath)
		// Does it exist?
		fileInfo, err := os.Stat(fullPath)
		if err == nil {
			mode := fileInfo.Mode()
			fmt.Printf("********************* Found: %v  %v\n", mode, mode.IsRegular())
			// Is it a regular file?
			if mode.IsRegular() {
				fmt.Printf("********************* Regular: %v\n", mode)
				// Is it executable?
				if mode&0111 != 0 {
					fmt.Println(fullPath)
					return
				}
			}
		}
	}
}

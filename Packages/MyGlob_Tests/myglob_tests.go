// myglob_tests.go
// Tests for package TestAutoDecode
//
// go mod edit -replace github.com/PieVio/MyGlob=C:\Development\GitHub\Go\Packages\MyGlob
// go mod tidy
//
// 2025-06-23	PV		First version
// 2025-09-07	PV		Test MaxDepth

package main

import (
	"fmt"
	"os"
	"time"

	myglob "github.com/PieVio/MyGlob"
)

func main() {
	fmt.Printf("MyGlob lib version: %s\n\n", myglob.Version())

	//testMyglob(`C:\Development\*.*`, false, []string{"d2"}, 0, 1)
	testMyglob(`S:\MaxDepth`, true, []string{}, 1, 1)
}

func testMyglob(pattern string, autorecurse bool, ignoreDirs []string, maxDepth int, loops int) {
	var durations []float64
	for pass := 0; pass < loops; pass++ {
		fmt.Printf("\nTest #%d\n", pass)

		start := time.Now()
		builder := myglob.New(pattern).Autorecurse(autorecurse).MaxDepth(maxDepth)
		for _, ignoreDir := range ignoreDirs {
			builder.AddIgnoreDir(ignoreDir)
		}
		gs, err := builder.Compile()
		if err != nil {
			fmt.Printf("Error building MyGlob: %s\n", err)
			return
		}

		nf := 0
		nd := 0
		for ma := range gs.Explore() {
			if ma.Err != nil {
				fmt.Println(ma.Err)
				continue
			}
			if ma.IsDir {
				fmt.Printf("%s%c\n", ma.Path, os.PathSeparator)
				nd++
			} else {
				fmt.Println(ma.Path)
				nf++
			}
		}
		duration := time.Since(start)
		fmt.Printf("%d file(s) found\n", nf)
		fmt.Printf("%d dir(s) found\n", nd)
		fmt.Printf("Iterator search in %.3fs\n\n", duration.Seconds())
		durations = append(durations, duration.Seconds())
	}
}

// g49_rotor_router.go
// Learning go, Generating a png image
// From Complexités, JP Delahaye, P.232, 236
//
// 2025-09-16	PV		First version in Go

// WOTAN on 2025-09-16: Done for 70000 iterations in 11.433s

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
)

func main() {
	const (
		r = 150
		n = 70_000		// max π.r² iterations
	)

	col0 := color.RGBA{255, 0, 0, 255}
	col1 := color.RGBA{0, 0, 255, 255}
	col2 := color.RGBA{0, 255, 0, 255}
	col3 := color.RGBA{0, 255, 255, 255}

	width := 2*r+1
	height := 2*r+1

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	start := time.Now()
	for i := 0; i < n; i++ {
		x:=r
		y:=r
InnerLoop:
		for {
			switch img.At(x,y) {
				case col0:
					img.Set(x,y,col1)
					y++
				case col1:
					img.Set(x,y,col2)
					x++
				case col2:
					img.Set(x,y,col3)
					y--
				case col3:
					img.Set(x,y,col0)
					x--
				default:
					img.Set(x,y,col0)
					break InnerLoop
			}
		}
	}
	duration := time.Since(start)

	f, err := os.Create("rotor_router.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatalf("png.Encode failed: %v", err)
	}

	fmt.Printf("Done for %d iterations in %.3fs\n", n, float64(duration.Milliseconds())/1000.0)
}

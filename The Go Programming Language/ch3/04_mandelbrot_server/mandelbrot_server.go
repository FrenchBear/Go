// mandelbrot_server
// http version of mandelbrot program
// Emits a PNG image of the Mandelbrot fractal.
//
// http://localhost:8000/?width=800&ss=0
//
// 2019-10-08	PV

package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		}

		width := 1024
		if widthStr, found := r.Form["width"]; found {
			if v, err := strconv.Atoi(widthStr[0]); err == nil && v >= 10 && v <= 2000 {
				width = v
			} else {
				log.Printf("Invalid width argument %s, must be int between 10 and 2000", widthStr)
			}
		}
		// Supersampling activated by default (4 calculations per pixel)
		ss := true
		if ssStr, found := r.Form["ss"]; found {
			if strings.ToLower(ssStr[0]) == "false" || ssStr[0] == "0" {
				ss = false
			}
		}

		surface(w, width, ss)
	}

	http.HandleFunc("/", handler) // Each request call handler
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func surface(out io.Writer, w int, ss bool) {
	const (
		xmin, ymin, xmax, ymax = -2, -1.25, +0.5, +1.25
	)

	width, height := w, w

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	halfPixelX := float64(xmax-xmin) / (2.0 * float64(width))
	halfPixelY := float64(ymax-ymin) / (2.0 * float64(height))
	for py := 0; py < height; py++ {
		y := float64(py)/float64(height)*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/float64(width)*(xmax-xmin) + xmin

			if ss {
				// Supersampling
				var v1 uint = uint(mandelbrot(complex(x, y)))
				var v2 uint = uint(mandelbrot(complex(x+halfPixelX, y)))
				var v3 uint = uint(mandelbrot(complex(x, y+halfPixelY)))
				var v4 uint = uint(mandelbrot(complex(x+halfPixelX, y+halfPixelY)))
				img.Set(px, py, color.Gray{uint8((v1 + v2 + v3 + v4) >> 2)})
			} else {
				// Simple calculation
				z := complex(x, y)
				img.Set(px, py, color.Gray{mandelbrot(z)})
			}
		}
	}
	png.Encode(out, img) // NOTE: ignoring errors
}

// Returns level of gray 0..255, 0=black
func mandelbrot(z complex128) uint8 {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return uint8(255 - contrast*n) // Result modulo 256 -> several cycles of color
		}
	}
	return 0
}

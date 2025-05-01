// server4_lissajous.go
// Learning Go, Th Gro Programming Language, chap 1.7
// Minimal web server showing an animated gif
// http://localhost:8000/?cycles=20
//
// 2019-09-26	PV

package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		}

		cycles := 5
		if cyclesStr, found := r.Form["cycles"]; found {
			if v, err := strconv.Atoi(cyclesStr[0]); err == nil && v > 0 && v <= 50 {
				cycles = v
			} else {
				log.Print("Invalid cycles argument, must be int between 1 and 50")
			}
		}

		lissajous(w, cycles)
	}

	http.HandleFunc("/", handler) // Each request call handler
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

var palette = []color.Color{color.Black, color.RGBA{0x00, 0x80, 0x00, 0xFF}}

const (
	blackIndex = 0 // in palette
	penIndex   = 1
)

func lissajous(out io.Writer, cycles int) {
	const (
		res     = 0.001 // angular resolution
		size    = 200   // image canvas covers [-size..size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), penIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // ignore errors
}

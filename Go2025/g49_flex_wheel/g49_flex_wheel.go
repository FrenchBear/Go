// g49_flex_wheel.c
// Adaptation of a C RPi program (fb4wheel.c) comparing lines drawing with/without anti-aliasing
// From raspberrycompote.blogspot.com/2014/04/low-level-graphics-on-raspberry-pi.html, Code from fbtestXX.c and fbtest6.c
// Compare Bresenham line drawing (non-aliasing) and  Xiaolin Wu line (aliasing) algorithms
//
// 2016-06-05	PV		Adapted fbtestXX.c to support all depths and not only 8-bit palette
// 2016-06-07	PV		Retrieve existing pixel colog for smooth blending in final version
// 2025-09-19	PV 		Go translation by Gemini

package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math"
	"os"
)

const (
	width            = 1600
	height           = 800
	numFrames        = 100
	fps              = 30
	numPaletteColors = 256
	spokes           = 72
)

var (
	// Background color (black)
	bgColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	// Foreground color (light green)
	fgColor = color.RGBA{R: 84, G: 255, B: 84, A: 255}
)

func main() {
	// Create a color palette with a gradient from background to foreground
	palette := createPalette(bgColor, fgColor, numPaletteColors)

	// Create the animated GIF structure
	anim := &gif.GIF{}

	// Generate each frame
	for i := 0; i < numFrames; i++ {
		// Create a new paletted image for the frame
		img := image.NewPaletted(image.Rect(0, 0, width, height), palette)
		// Fill background
		draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

		// Calculate animation phase for rotation (3 spokes roration over the whole frames)
		phase := (3*float64(i)/float64(spokes) / float64(numFrames)) * 2 * math.Pi

		// Define wheel parameters
		radius := 0.45 * math.Min(float64(height), float64(width)/2)
		xc1 := width / 4
		xc2 := 3 * width / 4
		yc := height / 2

		// Draw the two wheels
		for j := 0; j < spokes; j++ {
			angle := (float64(j) / float64(spokes)) * 2 * math.Pi
			endX := math.Cos(angle+phase)*radius + 0.5
			endY := math.Sin(angle+phase)*radius + 0.5

			// Left wheel: Bresenham's algorithm (no anti-aliasing)
			// The foreground color is the last one in the palette.
			drawLineBresenham(img, xc1, yc, xc1+int(endX), yc+int(endY), uint8(numPaletteColors-1))

			// Right wheel: Xiaolin Wu's algorithm (anti-aliasing)
			drawLineWu(img, xc2, yc, xc2+int(endX), yc+int(endY), fgColor, bgColor, palette)
		}

		// Add the completed frame to the GIF
		anim.Image = append(anim.Image, img)
		anim.Delay = append(anim.Delay, 100/fps) // Delay in 1/100s of a second
	}

	// Save the animated GIF
	f, err := os.Create("flexwheel.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gif.EncodeAll(f, anim)
}

// createPalette generates a palette with a linear gradient.
func createPalette(start, end color.RGBA, steps int) color.Palette {
	p := make(color.Palette, steps)
	for i := 0; i < steps; i++ {
		ratio := float64(i) / float64(steps-1)
		r := uint8(float64(start.R) + ratio*(float64(end.R)-float64(start.R)))
		g := uint8(float64(start.G) + ratio*(float64(end.G)-float64(start.G)))
		b := uint8(float64(start.B) + ratio*(float64(end.B)-float64(start.B)))
		p[i] = color.RGBA{R: r, G: g, B: b, A: 255}
	}
	return p
}

// findClosestColor finds the index of the color in the palette that is closest to c.
func findClosestColor(c color.RGBA, p color.Palette) uint8 {
	closestIndex := 0
	minDist := int(^uint(0) >> 1) // Max integer

	for i, palColor := range p {
		r, g, b, _ := palColor.RGBA()
		dr := int(r>>8) - int(c.R)
		dg := int(g>>8) - int(c.G)
		db := int(b>>8) - int(c.B)
		dist := dr*dr + dg*dg + db*db
		if dist < minDist {
			minDist = dist
			closestIndex = i
		}
	}
	return uint8(closestIndex)
}

// drawLineBresenham draws a line using Bresenham's algorithm (no anti-aliasing) in given color (index in palette)
func drawLineBresenham(img *image.Paletted, x0, y0, x1, y1 int, colorIndex uint8) {
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}

	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		if x0 < 0 || x0 >= width || y0 < 0 || y0 >= height {
			// bounds check
		} else {
			img.SetColorIndex(x0, y0, colorIndex)
		}

		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// plotWu draws a single anti-aliased pixel.
func plotWu(img *image.Paletted, x, y int, brightness float64, fg, bg color.RGBA, p color.Palette) {
	if x < 0 || x >= width || y < 0 || y >= height {
		return // bounds check
	}
	// Linear interpolation between background and foreground color
	r := float64(bg.R) + float64(fg.R-bg.R)*brightness
	g := float64(bg.G) + float64(fg.G-bg.G)*brightness
	b := float64(bg.B) + float64(fg.B-bg.B)*brightness
	
	// Find the best-matching color in the palette and draw the pixel
	blendedColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	colorIndex := findClosestColor(blendedColor, p)
	img.SetColorIndex(x, y, colorIndex)
}

// drawLineWu draws a line using Xiaolin Wu's algorithm (with anti-aliasing).
// Basic implementation from https://en.wikipedia.org/wiki/Xiaolin_Wu%27s_line_algorithm
func drawLineWu(img *image.Paletted, x0, y0, x1, y1 int, fg, bg color.RGBA, p color.Palette) {
	dx := float64(x1 - x0)
	dy := float64(y1 - y0)

	swapped := false
	if math.Abs(dx) < math.Abs(dy) {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
		dx, dy = dy, dx
		swapped = true
	}

	if x1 < x0 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	gradient := dy / dx
	if dx == 0 {
		gradient = 1.0
	}

	// plot function that handles swapping
	plot := func(x, y int, b float64) {
		if swapped {
			plotWu(img, y, x, b, fg, bg, p)
		} else {
			plotWu(img, x, y, b, fg, bg, p)
		}
	}

	// handle first endpoint
	xend := round(float64(x0))
	yend := float64(y0) + gradient*(xend-float64(x0))
	xgap := rfpart(float64(x0) + 0.5)
	xpxl1 := int(xend)
	ypxl1 := int(yend)
	plot(xpxl1, ypxl1, rfpart(yend)*xgap)
	plot(xpxl1, ypxl1+1, fpart(yend)*xgap)
	intery := yend + gradient

	// handle second endpoint
	xend = round(float64(x1))
	yend = float64(y1) + gradient*(xend-float64(x1))
	xgap = fpart(float64(x1) + 0.5)
	xpxl2 := int(xend)
	ypxl2 := int(yend)
	plot(xpxl2, ypxl2, rfpart(yend)*xgap)
	plot(xpxl2, ypxl2+1, fpart(yend)*xgap)

	// main loop
	for x := xpxl1 + 1; x < xpxl2; x++ {
		plot(x, int(intery), rfpart(intery))
		plot(x, int(intery)+1, fpart(intery))
		intery = intery + gradient
	}
}

// Xiaolin Wu's algorithm helper functions
//func ipart(x float64) float64     { return math.Floor(x) }
func round(x float64) float64     { return math.Round(x) }
func fpart(x float64) float64     { return x - math.Floor(x) }
func rfpart(x float64) float64    { return 1.0 - fpart(x) }

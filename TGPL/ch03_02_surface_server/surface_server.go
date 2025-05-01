// surface_server.go
// Computes an SVG rendering of a 3-D surface function in an http server
//
// 2019-09-07	PV

package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	width, height = 1200, 640           // canvas size in pixels
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var cells = 100 // number of grid cells

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
		}

		cells = 100
		if cyclesStr, found := r.Form["cells"]; found {
			if v, err := strconv.Atoi(cyclesStr[0]); err == nil && v >= 10 && v <= 500 {
				cells = v
			} else {
				log.Print("Invalid cells argument, must be int between 10 and 500")
			}
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		surface(w)
	}

	http.HandleFunc("/", handler) // Each request call handler
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)
func surface(out io.Writer) {
	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	var color string
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az := corner(i+1, j)
			bx, by, bz := corner(i, j)
			cx, cy, cz := corner(i, j+1)
			dx, dy, dz := corner(i+1, j+1)
			z := (az + bz + cz + dz) / 4
			if z >= 0 {
				color = fmt.Sprintf("#%02x0000", normalize(z))
			} else {
				color = fmt.Sprintf("#0000%02x", normalize(z))
			}
			fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g' style='fill:%s' />\n",
				ax, ay, bx, by, cx, cy, dx, dy, color)
		}
	}
	fmt.Fprintln(out, "</svg>")
}

func normalize(z float64) int {
	i := int(1000.0 * math.Abs(z))
	if i > 255 {
		i = 255
	}
	return i
}

func corner(i, j int) (float64, float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/float64(cells) - 0.5)
	y := xyrange * (float64(j)/float64(cells) - 0.5)

	// Compute surface height z.
	z := f3(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z
}

// Amortized wave
func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

// eggbox
func f2(x, y float64) float64 {
	//r := math.Hypot(x, y) // distance from (0,0)
	return -math.Abs(math.Sin(x/5.0) * math.Sin(y/5.0))
}

// hills
func f3(x, y float64) float64 {
	//r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(x/3.5) * math.Sin(y/3.5) / 3.0
}

// xkcd2 - Learning go
// Download xkcd into in json format and use a template to produce HTML output
//
// 2019-10-25	PV

package main

import (
	"html/template"
	"log"
	"os"
	"strconv"

	"go.pv/xkcd2/xkcd"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: xkcd2 number")
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Usage: xkcd2 number")
	}

	result, err := xkcd.SearchComic(n)
	if err != nil {
		log.Fatal(err)
	}

	var webPage = template.Must(template.New("issuelist").Parse(`
<h1>{{.Num}}: {{.Title}} issues</h1>
<img src="{{.Img}}"/></p>
{{.Alt}}
`))

	if err := webPage.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}

}

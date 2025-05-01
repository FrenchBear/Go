// Learning go
// kxcd viewer: download xkcd into in json format and shows detail
//
// 2019-10-23	PV

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"go.pv/learning/xkcd"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: xkcdinfo number")
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Usage: xkcdinfo number")
	}

	fmt.Printf("xkcd %d\n", n)
	result, err := xkcd.SearchComic(n)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", *result)
}

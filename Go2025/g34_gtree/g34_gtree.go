// g34_gtree.go
// Learning go, System programming, Tree command in Go
//
// 2025-06-29	PV		First version

// ToDo: Skip well hidden files (SYSTEM+DIRECTORY)
// ToDo: Implement showall option ?

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Options
var h1, h2 bool
var verbose bool
var showall bool

// Global constants
const APP_NAME string = "rtree"
const APP_VERSION string = "1.0.0"

// Usage overrides default flag version
func Usage() {
	fmt.Printf("%s %s\nVisual directory structure in Go\n\n", APP_NAME, APP_VERSION)
	fmt.Printf("Usage: %s [-?|-h] [-v] [-a] directory\nOptions", APP_NAME)
	fmt.Println("")
	flag.PrintDefaults()
}

type DataBag = struct {
	DirsCount  int
	LinksCount int
}

func main() {
	flag.BoolVar(&h1, "h", false, "Shows this message")
	flag.BoolVar(&h2, "?", false, "Shows this message")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&showall, "a", false, "Include hidden directories")

	flag.Usage = Usage
	flag.Parse()

	// First process help
	if h1 || h2 || flag.NArg() > 0 && (flag.Args()[0] == "?" || flag.Args()[0] == "help") {
		flag.Usage()
		os.Exit(0)
	}

	var root string
	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	} else if flag.NArg() == 1 {
		root = flag.Args()[0]
	} else {
		root = "."
	}

	b := DataBag{}
	start := time.Now()
	do_print(&b, root)

	duration := time.Since(start)
	if verbose {
		fmt.Printf("%s: %d directorie(s)", APP_NAME, b.DirsCount)
		if b.LinksCount > 0 {
			fmt.Printf(", %d link(s)", b.LinksCount)
		}
		fmt.Printf(" in %.3fs\n", duration.Seconds())
	}
}

func do_print(b *DataBag, root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a valid directory: %v:\n", APP_NAME, root, err)
		return
	}
	infos := make([]fs.FileInfo, 0)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Errror processing '%s' entry: %s\n", APP_NAME, root, err)
		} else if info.IsDir() {
			infos = append(infos, info)
		}
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name() < infos[j].Name()
	})

	fmt.Println(root)

	for i, info := range infos {
		print_tree(b, root, info.Name(), "", i == len(infos)-1)
	}
}

func print_tree(b *DataBag, root string, subdir string, prefix string, is_last bool) {
	var entry_prefix string
	var new_prefix string
	if is_last {
		entry_prefix = "└── "
		new_prefix = prefix + "    "
	} else {
		entry_prefix = "├── "
		new_prefix = prefix + "│   "
	}

	b.DirsCount++
	fmt.Println(prefix + entry_prefix + subdir)

	subdir_fp := filepath.Join(root, subdir)
	entries, err := os.ReadDir(subdir_fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a valid directory: %v:\n", APP_NAME, subdir_fp, err)
		return
	}

	infos := make([]fs.FileInfo, 0)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Errror processing '%s' entry: %s\n", APP_NAME, root, err)
		} else if info.IsDir() {
			infos = append(infos, info)
		}
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name() < infos[j].Name()
	})

	for i, info := range infos {
		print_tree(b, subdir_fp, info.Name(), new_prefix, i == len(infos)-1)
	}
}

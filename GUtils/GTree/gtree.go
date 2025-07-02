// gtree.go
// Tree command in Go
// without -a option, doesn't show hidden directories or directories starting with a dot
// SYSTEM+HIDDEN directories are always skipped
// For now, the code in Windows-only until I learn how to compile code conditionally
//
// 2025-06-29	PV		First version
// 2025-06-30	PV		Sort files and folders using StrCmpLogicalW used by Windows file explorer
// 2025-07-02	PV		1.2.0 Separate Linux and Windows code

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
const APP_NAME string = "gtree"
const APP_VERSION string = "1.2.0"

// usage overrides default flag version
func usage() {
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
	flag.BoolVar(&showall, "a", false, "Show all directories, including hidden directories and directories starting with a dot")

	flag.Usage = usage
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
	doPrint(&b, root)

	duration := time.Since(start)
	if verbose {
		fmt.Printf("%d directorie(s)", b.DirsCount)
		if b.LinksCount > 0 {
			fmt.Printf(", %d link(s)", b.LinksCount)
		}
		fmt.Printf(" in %.3fs\n", duration.Seconds())
	}
}



func doPrint(b *DataBag, root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a valid directory: %v\n", APP_NAME, root, err)
		return
	}
	infos := make([]fs.FileInfo, 0)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Errror processing '%s' entry: %s\n", APP_NAME, root, err)
			continue
		} else if info.IsDir() {
			h, sh := is_hidden_folder(filepath.Join(root, info.Name()))
			if sh || h && !showall {
				continue
			}
			infos = append(infos, info)
		}
	}
	sort.Slice(infos, func(i, j int) bool {
		return path_comparer(infos[i].Name(), infos[j].Name()) < 0
	})

	fmt.Println(root)

	for i, info := range infos {
		printTree(b, root, info.Name(), "", i == len(infos)-1)
	}
}

func printTree(b *DataBag, root string, subdir string, prefix string, is_last bool) {
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
			h, sh := is_hidden_folder(filepath.Join(subdir_fp, info.Name()))
			// Ignore well-hidden directories such as $RECYCLE.BIN
			if sh || h && !showall {
				continue
			}
			infos = append(infos, info)
		}
	}
	sort.Slice(infos, func(i, j int) bool {
		return path_comparer(infos[i].Name(), infos[j].Name()) < 0
	})

	for i, info := range infos {
		printTree(b, subdir_fp, info.Name(), new_prefix, i == len(infos)-1)
	}
}

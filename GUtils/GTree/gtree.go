// gtree.go
// Tree command in Go
// without -a option, doesn't show hidden directories or directories starting with a dot
// SYSTEM+HIDDEN directories are always skipped
// For now, the code in Windows-only until I learn how to compile code conditionally
//
// 2025-06-29	PV		First version
// 2025-06-30	PV		Sort files and folders using StrCmpLogicalW used by Windows file explorer
// 2025-07-02	PV		1.2.0 Separate Linux and Windows code
// 2025-07-02	PV		1.2.1 Usage using MyMarkup
// 2025-07-02	PV		1.2.2 Print links
// 2025-07-03	PV		1.3.0 Junctions, use sortmethod, maxdepth

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	MyMarkup "github.com/PieVio/MyMarkup"
)

// Options
var h1, h2 bool
var verbose bool
var show_hidden bool
var show_hidden_and_system bool
var sortmethod int
var sortmethod1 bool
var sortmethod2 bool
var maxdepth int

// Global constants
const APP_NAME string = "gtree"
const APP_VERSION string = "1.3.0"

// usage overrides default flag version
func usage() {
	fmt.Printf("%s %s\nVisual directory structure in Go\n\n", APP_NAME, APP_VERSION)

	text := `⌊Usage⌋: {APP_NAME} ¬[⦃?⦄|⦃-?⦄|⦃-h⦄] [-⦃a⦄|-⦃A⦄] [-⦃s⦄ ⦃0⦄|⦃1⦄|⦃2⦄] [⦃-d⦄ ⟨max_depth⟩] [-⦃v⦄] [⟨dir⟩]

⌊Options⌋:
⦃?⦄|⦃-?⦄|⦃-h⦄      ¬Show this message
⦃-a⦄           ¬Show hidden directories and directories starting with a dot
⦃-A⦄           ¬Show system+hidden directories and directories starting with a dollar sign
⦃-s⦄ ⦃0⦄|⦃1⦄|⦃2⦄     ¬Sort method: 0=Default, 1=Windows File Explorer (Windows only), 2=Case fold
⦃-d⦄ ⟨max_depth⟩ ¬Limits recursion to max_depth folders, default is 0 meaning no limitation
⦃-v⦄           ¬Verbose output
⟨dir⟩          ¬Starting directory`

	MyMarkup.RenderMarkup(strings.Replace(text, "{APP_NAME}", APP_NAME, -1))
}

type DataBag = struct {
	DirCount      int
	SymLinkDCount int
	JunctionCount int
}

// func main() {
// 	h, s := is_hidden_folder(`C:\Users\Pierr\Cookies`)
// 	fmt.Printf("h=%v, s=%v\n", h, s)
// }

func main() {
	flag.BoolVar(&h1, "h", false, "Shows this message")
	flag.BoolVar(&h2, "?", false, "Shows this message")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&show_hidden, "a", false, "Show all (hidden) directories")
	flag.BoolVar(&show_hidden_and_system, "A", false, "Show all (hidden and hidden+system) directories")
	flag.IntVar(&sortmethod, "s", 0, "Sort method: 0=Default, 1=Windows File Explorer, 2=Case fold")
	flag.BoolVar(&sortmethod1, "s1", false, "Sort method 1")
	flag.BoolVar(&sortmethod2, "s2", false, "Sort method 2")
	flag.IntVar(&maxdepth, "d", 0, "Max recursion depth, 0=no limit")

	flag.Usage = usage
	flag.Parse()

	// First process help
	if h1 || h2 || flag.NArg() > 0 && (flag.Args()[0] == "?" || flag.Args()[0] == "help") {
		flag.Usage()
		os.Exit(0)
	}

	if sortmethod1 {
		sortmethod = 1
	}
	if sortmethod2 {
		sortmethod = 2
	}

	// show_hidden_and_system implies show_hidden
	if show_hidden_and_system {
		show_hidden = true
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
	doPrint(&b, root, maxdepth)

	duration := time.Since(start)
	if verbose {
		fmt.Printf("%d subdirectorie(s)", b.DirCount)
		if b.SymLinkDCount > 0 {
			fmt.Printf(", %d SymLinkD(s)", b.SymLinkDCount)
		}
		if b.JunctionCount > 0 {
			fmt.Printf(", %d Junction(s)", b.JunctionCount)
		}
		fmt.Printf(" in %.3fs\n", duration.Seconds())
	}
}

type DirEntryType int

const (
	DET_Dir DirEntryType = iota
	DET_SymLinkDcujurs
	DET_Junction
)

type DirEntryData struct {
	Type   DirEntryType
	Name   string
	Target string
}

func doPrint(b *DataBag, root string, maxdepth int) {
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", APP_NAME, err)
		return
	}
	infos := make([]DirEntryData, 0)
	for _, entry := range entries {
		fp := filepath.Join(root, entry.Name())
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Errror processing '%s' entry: %s\n", APP_NAME, root, err)
			continue
		} else {
			h, s := is_hidden_folder(fp)

			// fmt.Printf("%s  h=%v, s=%v\n", fp, h, s)

			if s && !show_hidden_and_system || h && !show_hidden {
				continue
			}

			if info.IsDir() {
				infos = append(infos, DirEntryData{Type: DET_Dir, Name: info.Name()})
			} else if ok, target, _ := IsSymLinkD(fp); ok {
				infos = append(infos, DirEntryData{Type: DET_SymLinkD, Name: info.Name(), Target: target})
			} else if ok, target, _ := IsJunction(fp); ok {
				infos = append(infos, DirEntryData{Type: DET_Junction, Name: info.Name(), Target: target})
			}
		}
	}
	sort.Slice(infos, func(i, j int) bool {
		return path_comparer(sortmethod, infos[i].Name, infos[j].Name) < 0
	})

	fmt.Println(root)

	for i, info := range infos {
		printTree(b, root, info, "", i == len(infos)-1, maxdepth-1)
	}
}

func printTree(b *DataBag, root string, subdir DirEntryData, prefix string, is_last bool, depth int) {
	var entry_prefix string
	var new_prefix string
	if is_last {
		entry_prefix = "└── "
		new_prefix = prefix + "    "
	} else {
		entry_prefix = "├── "
		new_prefix = prefix + "│   "
	}

	fmt.Print(prefix + entry_prefix)
	fmt.Print(subdir.Name)

	switch subdir.Type {
	case DET_SymLinkD:
		fmt.Println(" ->", subdir.Target, " [SymLinkD]")
		b.SymLinkDCount++
		return
	case DET_Junction:
		fmt.Println(" -> ", subdir.Target, " [Junction]")
		b.JunctionCount++
		return
	}

	b.DirCount++

	subdir_fp := filepath.Join(root, subdir.Name)
	entries, err := os.ReadDir(subdir_fp)
	if err != nil {
		fmt.Println("  ... ?")
		return
	}

	infos := make([]DirEntryData, 0)
	for _, entry := range entries {
		fp := filepath.Join(root, entry.Name())
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Errror processing '%s' entry: %s\n", APP_NAME, root, err)
		} else {
			h, s := is_hidden_folder(fp)
			if s && !show_hidden_and_system || h && !show_hidden {
				continue
			}

			if info.IsDir() {
				infos = append(infos, DirEntryData{Type: DET_Dir, Name: info.Name()})
			} else if ok, target, _ := IsSymLinkD(fp); ok {
				infos = append(infos, DirEntryData{Type: DET_SymLinkD, Name: info.Name(), Target: target})
			} else if ok, target, _ := IsJunction(fp); ok {
				infos = append(infos, DirEntryData{Type: DET_Junction, Name: info.Name(), Target: target})
			}
		}
	}

	if depth == 0 {
		if len(infos) > 0 {
			fmt.Println(" ...")
		} else {
			fmt.Println()
		}
		return
	}
	fmt.Println("")

	sort.Slice(infos, func(i, j int) bool {
		return path_comparer(sortmethod, infos[i].Name, infos[j].Name) < 0
	})

	for i, info := range infos {
		printTree(b, subdir_fp, info, new_prefix, i == len(infos)-1, depth-1)
	}
}

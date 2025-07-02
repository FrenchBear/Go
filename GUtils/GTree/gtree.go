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

// ToDo: use sortmethod
// ToDo: hidden links
// ToDo: Show links, hidden and well-hidden folders in color
// ToDo: Show development reparse data
// ToDo: Option -d max_depth
// ToDo: Process junctions (in C:\)

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
var showall bool
var sortmethod int
var sortmethod1 bool
var sortmethod2 bool

// Global constants
const APP_NAME string = "gtree"
const APP_VERSION string = "1.2.2"

// usage overrides default flag version
func usage() {
	fmt.Printf("%s %s\nVisual directory structure in Go\n\n", APP_NAME, APP_VERSION)
	// fmt.Printf("Usage: %s [-?|-h] [-v] [-a] directory\nOptions", APP_NAME)
	// fmt.Println("")

	text := `⌊Usage⌋: {APP_NAME} ¬[⦃?⦄|⦃-?⦄|⦃-h⦄] [-⦃a⦄] [-⦃s⦄ ⦃0⦄|⦃1⦄|⦃2⦄] [-⦃v⦄] [⟨dir⟩]

⌊Options⌋:
⦃?⦄|⦃-?⦄|⦃-h⦄  ¬Show this message
⦃-a⦄       ¬Show all directories, including hidden directories and directories starting with a dot
⦃-s⦄ ⦃0⦄|⦃1⦄|⦃2⦄ ¬Sort method: 0=Default, 1=Windows File Explorer (Windows only), 2=Case fold
⦃-v⦄       ¬Verbose output
⟨dir⟩      ¬Starting directory`

	MyMarkup.RenderMarkup(strings.Replace(text, "{APP_NAME}", APP_NAME, -1))
}

type DataBag = struct {
	DirsCount  int
	LinksCount int
}

func main() {
	flag.BoolVar(&h1, "h", false, "Shows this message")
	flag.BoolVar(&h2, "?", false, "Shows this message")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&showall, "a", false, "Show all (hidden) directories")
	flag.IntVar(&sortmethod, "s", 0, "Sort method: 0=Default, 1=Windows File Explorer, 2=Case fold")
	flag.BoolVar(&sortmethod1, "s1", false, "Sort method 1")
	flag.BoolVar(&sortmethod2, "s2", false, "Sort method 2")

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

	var root string
	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	} else if flag.NArg() == 1 {
		root = flag.Args()[0]
	} else {
		root = "."
	}

	b := DataBag{DirsCount: 1, LinksCount: 0}
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

type DirLink struct {
	IsLink  bool
	DirName string
	Target  string
}

func doPrint(b *DataBag, root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a valid directory: %v\n", APP_NAME, root, err)
		return
	}
	infos := make([]DirLink, 0)
	for _, entry := range entries {
		fp := filepath.Join(root, entry.Name())
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Errror processing '%s' entry: %s\n", APP_NAME, root, err)
			continue
		} else if info.IsDir() {
			h, sh := is_hidden_folder(fp)
			if sh || h && !showall {
				continue
			}
			infos = append(infos, DirLink{IsLink: false, DirName: info.Name()})
		} else if info.Mode()&os.ModeSymlink != 0 {
			temp, err1 := os.Readlink(fp)
			newPath, err2 := filepath.EvalSymlinks(temp)
			if err1 == nil && err2 == nil {
				infos = append(infos, DirLink{IsLink: true, DirName: info.Name(), Target: newPath})
			}
		}
	}
	sort.Slice(infos, func(i, j int) bool {
		return path_comparer(infos[i].DirName, infos[j].DirName) < 0
	})

	fmt.Println(root)

	for i, info := range infos {
		printTree(b, root, info, "", i == len(infos)-1)
	}
}

func printTree(b *DataBag, root string, subdir DirLink, prefix string, is_last bool) {
	var entry_prefix string
	var new_prefix string
	if is_last {
		entry_prefix = "└── "
		new_prefix = prefix + "    "
	} else {
		entry_prefix = "├── "
		new_prefix = prefix + "│   "
	}

	fmt.Print(prefix + entry_prefix + subdir.DirName)
	if subdir.IsLink {
		fmt.Println(" -> " + subdir.Target)
		b.LinksCount++
		return
	}

	b.DirsCount++
	fmt.Println("")

	subdir_fp := filepath.Join(root, subdir.DirName)
	entries, err := os.ReadDir(subdir_fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a valid directory: %v:\n", APP_NAME, subdir_fp, err)
		return
	}

	infos := make([]DirLink, 0)
	for _, entry := range entries {
		fp := filepath.Join(root, entry.Name())
		info, err := entry.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Errror processing '%s' entry: %s\n", APP_NAME, root, err)
		} else if info.IsDir() {
			h, sh := is_hidden_folder(fp)
			// Ignore well-hidden directories such as $RECYCLE.BIN
			if sh || h && !showall {
				continue
			}
			infos = append(infos, DirLink{IsLink: false, DirName: info.Name()})
		} else if info.Mode()&os.ModeSymlink != 0 {
			temp, err1 := os.Readlink(fp)
			newPath, err2 := filepath.EvalSymlinks(temp)
			if err1 == nil && err2 == nil {
				infos = append(infos, DirLink{IsLink: true, DirName: info.Name(), Target: newPath})
			}
		}
	}
	sort.Slice(infos, func(i, j int) bool {
		return path_comparer(infos[i].DirName, infos[j].DirName) < 0
	})

	for i, info := range infos {
		printTree(b, subdir_fp, info, new_prefix, i == len(infos)-1)
	}
}

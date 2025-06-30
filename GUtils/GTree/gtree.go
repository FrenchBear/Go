// gtree.go
// Tree command in Go
// without -a option, doesn't show hidden directories or directories starting with a dot
// SYSTEM+HIDDEN directories are always skipped
// For now, the code in Windows-only until I learn how to compile code conditionally
//
// 2025-06-29	PV		First version
// 2025-06-30	PV		Sort files and folders using StrCmpLogicalW used by Windows file explorer

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"
	"unsafe"
)

// Options
var h1, h2 bool
var verbose bool
var showall bool

// Global constants
const APP_NAME string = "gtree"
const APP_VERSION string = "1.1.0"

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

// Define the signature of StrCmpLogicalW
// int StrCmpLogicalW(LPCWSTR psz1, LPCWSTR psz2);
// Returns:
//   -1 if psz1 comes before psz2
//    0 if psz1 is identical to psz2
//    1 if psz1 comes after psz2
var (
	shlwapi        = syscall.NewLazyDLL("shlwapi.dll")
	strCmpLogicalW = shlwapi.NewProc("StrCmpLogicalW")
)

// StrCmpLogicalWGo is a Go wrapper for the Windows API StrCmpLogicalW function.
// It compares two strings using the natural sort algorithm, similar to Windows File Explorer.
func StrCmpLogicalWGo(s1, s2 string) int {
	// Convert Go strings to null-terminated UTF-16 pointers for Windows API
	// syscall.UTF16PtrFromString allocates memory that needs to be freed
	// It's safer to use the x/sys/windows package for this.
	// For simplicity in this example, we'll use syscall directly,
	// but be aware of potential memory considerations in very high-performance loops.
	// In most cases, the garbage collector will handle it.

	// Create UTF-16 pointers
	p1, err := syscall.UTF16PtrFromString(s1)
	if err != nil {
		// Handle error: perhaps log or panic, depending on your application's needs.
		// For a comparison function, panicking might be acceptable if input is guaranteed valid.
		panic(fmt.Sprintf("Failed to convert string 1 to UTF16: %v", err))
	}
	p2, err := syscall.UTF16PtrFromString(s2)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert string 2 to UTF16: %v", err))
	}

	// Call the underlying Windows API function
	// The uintptr(unsafe.Pointer(p1)) converts the pointer to a uintptr,
	// which is what NewProc.Call expects for arguments.
	ret, _, _ := strCmpLogicalW.Call(uintptr(unsafe.Pointer(p1)), uintptr(unsafe.Pointer(p2)))

	// StrCmpLogicalW returns:
	// < 0 if psz1 comes before psz2
	//   0 if psz1 is identical to psz2
	// > 0 if psz1 comes after psz2
	// We cast the result to int64 first to handle potential negative return values correctly.
	return int(int32(ret))
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
			h, s, e := isHiddenOrSystemWindows(filepath.Join(root, info.Name()))
			if e != nil {
				fmt.Fprintf(os.Stderr, "%s: Errror processing '%s' entry: %s\n", APP_NAME, root, e)
				continue
			}
			// Ignore well-hidden directories such as $RECYCLE.BIN
			if h && s {
				continue
			}
			if !showall && (h || info.Name()[0] == '.') {
				continue
			}

			infos = append(infos, info)
		}
	}
	sort.Slice(infos, func(i, j int) bool {
		return StrCmpLogicalWGo(infos[i].Name(), infos[j].Name()) < 0
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
			h, s, e := isHiddenOrSystemWindows(filepath.Join(subdir_fp, info.Name()))
			if e != nil {
				fmt.Fprintf(os.Stderr, "%s: Errror processing '%s' entry: %s\n", APP_NAME, root, e)
				continue
			}
			// Ignore well-hidden directories such as $RECYCLE.BIN
			if h && s {
				continue
			}
			if !showall && (h || info.Name()[0] == '.') {
				continue
			}
			infos = append(infos, info)
		}
	}
	sort.Slice(infos, func(i, j int) bool {
		return StrCmpLogicalWGo(infos[i].Name(), infos[j].Name()) < 0
	})

	for i, info := range infos {
		printTree(b, subdir_fp, info.Name(), new_prefix, i == len(infos)-1)
	}
}

// Check if a file or directory has the HIDDEN or SYSTEM attribute on Windows.
// This function is Windows-specific.
func isHiddenOrSystemWindows(path string) (bool, bool, error) {
	// Convert the path to a UTF-16 pointer for Windows API.
	// Use windows.UTF16PtrFromString for safer conversion if available (golang.org/x/sys/windows)
	// Otherwise, syscall.UTF16PtrFromString.
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, false, fmt.Errorf("failed to convert path to UTF16: %w", err)
	}

	// Get file attributes using GetFileAttributesW.
	attributes, err := syscall.GetFileAttributes(pathPtr)
	if err != nil {
		return false, false, fmt.Errorf("failed to get file attributes for %s: %w", path, err)
	}

	// Check the HIDDEN and SYSTEM bits.
	isHidden := (attributes & syscall.FILE_ATTRIBUTE_HIDDEN) != 0
	isSystem := (attributes & syscall.FILE_ATTRIBUTE_SYSTEM) != 0

	return isHidden, isSystem, nil
}

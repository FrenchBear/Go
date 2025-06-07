// g08_gwhich.go
// Learning go, which utility
//
// 2025-06-04	PV		First version

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide an argument!")
		return
	}
	file := arguments[1]
	path := os.Getenv("PATH")
	pathSplit := filepath.SplitList(path)
	seenPaths := make(map[string]struct{})

	for _, directory := range pathSplit {
		// If path ends with a ;, range returns an empty directory, just ignore it
		if directory == "" {
			continue
		}

		// Check that PATH segment exists as a directory
		dirExists, err := DirExists(directory)
		if !dirExists {
			fmt.Println("Inexistent dir in PATH:", directory)
			continue
		}

		// Normalize to detect duplicates
		normalized, err := NormalizeDirPath(directory)
		if err != nil {
			// Decide how to handle errors during normalization.
			// For this example, we'll report them but continue processing.
			fmt.Printf("Warning: Could not normalize path '%s': %v\n", directory, err)
			continue
		}
		if _, found := seenPaths[normalized]; found {
			// This normalized path has been seen before, so the original 'dir' is a duplicate.
			fmt.Println("Duplicate dir in PATH:", directory)
		} else {
			seenPaths[normalized] = struct{}{}
		}

		// Does it exist?
		fullPath := filepath.Join(directory, file)
		fileInfo, err := os.Stat(fullPath)
		if err == nil {
			mode := fileInfo.Mode()
			//fmt.Printf("********************* Found: %v  %v\n", mode, mode.IsRegular())
			// Is it a regular file?
			if mode.IsRegular() {
				//fmt.Printf("********************* Regular: %v\n", mode)
				if runtime.GOOS == "windows" {
					fmt.Println("Found:", fullPath)
				} else {
					// Is it executable?
					if mode&0111 != 0 {
						fmt.Println("Found:", fullPath)
					}
				}
			}
		}
	}
}

// DirExists checks if a directory exists and is actually a directory.
func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		// Path exists, now check if it's a directory
		return info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		// Path does not exist
		return false, nil
	}
	// Some other error occurred (e.g., permissions)
	return false, err
}

// NormalizeDirPath prepares a directory path for case-insensitive and trailing-slash-agnostic comparison.
func NormalizeDirPath(path string) (string, error) {
	// 1. Clean the path: resolves ".." and "." and removes redundant slashes.
	// It also typically removes trailing slashes unless it's a root.
	cleanedPath := filepath.Clean(path)

	// 2. Convert to absolute path for consistent comparison, handling relative paths.
	// If you absolutely need to preserve relative paths, you'd skip this and
	// ensure your input relative paths are consistent, potentially by resolving them
	// against a known base directory.
	absPath, err := filepath.Abs(cleanedPath)
	if err != nil {
		// Custom error message with fmt.Errorf()
		return "", fmt.Errorf("failed to get absolute path for '%s': %w", path, err)
	}

	// 3. Convert to lowercase for case-insensitive comparison.
	normalizedPath := strings.ToLower(absPath)

	// filepath.Clean already handles most trailing slash cases.
	// For example, on Unix, filepath.Clean("/foo/bar/") results in "/foo/bar".
	// On Windows, filepath.Clean("C:\\foo\\bar\\") results in "C:\\foo\\bar".
	// The only exception is the root directory, which will remain with a trailing slash (e.g., "/" or "C:\").
	// We usually want to treat "C:\", "C:", and "C" as the same root.
	// For this, we'll ensure that only the actual root (e.g., "C:\") retains its trailing slash,
	// and any other path that *could* be a root (like "C:") is normalized to "C:\".

	// Special handling for Windows drive roots if `filepath.Abs` returns without trailing slash
	// for a root drive (e.g., "C:"). `filepath.Clean` usually adds the backslash for drive roots.
	// Let's ensure it's consistent.
	if len(normalizedPath) == 2 && normalizedPath[1] == ':' { // e.g., "c:"
		normalizedPath += string(filepath.Separator) // Add trailing slash for drive root
	}

	return normalizedPath, nil
}

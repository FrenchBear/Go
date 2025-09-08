// GFind, go version of Search/find/XFind/RFind utility
//
// 2025-07-12 	PV 		First version
// 2025-07-13 	PV 		1.1.0 Option -nop
// 2025-08-13 	PV 		1.2.0 Support for Windows Recycle Bin
// 2025-09-07 	PV 		1.3.0 Option -maxdepth
// 2025-09-08 	PV 		1.3.1 Use MyGlob 1.5 with a queue instead of a stack for a more logical output order

// go mod edit -replace github.com/PieVio/MyMarkup=../../Packages/MyMarkup
// go mod tidy

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PieVio/MyGlob"
)

const (
	APP_NAME        = "gfind"
	APP_VERSION     = "1.3.1"
	APP_DESCRIPTION = "Searching files in Go"
)

func main() {
	options, err := NewOptions()
	if err != nil {
		msg := err.Error()
		if msg == "" {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "%s: Problem parsing arguments: %s\n", APP_NAME, msg)
		os.Exit(1)
	}

	start := time.Now()

	// Adjust sources if option -name is used
	if len(options.names) > 0 {
		var name string
		if len(options.names) == 1 {
			name = options.names[0]
		} else {
			name = strings.Join(options.names, ",")
		}

		for i, source := range options.sources {
			info, err := os.Stat(source)
			if err == nil && info.IsDir() {
				options.sources[i] = filepath.Join(options.sources[i], "**", name)
			}
		}
	}

	// Convert String sources into MyGlobSearch structs
	sources := make([]*MyGlob.MyGlobSearch, len(options.sources))
	for i, source := range options.sources {
		mg, err := MyGlob.New(source).Autorecurse(options.autorecurse).MaxDepth(options.maxdepth).Compile()
		if err != nil {
			fmt.Printf("*** Error building MyGlob: %v\n", err)
			continue
		}
		sources[i] = mg
	}

	if len(sources) == 0 {
		fmt.Fprintf(os.Stderr, "*** No source specified. Use %s ? to show usage.", APP_NAME)
		os.Exit(1)
	}

	if options.verbose {
		fmt.Print("Sources(s): ")
		if options.search_dirs && options.search_files {
			fmt.Println("(search for files and directories)")
		} else if options.search_dirs {
			fmt.Println("(search for directories)")
		} else {
			fmt.Println("(search for files)")
		}

		for _, source := range options.sources {
			fmt.Println("- ", source)
		}
	}

	actions := make([]IAction, 0) //, len(options.actions_names))
	for action := range options.actions_names {
		switch action {
		case "print":
			if _, ok := options.actions_names["dir"]; ok {
				fmt.Println("*** Both actions print and dir used, action print ignored.")
			} else {
				actions = append(actions, &action_print{detailed_output: false})
			}

		case "nop":
			// Do nothing

		case "dir":
			actions = append(actions, &action_print{detailed_output: true})

		case "delete":
			actions = append(actions, &action_delete{recycle: options.recycle})

		case "rmdir":
			actions = append(actions, &action_rmdir{recycle: options.recycle})

		default:
			panic("Invalid action: " + action)
		}
	}

	if options.verbose {
		fmt.Print("\nAction(s): ")
		if options.noaction {
			fmt.Println("(no action will actually be performed)")
		} else {
			fmt.Println()
		}

		for _, action := range actions {
			fmt.Println("-", action.name())
		}
		fmt.Println()
		if options.isempty {
			fmt.Println("Only search for empty files or directories")
		}
	}

	files_count := 0
	dirs_count := 0

	for _, gs := range sources {
		for ma := range gs.Explore() {
			if ma.Err != nil {
				if options.verbose {
					fmt.Printf("%s: MyGlobMatch error %v\n", APP_NAME, ma.Err)
				}
				continue
			}
			info, err := os.Stat(ma.Path)
			if err != nil {
				continue
			}

			if !ma.IsDir {
				if options.search_files && (!options.isempty || info.Size() == 0) {
					files_count++
					for _, ba := range actions {
						ba.action(ma.Path, info, options.noaction, options.verbose)
					}
				}
			} else {
				if options.search_dirs && (!options.isempty || IsDirEmpty(ma.Path)) {
					dirs_count++
					for _, ba := range actions {
						ba.action(ma.Path, info, options.noaction, options.verbose)
					}
				}
			}
		}
	}

	duration := time.Since(start)

	if options.verbose {
		if files_count+dirs_count > 0 {
			fmt.Println()
		}

		if options.search_files {
			fmt.Printf("%d file(s)", files_count)
		}
		if options.search_dirs {
			if options.search_files {
				fmt.Print(", ")
			}
			fmt.Printf("%d dir(s)", dirs_count)
		}
		fmt.Printf(" found in %.3fs\n", float64(duration.Milliseconds())/1000.0)
	}
}

// IsDirEmpty checks if a directory is empty. It returns true if the directory
// is empty, and false otherwise. An error is returned if the path does not
// exist or is not a directory.
func IsDirEmpty(path string) bool {
	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		return false
	}
	defer dir.Close()

	// Attempt to read just one directory entry.
	// Readdirnames(1) is more efficient than Readdir(-1) as it stops after the first entry.
	_, err = dir.Readdirnames(1)

	// If we get an End-Of-File error, it means the directory is empty.
	if err == io.EOF {
		return true
	}

	// If there was another error, return it.
	if err != nil {
		return false
	}

	// Otherwise, the directory is not empty.
	return false
}

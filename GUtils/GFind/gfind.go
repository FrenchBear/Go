// go mod edit -replace github.com/PieVio/MyMarkup=../../Packages/MyMarkup
// go mod tidy

package main

import (
	"fmt"
	"os"
	// "github.com/PieVio/MyMarkup"
)

const (
	appName    = "gfind"
	appVersion = "1.0.0"
)

const (
	APP_NAME        = "gfind"
	APP_VERSION     = "1.0.0"
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

	fmt.Println(options)

}

	/*
	start := time.Now()

	// Adjust sources if option -name is used
	if len(args.Names) > 0 {
		name := strings.Join(args.Names, ",")
		if len(args.Names) > 1 {
			name = "{" + name + "}"
		}

		for i, source := range args.Sources {
			info, err := os.Stat(source)
			if err == nil && info.IsDir() {
				args.Sources[i] = filepath.Join(source, "**", name)
			}
		}
	}

	// devtemp
	fmt.Printf("%q\n", args)
	for _, source := range args.Sources {
		fmt.Println("- ", source)
	}
	duration := time.Since(start)


	fmt.Println(duration.Seconds())

	/*
	var filesCount, dirsCount int

	for _, source := range args.Sources {
		mg, err := MyGlob.New(source).Autorecurse(args.AutoRecurse).Compile()
		if err != nil {
			fmt.Printf("*** Error building MyGlob: %v\n", err)
			continue
		}

		for match := range mg.Explore() {
			if match.Err != nil {
				if args.Verbose {
					fmt.Printf("%s: MyGlobMatch error %v\n", appName, match.Err)
				}
				continue
			}

			info, err := os.Stat(match.Path)
			if err != nil {
				continue
			}

			isFile := !info.IsDir()
			isDir := info.IsDir()

			searchFile := args.Search == "files" || args.Search == "both"
			searchDir := args.Search == "dirs" || args.Search == "both"

			if (searchFile && isFile) || (searchDir && isDir) {
				if args.IsEmpty {
					if isFile && info.Size() != 0 {
						continue
					}
					if isDir {
						f, err := os.Open(match.Path)
						if err != nil {
							continue
						}
						_, err = f.Readdirnames(1)
						f.Close()
						if err == nil { // Not empty
							continue
						}
					}
				}

				if isFile {
					filesCount++
				} else {
					dirsCount++
				}

				action(match.Path, info)
			}
		}
	}

	duration := time.Since(start)

	if args.Verbose {
		if filesCount+dirsCount > 0 {
			fmt.Println()
		}
		if args.Search == "files" || args.Search == "both" {
			fmt.Printf("%d files(s)", filesCount)
		}
		if args.Search == "dirs" || args.Search == "both" {
			if args.Search == "both" {
				fmt.Print(", ")
			}
			fmt.Printf("%d dir(s)", dirsCount)
		}
		fmt.Printf(" found in %.3fs\n", duration.Seconds())
	}
}

func action(path string, info os.FileInfo) {
	if args.Print || (!args.Dir && !args.Delete && !args.RmDir) {
		fmt.Println(path)
	}

	if args.Dir {
		modTime := info.ModTime().Format("2006-01-02 15:04:05")
		size := ""
		if !info.IsDir() {
			size = humanize.Bytes(uint64(info.Size()))
		} else {
			size = "<DIR>"
		}
		fmt.Printf("%-20s %15s %s\n", modTime, size, path)
	}

	if args.Delete && !info.IsDir() {
		performDelete(path)
	}

	if args.RmDir && info.IsDir() {
		performDelete(path)
	}
}

func performDelete(path string) {
	quotedPath := path
	if strings.Contains(quotedPath, " ") {
		quotedPath = "\"" + quotedPath + "\""
	}

	logMsg := ""
	if args.Recycle {
		logMsg = fmt.Sprintf("RECYCLE %s", quotedPath)
	} else {
		logMsg = fmt.Sprintf("DEL %s", quotedPath)
	}
	fmt.Println(MyMarkup.BuildMarkup(logMsg))

	if !args.NoAction {
		var err error
		if args.Recycle {
			// Go does not have a built-in recycle bin functionality.
			// For now, we will just print a message.
			fmt.Println(MyMarkup.BuildMarkup("  `-> Recycle bin not implemented in this Go version."))
		} else {
			err = os.RemoveAll(path)
		}

		if err != nil {
			fmt.Printf("*** Error deleting %s: %v\n", quotedPath, err)
		} else if args.Verbose {
			fmt.Printf("File %s deleted successfully.\n", quotedPath)
		}
	}
		*/

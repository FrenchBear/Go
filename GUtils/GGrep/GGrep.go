// GGrep tool, Go version og grep
//
// 2025-08-130 	PV 		First version

package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/PieVio/MyGlob"
	"github.com/PieVio/TextAutoDecode"
)

const (
	APP_NAME        = "ggrep"
	APP_VERSION     = "1.0.0"
	APP_DESCRIPTION = "Grep utility in Go"
)

type DataBag struct {
	files_count int
}

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

	re, err := BuildRegexp(options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Problem building regexp: %v\n", APP_NAME, err)
		os.Exit(1)
	}

	start := time.Now()

    // Building list of files
    // It could be better to process file just when it's returned by iterator rather than stored in a Vec and processed
    // later... but then we don't know when processing the first file whether there's more than one, to print paths...
	files := []string{}
	for _, source := range options.Sources {
		gs, err := MyGlob.New(source).Autorecurse(options.Autorecurse).Compile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Error building MyGlob: %v\n", APP_NAME, err)
			continue
		}

		for ma := range gs.Explore() {
			if ma.Err != nil {
				if options.Verbose {
					fmt.Fprintf(os.Stderr, "%s: error %v\n", APP_NAME, ma.Err)
				}
				continue
			}
			if !ma.IsDir { // We ignore matching directories in rgrep, we only look for files
				//processPath(&b, re, ma.Path, options)
				files = append(files, ma.Path)
			}
		}
	}

	// If no source has been provided, use stdin
	b := DataBag{}
	if len(options.Sources) == 0 {
		err := processStdin(&b, re, options)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Error reading from stdin: %v\n", APP_NAME, err)
		}
	} else {
		if len(files) > 1 {
			options.ShowPath = true
		}
		for _, path := range files {
			processPath(&b, re, path, options)
		}
	}

	duration := time.Since(start)

	if options.Verbose {
		if len(options.Sources) == 0 {
			fmt.Print("\nstdin")
		} else {
			fmt.Printf("\n%d file", b.files_count)
			if b.files_count > 1 {
				fmt.Printf("s")
			}
		}
		fmt.Printf(" searched in %.3fs\n", duration.Seconds())
	}
}

// Helper, build Regex according to options (case, fixed string, whole word).
// Return an error in case of invalid Regex.
func BuildRegexp(options *Options) (*regexp.Regexp, error) {
	var spat string

	if options.FixedString {
		spat = regexp.QuoteMeta(options.Pattern)
	} else {
		spat = options.Pattern
	}

	if options.WholeWord {
		spat = fmt.Sprintf("\\b%s\\b", spat)
	}

	if options.IgnoreCase {
		spat = "(?im)" + spat
	} else {
		spat = "(?m)" + spat
	}

	return regexp.Compile(spat)
}

func processStdin(b *DataBag, re *regexp.Regexp, options *Options) error {
	if options.Verbose {
		fmt.Println("Reading from stdin")
	}

	byteData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	inputString := string(byteData)

	processText(b, re, inputString, "(stdin)", options)
	return nil
}

func processPath(b *DataBag, re *regexp.Regexp, path string, options *Options) {
	tadRes, err := TextAutoDecode.ReadTextFile(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "*** Error reading file %s: %v\n", path, err)
		return
	}

	if tadRes.Encoding == TextAutoDecode.TFE_NotText {
		if options.Verbose {
			fmt.Printf("%s: ignored non-text file %s\n", APP_NAME, path)
		}
	} else {
		// fileInfo, err := os.Stat(path)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "%s: Error getting info for file %s: %v\n", APP_NAME, path, err)
		// 	return
		// }
		// if fileInfo.Size() > 1024*1024*1024 {
		// 	if options.Verbose {
		// 		fmt.Printf("%s: ignored very large file %s, size: %d bytes\n", APP_NAME, path, fileInfo.Size())
		// 	}
		// } else {
		processText(b, re, tadRes.Text, path, options)
	}

}

func processText(b *DataBag, re *regexp.Regexp, txt, path string, options *Options) {
	matchlinecount := 0
	for gi := range Grep(txt, re) {
		matchlinecount++

		// ToDo: add colored output
		if options.OutLevel == 1 {
			fmt.Printf("%s\n", path)
			return
		}

		if options.OutLevel == 0 {
			if options.ShowPath {
				fmt.Printf("%s: ", path)
			}
			p := 0
			for _, ma := range gi.Ranges {
				if ma.Start < len(gi.Line) {
					e := ma.End
					fmt.Printf("%s«%s»", gi.Line[p:ma.Start], gi.Line[ma.Start:ma.End])
					p = e
				}
			}
			fmt.Println(gi.Line[p:])
		}
	}

	// Note: Using together options -c and -l (out_level==3) is not supported by Linux grep command
	if options.OutLevel == 2 || (options.OutLevel == 3 && matchlinecount > 0) {
		fmt.Printf("%s: %d\n", path, matchlinecount)
	}

	b.files_count++
}

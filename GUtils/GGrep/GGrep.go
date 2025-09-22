// GGrep tool, Go version of grep
//
// 2025-08-13 	PV 		First version
// 2025-08-18	PV 		1.1 Process files while enumerating; use MyGlob.SetChannelSize(25) to speed up globbing
// 2025-09-22   PV      Option -v -> -t to show execution time. Option -v to invert the sense of matching, to select non-matching lines

package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/PieVio/MyGlob"
	"github.com/PieVio/TextAutoDecode"
	"github.com/mattn/go-isatty"
)

const (
	APP_NAME        = "ggrep"
	APP_VERSION     = "1.2.0"
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

	// Need to wait for 2nd file to call processPath, since if there is a 2nd file, we set options.ShowPath to true
	// to show filename before matches.  file_to_process is the file from the previous loop
	file_to_process := ""
	b := DataBag{}
	for _, source := range options.Sources {
		gs, err := MyGlob.New(source).Autorecurse(options.Autorecurse).ChannelSize(25).Compile()
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
				// We've met out second file!
				if file_to_process != "" {
					options.ShowPath = true
					processPath(&b, re, file_to_process, options)
				}
				file_to_process = ma.Path
			}
		}
	}
	if file_to_process != "" {
		processPath(&b, re, file_to_process, options)
	}

	// If no source has been provided, use stdin
	if len(options.Sources) == 0 {
		err := processStdin(&b, re, options)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Error reading from stdin: %v\n", APP_NAME, err)
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

	if isatty.IsTerminal(os.Stdout.Fd()) {
		// tty output in color
		const BrightBlack string = "\033[90m"
		const BoldRed string = "\033[1;31m"
		const NormalColor string = "\033[0;37m"

		var iter <-chan GrepLineMatches
		if options.InvertMatch {
			iter = GrepInvert(txt, re)
		} else {
			iter = Grep(txt, re)
		}

		for gi := range iter {
			matchlinecount++

			if options.OutLevel == 1 {
				fmt.Printf("%s\n", path)
				return
			}

			if options.OutLevel == 0 && options.InvertMatch {
				if options.ShowPath {
					fmt.Printf("%s%s:%s ", BrightBlack, path, NormalColor)
				}
				fmt.Println(gi.Line)
			} else if options.OutLevel == 0 {
				if options.ShowPath {
					fmt.Printf("%s%s:%s ", BrightBlack, path, NormalColor)
				}
				p := 0
				for _, ma := range gi.Ranges {
					if ma.Start < len(gi.Line) {
						e := ma.End
						fmt.Printf("%s%s%s%s", gi.Line[p:ma.Start], BoldRed, gi.Line[ma.Start:ma.End], NormalColor)
						p = e
					}
				}
				fmt.Println(gi.Line[p:])
			}
		}
	} else {
		var iter <-chan GrepLineMatches
		if options.InvertMatch {
			iter = GrepInvert(txt, re)
		} else {
			iter = Grep(txt, re)
		}

		// Not a tty, monochrome output
		for gi := range iter {
			matchlinecount++

			if options.OutLevel == 1 {
				fmt.Printf("%s\n", path)
				return
			}

			if options.OutLevel == 0 {
				if options.ShowPath {
					fmt.Printf("%s: ", path)
				}
				fmt.Println(gi.Line)
			}
		}
	}

	// Note: Using together options -c and -l (out_level==3) is not supported by Linux grep command
	if options.OutLevel == 2 || (options.OutLevel == 3 && matchlinecount > 0) {
		fmt.Printf("%s: %d\n", path, matchlinecount)
	}

	b.files_count++
}

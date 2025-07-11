// Gwc tool, prints text type for text files
// Translation of Rwc utility
//
// 2027-07-10 	PV 		First version
// 2027-07-11 	PV 		1.1 Parallel version of ProcessText

/* Before parallelism, on WOTAN:

gwc -v "C:\Development\TestFiles\Text\Les secrets d'Hermione.txt"
  56337 1363732  7946200  8490462  C:\Development\TestFiles\Text\Les secrets d'Hermione.txt
1 files(s) searched in 0.071s

After parallelism:
  56337 1363732  7946200  8490462  C:\Development\TestFiles\Text\Les secrets d'Hermione.txt
1 files(s) searched in 0.046s

*/

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PieVio/MyGlob"
	"github.com/PieVio/TextAutoDecode"
)

const (
	APP_NAME        = "gwc"
	APP_VERSION     = "1.1.0"
	APP_DESCRIPTION = "Word Count utility in Go"
)

type DataBag struct {
	files_count int
	lines_count int
	words_count int
	chars_count int
	bytes_count int
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

	start := time.Now()

	bTotal := DataBag{}

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
				processFile(&bTotal, ma.Path, options)
			}
		}
	}

	// If no source has been provided, use stdin
	if len(options.Sources) == 0 {
		err := processStdin(options)
		if err!=nil {
			fmt.Fprintf(os.Stderr, "%s: Error reading from stdin: %v\n", APP_NAME, err)
		}
	}

	duration := time.Since(start)

	if bTotal.files_count > 1 || options.ShowOnlyTotal {
		name := "total"
		if bTotal.files_count > 1 {
			name += fmt.Sprintf(" (%d files)", bTotal.files_count)
		}
		printResultOneFile(&bTotal, name)
	}

	if options.Verbose {
		fmt.Printf("\n%d files(s) searched in %.3fs\n", bTotal.files_count, duration.Seconds())
	}
}

func printResultOneFile(b *DataBag, filename string) {
	printLine(b.lines_count, b.words_count, b.chars_count, b.bytes_count, filename)
}

func printLine(lines, words, chars, bytes int, filename string) {
	fmt.Printf("%7d %7d %8d %8d  %s\n", lines, words, chars, bytes, filename)
}

func processStdin(options *Options) error {
	if options.Verbose {
		fmt.Println("Reading from stdin")
	}

	byteData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	inputString := string(byteData)

	b := DataBag{}
	processText(&b, inputString, "(stdin)", options, int64(len(inputString)))
	return nil
}

func processFile(b *DataBag, path string, options *Options) {
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
		fileInfo, err := os.Stat(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Error getting info for file %s: %v\n", APP_NAME, path, err)
			return
		}
		if fileInfo.Size() > 1024*1024*1024 {
			if options.Verbose {
				fmt.Printf("%s: ignored very large file %s, size: %d bytes\n", APP_NAME, path, fileInfo.Size())
			}
		} else {
			processText(b, tadRes.Text, path, options, fileInfo.Size())
		}
	}
}

func processText(b *DataBag, txt, path string, options *Options, filesize int64) {
	normalized := strings.ReplaceAll(strings.ReplaceAll(txt, "\r\n", "\n"), "\r", "\n")
	textLines := strings.Split(normalized, "\n")

	lines := len(textLines)
	chars := utf8.RuneCountInString(txt)
	bytes := int(filesize) // sizes longer than 1GB are skipped

	// Special correction: If last line ends with \n, strings.Split counts an extra empty line in Go, while it's not
	// counted in Rust version and also in Linux wc command, hence this manual correction
	if strings.HasSuffix(normalized, "\n") {
		lines--
	}

	// To count words, we use a goroutine to count in blocks of 6000 lines since empirically that's near the most efficient size
	SLICESIZE := 6000
	blocks := len(textLines)/SLICESIZE + 1
	reschan := make(chan int, blocks)
	sl := 0
	for i := 0; i < len(textLines); i += SLICESIZE {
		end := i + SLICESIZE
		if end > len(textLines) {
			end = len(textLines)
		}
		sl++
		go count_slice_words_to_reschan(textLines[i:end], reschan)
	}

	words := 0
	for i := 0; i < sl; i++ {
		words += <-reschan
	}
	close(reschan)


	if !options.ShowOnlyTotal {
		printLine(lines, words, chars, bytes, path)
	}

	b.files_count++
	b.lines_count += lines
	b.words_count += words
	b.chars_count += chars
	b.bytes_count += bytes
}

func count_slice_words_to_reschan(lines []string, reschan chan <- int) {
	words := 0
	for _, line := range lines {
		splitFunc := func(r rune) bool {
			return r == ' ' || r == '\t'
		}
		words += len(strings.FieldsFunc(strings.Trim(line, " \t"), splitFunc))
	}
	reschan <- words
}

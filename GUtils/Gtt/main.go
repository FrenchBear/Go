// Gtt tool, prints text type for text files
// Translation of Rtt utility
//
// 2027-07-05 	PV 		Initial translation by Gemini

/*
I need to translate a simple command line Rust program into its equivalent in Go.

The original Rust source uses three external crates I developed myself, myGlob, myMarkup and TextAutoDecode, and I have
the equivalent in Go with exactlty the same name and identical public API:

- MyGlob package is defined here: @C:\Users\Pierr\Gemini\Packages\MyGlob\go.mod
  @C:\Users\Pierr\Gemini\Packages\MyGlob\myglob.go

- MyMarkup package is defined here: @C:\Users\Pierr\Gemini\Packages\MyMarkup\go.mod
  @C:\Users\Pierr\Gemini\Packages\MyMarkup\mymarkup.go

- TextAutoDecode is defined here: @C:\Users\Pierr\Gemini\Packages\TextAutoDecode\go.mod
  C:\Users\Pierr\Gemini\Packages\TextAutoDecode\textautodecode.go

Three other crate dependencies shoud be replaced as follows:
- getopt crate should be replaces by Go standard Flags package.
- colored crate should be replaced by Go external crate github.com/fatih/color.
- tempfile crate should be replaced by its Go equivalent.

Finally, source code to convert in Go is in in three files: @C:\Users\Pierr\Gemini\Rtt\src\main.rs
@C:\Users\Pierr\Gemini\Rtt\src\options.rs @C:\Users\Pierr\Gemini\Rtt\src\tests.rs

No need to actually verify that the tests are running correctly, just produce correct Go equivalent, I'll run the tests
and debugging myself.
*/

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/PieVio/MyGlob"
	"github.com/PieVio/TextAutoDecode"
	"github.com/fatih/color"
)

const (
	APP_NAME    = "rtt"
	APP_VERSION = "1.0.1"
)

// These extensions should indicate a text content
var TEXT_EXT = []string{
	// Sources
	"awk", "c", "cpp", "cs", "fs", "go", "h", "java", "jl", "js", "lua", "py", "rs", "sql", "ts", "vb", "xaml",
	// VB6
	"bas", "frm", "cls", "ctl", "vbp", "vbg",
	// Projects
	"sln", "csproj", "vbproj", "fsproj", "pyproj", "vcxproj",
	// Misc
	"appxmanifest", "clang-format", "classpath", "ruleset", "editorconfig", "gitignore", "globalconfig", "resx", "targets", "pubxml", "filters",
	// Config
	"ini", "xml", "yml", "yaml", "json", "toml",
	// Scripts
	"bat", "cmd", "ps1", "sh", "vbs",
	// Text
	"txt", "md",
}

type DataBag struct {
	FilesTypes FileTypeCounts
	EolStyles  EOLStyleCounts
	Counters   map[string]map[string]*DirectoryExtCounts
}

func NewDataBag() *DataBag {
	return &DataBag{
		FilesTypes: FileTypeCounts{},
		EolStyles:  EOLStyleCounts{},
		Counters:   make(map[string]map[string]*DirectoryExtCounts),
	}
}

type DirectoryExtCounts struct {
	FilesTypes FileTypeCounts
	EolStyles  EOLStyleCounts
}

type FileTypeCounts struct {
	Total    int
	Empty    int
	Ascii    int
	Utf8     int
	Utf16    int
	EightBit int
	NonText  int
}

type EOLStyleCounts struct {
	Total   int
	Windows int
	Unix    int
	Mac     int
	Mixed   int
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

	b := NewDataBag()

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
				printResult(processFile(b, ma.Path, ma.Path), options)
			}
		}
	}

	// Warnings per directory+extension
	if len(options.Sources) > 0 {
		headerPrinted := false
		var fk []string
		for k := range b.Counters {
			fk = append(fk, k)
		}
		sort.Strings(fk)

		for _, f := range fk {
			var ek []string
			for k := range b.Counters[f] {
				ek = append(ek, k)
			}
			sort.Strings(ek)

			for _, e := range ek {
				filePrinted := false

				ft := b.Counters[f][e].FilesTypes
				if (ft.Utf8 > 0 && ft.Utf16 > 0) || (ft.Utf8 > 0 && ft.EightBit > 0) || (ft.Utf16 > 0 && ft.Ascii > 0) || (ft.Utf16 > 0 && ft.EightBit > 0) {
					if !headerPrinted {
						fmt.Println("\nMixed directory contents:")
						headerPrinted = true
					}
					fmt.Printf("%s, ext .%s: ", f, e)
					filePrinted = true

					color.Red("Mixed text file contents")
				}

				eol := b.Counters[f][e].EolStyles
				if eol.Total > 1 {
					if (eol.Windows > 0 && eol.Unix > 0) || (eol.Windows > 0 && eol.Mac > 0) || (eol.Unix > 0 && eol.Mac > 0) {
						if !headerPrinted {
							fmt.Println("\nMixed directory contents:")
							headerPrinted = true
						}

						if filePrinted {
							fmt.Print(", ")
						} else {
							filePrinted = true
							fmt.Printf("%s, ext .%s: ", f, e)
						}
						color.Red("Mixed EOF styles")
					}
				}

				if filePrinted {
					fmt.Println()
				}
			}
		}
	}

	// If no source has been provided, use stdin
	if len(options.Sources) == 0 {
		processStdin(b, options)
	}

	duration := time.Since(start)

	if options.Verbose {
		fmt.Println("\nGlobal stats:")
		printFilesTypesCounts(b.FilesTypes)
		printEolStylesCounts(b.EolStyles)

		fmt.Printf("\n%d files(s) searched in %.3fs\n", b.FilesTypes.Total, duration.Seconds())
	}
}

func processStdin(b *DataBag, options *Options) error {
	if options.Verbose {
		fmt.Println("Reading from stdin")
	}

	tempFile, err := os.CreateTemp("", "rtt-stdin-")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = io.Copy(tempFile, os.Stdin)
	if err != nil {
		return err
	}
	tempFile.Sync()

	printResult(processFile(b, tempFile.Name(), "(stdin)"), options)
	return nil
}

func printFilesTypesCounts(f FileTypeCounts) {
	tot := f.Empty + f.Ascii + f.Utf8 + f.Utf16 + f.EightBit + f.NonText
	fmt.Printf("Total files: %d\n", tot)
	if f.Empty > 0 {
		fmt.Printf("- Empty: %d\n", f.Empty)
	}
	if f.Ascii > 0 {
		fmt.Printf("- ASCII: %d\n", f.Ascii)
	}
	if f.Utf8 > 0 {
		fmt.Printf("- UTF-8: %d\n", f.Utf8)
	}
	if f.Utf16 > 0 {
		fmt.Printf("- UTF-16: %d\n", f.Utf16)
	}
	if f.EightBit > 0 {
		fmt.Printf("- 8-Bit: %d\n", f.EightBit)
	}
	if f.NonText > 0 {
		fmt.Printf("- Non text: %d\n", f.NonText)
	}
}

func printEolStylesCounts(e EOLStyleCounts) {
	tot := e.Windows + e.Unix + e.Mac + e.Mixed
	fmt.Printf("Total EOL styles: %d\n", tot)
	if e.Windows > 0 {
		fmt.Printf("- Windows: %d\n", e.Windows)
	}
	if e.Unix > 0 {
		fmt.Printf("- Unix: %d\n", e.Unix)
	}
	if e.Mac > 0 {
		fmt.Printf("- Mac: %d\n", e.Mac)
	}
	if e.Mixed > 0 {
		fmt.Printf("- Mixed: %d\n", e.Mixed)
	}
}

func printResult(msg string, options *Options) {
	if !options.ShowOnlyWarnings || strings.Contains(msg, "«") {
		printResultCore(msg)
	}
}

func printResultCore(msg string) {
	p0 := 0
	for {
		p1 := strings.Index(msg[p0:], "«")
		if p1 == -1 {
			fmt.Println(msg[p0:])
			return
		}
		p1 += p0

		if p1 > p0 {
			fmt.Print(msg[p0:p1])
		}

		p2 := strings.Index(msg[p1+1:], "»")
		if p2 == -1 {
			// This should not happen based on the Rust code's expect
			panic(fmt.Sprintf("Internal error, unbalanced « » in %s", msg))
		}
		p2 += p1 + 1

		color.Red(msg[p1+1 : p2])
		p0 = p2 + 1
	}
}

func processFile(b *DataBag, pathForRead string, pathForName string) string {
	res := ""
	tadRes, err := TextAutoDecode.ReadTextFile(pathForRead)

	fmt.Println(pathForName)
	fmt.Println(tadRes)
	fmt.Println(err)

	b.FilesTypes.Total++
	if err != nil {
		fmt.Fprintf(os.Stderr, "*** Error reading file %s: %v\n", pathForName, err)
		return res
	}

	ext := strings.ToLower(filepath.Ext(pathForName))
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	dir := strings.ToLower(filepath.Dir(pathForName))

	if _, ok := b.Counters[dir]; !ok {
		b.Counters[dir] = make(map[string]*DirectoryExtCounts)
	}
	if _, ok := b.Counters[dir][ext]; !ok {
		b.Counters[dir][ext] = &DirectoryExtCounts{}
	}
	fc := b.Counters[dir][ext]

	var enc string
	var war string

	switch tadRes.Encoding {
	case TextAutoDecode.TFE_NotText:
		b.FilesTypes.NonText++
		fc.FilesTypes.NonText++
		// Silently ignore non-text files, but check whether it should have contained text
		isTextExt := false
		for _, te := range TEXT_EXT {
			if te == ext {
				isTextExt = true
				break
			}
		}
		if isTextExt {
			return fmt.Sprintf("%s: «Non-text file detected, but extension %s is usually a text file»", pathForName, ext)
		}
		return res
	case TextAutoDecode.TFE_Empty:
		b.FilesTypes.Empty++
		// Don't collect infos per directory+ext for empty files
		// No need to continue if it's empty
		return fmt.Sprintf("%s: «Empty file»", pathForName)
	case TextAutoDecode.TFE_ASCII:
		b.FilesTypes.Ascii++
		fc.FilesTypes.Ascii++
		enc = "ASCII"
		war = ""
	case TextAutoDecode.TFE_EightBit:
		b.FilesTypes.EightBit++
		fc.FilesTypes.EightBit++
		enc = "8-Bit text"
		war = ""
	case TextAutoDecode.TFE_UTF8, TextAutoDecode.TFE_UTF8BOM:
		b.FilesTypes.Utf8++
		fc.FilesTypes.Utf8++
		enc = "UTF-8"
		if tadRes.Encoding == TextAutoDecode.TFE_UTF8BOM {
			war = "with BOM"
		} else {
			war = ""
		}
	case TextAutoDecode.TFE_UTF16LE, TextAutoDecode.TFE_UTF16BE, TextAutoDecode.TFE_UTF16LEBOM, TextAutoDecode.TFE_UTF16BEBOM:
		b.FilesTypes.Utf16++
		fc.FilesTypes.Utf16++
		if tadRes.Encoding == TextAutoDecode.TFE_UTF16LE || tadRes.Encoding == TextAutoDecode.TFE_UTF16LEBOM {
			enc = "UTF-16 LE"
		} else {
			enc = "UTF-16 BE"
		}
		if tadRes.Encoding == TextAutoDecode.TFE_UTF16LE || tadRes.Encoding == TextAutoDecode.TFE_UTF16BE {
			war = "without BOM"
		} else {
			war = ""
		}
	}

	eol := getEol(tadRes.Text)

	fc.EolStyles.Windows += eol.Windows
	fc.EolStyles.Unix += eol.Unix
	fc.EolStyles.Mac += eol.Mac
	fc.EolStyles.Mixed += eol.Mixed
	fc.EolStyles.Total += eol.Total

	b.EolStyles.Windows += eol.Windows
	b.EolStyles.Unix += eol.Unix
	b.EolStyles.Mac += eol.Mac
	b.EolStyles.Mixed += eol.Mixed
	b.EolStyles.Total += eol.Total

	res = fmt.Sprintf("%s: %s", pathForName, enc)
	if war != "" {
		res += fmt.Sprintf(" «%s»", war)
	}
	res += ", "

	if eol.Mixed > 0 {
		res += "«Mixed EOL styles»"
	} else if eol.Windows+eol.Unix+eol.Mac == 0 {
		res += "No EOL detected"
	} else if eol.Windows > 0 {
		res += "Windows"
	} else if eol.Unix > 0 {
		res += "Unix"
	} else if eol.Mac > 0 {
		res += "Mac"
	}

	return res
}

func getEol(txt string) EOLStyleCounts {
	eol := EOLStyleCounts{}
	bytes := []byte(txt)
	for i := 0; i < len(bytes); i++ {
		c := bytes[i]
		switch c {
		case '\n':
			eol.Unix = 1
		case '\r':
			if i+1 < len(bytes) && bytes[i+1] == '\n' {
				i++
				eol.Windows = 1
			} else {
				eol.Mac = 1
			}
		}
	}

	// Don't count files without EOL detected in total
	if eol.Windows+eol.Unix+eol.Mac > 0 {
		eol.Total++
	}

	// Helper
	if eol.Windows+eol.Unix+eol.Mac > 1 {
		eol.Mixed = 1
	}

	return eol
}

// gtt options.go
// Parse and validate command line options, returning a clean Options struct
//
// 2025-07-05	PV 		First version, translated from Rust by Gemini
// 2025-07-07 	PV 		Compact options -a+ and -a-

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/PieVio/MyGlob"
	"github.com/PieVio/MyMarkup"
	"github.com/PieVio/TextAutoDecode"
)

type Options struct {
	Sources          []string
	Autorecurse      bool
	ShowOnlyWarnings bool
	Verbose          bool
}

func header() {
	fmt.Printf("%s %s\n", APP_NAME, APP_VERSION)
	fmt.Println(APP_DESCRIPTION)
}

func usage() {
	header()
	fmt.Println()
	text := `⌊Usage⌋: {APP_NAME} ¬[⦃?⦄|⦃-?⦄|⦃-h⦄|⦃??⦄|⦃-??⦄] [⦃-a+⦄|⦃-a-⦄] [⦃-w⦄] [⦃-v⦄] [⟨source⟩...]

⌊Options⌋:
⦃?⦄|⦃-?⦄|⦃-h⦄  ¬Show this message
⦃??⦄|⦃-??⦄   ¬Show advanced usage notes
⦃-a+⦄|⦃-a-⦄  ¬Enable (default) or disable glob autorecurse mode (see extended usage)
⦃-w⦄       ¬Only show warnings
⦃-v⦄       ¬Verbose output
⟨source⟩   ¬File or directory to search, glob syntax supported (see extended usage). Without source, search stdin.
`

	MyMarkup.RenderMarkup(strings.ReplaceAll(text, "{APP_NAME}", APP_NAME))
}

func extendedUsage() {
	header()
	fmt.Println("Copyright ©2025 Pierre Violent")
	fmt.Println()
	MyMarkup.RenderMarkup("⌊Dependencies⌋:")
	fmt.Println("MyGlob:", MyGlob.Version())
	fmt.Println("TextAutoDecode:", TextAutoDecode.Version())
	fmt.Println("MyMarkup:", MyMarkup.Version())
	fmt.Println()
	
	text := `⟪⌊Advanced usage notes⌋⟫

Counts include with and without BOM variants.
8-bit text files are likely Windows 1252/Latin-1/ANSI or OEM 850/OEM 437, there is no detailed analysis.

⌊EOL styles⌋:
- ¬⟪Windows⟫: ⟦\r\n⟧
- ¬⟪Unix⟫: ⟦\n⟧
- ¬⟪Mac⟫: ⟦\r⟧
 
⌊Warnings report⌋:
- ¬Empty files
- ¬Source text files (based on extension) that should contain text, but with unrecognized content
- ¬UTF-8 files with BOM
- ¬UTF-16 files without BOM
- ¬Different encodings for a given file type (extension) in a directory
- ¬Mixed EOL styles in a file
- ¬Different EOL styles for a given file type (extension) in a directory`

	MyMarkup.RenderMarkup(strings.ReplaceAll(text, "{APP_NAME}", APP_NAME))
	fmt.Println()
	MyMarkup.RenderMarkup(MyGlob.GlobSyntax())
}

func NewOptions() (*Options, error) {
	options := &Options{
		Autorecurse: true,
	}

	flag.Usage = func() {
		usage()
	}

	showHelp := flag.Bool("h", false, "Show this message")
	showHelp2 := flag.Bool("?", false, "Show this message")
	showExtendedHelp := flag.Bool("??", false, "Show advanced usage notes")
	autorecurseStr := flag.String("a", "+", "Enable (+) or disable (-) glob autorecurse mode")
	autorecursePlus := flag.Bool("a+", false, "Synonym for -a +")
	autorecurseMinus := flag.Bool("a-", false, "Synonym for -a -")
	flag.BoolVar(&options.ShowOnlyWarnings, "w", false, "Only show warnings")
	flag.BoolVar(&options.Verbose, "v", false, "Verbose output")

	flag.Parse()

	if *showHelp || *showHelp2 || flag.NArg() > 0 && (flag.Args()[0] == "?" || flag.Args()[0] == "help") {
		usage()
		os.Exit(0)
	}

	if *showExtendedHelp || flag.NArg() > 0 && flag.Args()[0] == "??" {
		extendedUsage()
		os.Exit(0)
	}

	// Separated option/argument: -a + and -a -
	switch *autorecurseStr {
	case "+":
		options.Autorecurse = true
	case "-":
		options.Autorecurse = false
	default:
		return nil, fmt.Errorf("Only -a+ and -a- (enable/disable autorecurse) are supported")
	}

	// Joined options, -a+ and -a-
	if *autorecursePlus {
		options.Autorecurse = true
	}
	if *autorecurseMinus {
		options.Autorecurse = false
	}


	options.Sources = flag.Args()

	return options, nil
}
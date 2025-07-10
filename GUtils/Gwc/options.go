// gwc options.go
// Parse and validate command line options, returning a clean Options struct
//
// 2027-07-10	PV 		First version

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
	ShowOnlyTotal bool
	Verbose          bool
}

func header() {
	fmt.Printf("%s %s\n", APP_NAME, APP_VERSION)
	fmt.Println(APP_DESCRIPTION)
}

func usage() {
	header()
	fmt.Println()
	text := `⌊Usage⌋: {APP_NAME} ¬[⦃?⦄|⦃-?⦄|⦃-h⦄|⦃??⦄|⦃-??⦄] [⦃-a+⦄|⦃-a-⦄] [⦃-t⦄] [⦃-v⦄] [⟨source⟩...]

⌊Options⌋:
⦃?⦄|⦃-?⦄|⦃-h⦄  ¬Show this message
⦃??⦄|⦃-??⦄   ¬Show advanced usage notes
⦃-a+⦄|⦃-a-⦄  ¬Enable (default) or disable glob autorecurse mode (see extended usage)
⦃-t⦄       ¬Only show total line
⦃-v⦄       ¬Verbose output
⟨source⟩   ¬File or directory to search, glob syntax supported (see extended usage). Without source, search stdin.`

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

The four numerical fields report lines, words, characters and bytes counts. For UTF-8 or UTF-16 encoded files, a character is a Unicode codepoint, so bytes and characters counts may be different. Characters count neither include line terminators, nor BOM if present. Bytes count is the total file size as reported by the operating system, including line terminators and BOM if present.

Words are series of character(s) separated by space(s), spaces are either ASCII 9 (tab) or 32 (regular space).  Unicode "fancy spaces" are not considered.

Lines end with ⟦\r⟧, ⟦\n⟧ or ⟦\r\n⟧. If the last line of the file ends with such termination character, an extra empty line is counted.`

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
	flag.BoolVar(&options.ShowOnlyTotal, "t", false, "Only show total line")
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
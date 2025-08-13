// gwc options.go
// Parse and validate command line options, returning a clean Options struct
//
// 2025-07-10	PV 		First version

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
	Pattern        string
	Sources        []string
	IgnoreCase     bool
	WholeWord      bool
	FixedString    bool
	OutLevel       int
	ShowPath       bool 	// Set to true by main if there is more than 1 file to search from
	Autorecurse    bool
	Verbose        bool
}

var ShowMatchCount bool
var ShowMatchPath  bool

func header() {
	fmt.Printf("%s %s\n", APP_NAME, APP_VERSION)
	fmt.Println(APP_DESCRIPTION)
}

func usage() {
	header()
	fmt.Println()
	text := `⌊Usage⌋: {APP_NAME} ¬[⦃?⦄|⦃-?⦄|⦃-h⦄|⦃??⦄|⦃-??⦄] [⦃-i⦄] [⦃-w⦄] [⦃-F⦄] [⦃-a+⦄|⦃-a-⦄] [⦃-v⦄] [⦃-c⦄] [⦃-l⦄] ⟨pattern⟩ [⟨source⟩...]

⌊Options⌋:
⦃?⦄|⦃-?⦄|⦃-h⦄  ¬Show this message
⦃??⦄|⦃-??⦄   ¬Show advanced usage notes
⦃-v⦄       ¬Verbose output
⦃-i⦄       ¬Ignore case during search
⦃-w⦄       ¬Whole word search
⦃-F⦄       ¬Fixed string search (no regexp interpretation), also for patterns starting with - ? or help
⦃-a+⦄|⦃-a-⦄  ¬Enable (default) or disable glob autorecurse mode (see extended usage)
⦃-c⦄       ¬Suppress normal output, show count of matching lines for each file
⦃-l⦄       ¬Suppress normal output, show matching file names only
⟨pattern⟩  ¬Regular expression to search
⟨source⟩   ¬File or directory to search, glob syntax supported. Without source, search stdin`

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

Options ⦃-c⦄ (show count of matching lines) and ⦃-l⦄ (show matching file names only) can be used together to show matching lines count only for matching files.
Put special characters such as ⟦.⟧, ⟦*⟧ or ⟦?⟧ between brackets such as ⟦[.]⟧, ⟦[*]⟧ or ⟦[?]⟧ to search them as is.
To search for ⟦[⟧ or ⟦]⟧, use ⟦[\\[]⟧ or ⟦[\\]]⟧.
To search for a string containing double quotes, surround string by double quotes, and double individual double quotes inside. To search for ⟦\"msg\"⟧: {APP_NAME} ⟦\"\"\"msg\"\"\"⟧ ⟦C:\\Sources\\**\\*.rs⟧
To search for the string help, use option ⦃-F⦄: {APP_NAME} ⦃-F⦄ help ⟦C:\\Sources\\**\\*.go⟧`

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
	flag.BoolVar(&options.IgnoreCase, "i", false, "Ignore case during search")
	flag.BoolVar(&options.WholeWord, "w", false, "Whole word search")
	flag.BoolVar(&options.FixedString, "F", false, "Fixed string search")
	autorecurseStr := flag.String("a", "+", "Enable (+) or disable (-) glob autorecurse mode")
	autorecursePlus := flag.Bool("a+", false, "Synonym for -a +")
	autorecurseMinus := flag.Bool("a-", false, "Synonym for -a -")
	flag.BoolVar(&ShowMatchCount, "c", false, "Show count of matching lines for each file")
	flag.BoolVar(&ShowMatchPath, "l", false, "Show matching file names only")
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

	if ShowMatchPath {
		options.OutLevel |= 1
	}
	if ShowMatchCount {
		options.OutLevel |= 2
	}

	for _, arg := range flag.Args() {
		if strings.HasPrefix(arg, "-") {
			return nil, fmt.Errorf("Invalid/unknown option %s", arg)
		}

		if len(options.Pattern) == 0 {
			options.Pattern = arg
		} else {
			options.Sources = append(options.Sources, arg)
		}
	}

	if len(options.Pattern) == 0 {
		header()
		fmt.Printf("\nNo pattern specified.\nUse %s ? to show options or %s ?? for advanced usage notes.\n", APP_NAME, APP_NAME)
		return nil, fmt.Errorf("")
	}

	return options, nil
}

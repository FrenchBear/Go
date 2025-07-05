package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/PieVio/MyMarkup"
	"github.com/PieVio/MyGlob"
)

type Options struct {
	Sources          []string
	Autorecurse      bool
	ShowOnlyWarnings bool
	Verbose          bool
}

func header() {
	fmt.Printf("%s %s\n", APP_NAME, APP_VERSION)
	fmt.Println("Text type information in Rust")
}

func usage() {
	header()
	fmt.Println()
	text := "⌊Usage⌋: {APP_NAME} ¬[⦃?⦄|⦃-?⦄|⦃-h⦄|⦃??⦄] [⦃-a+⦄|⦃-a-⦄] [⦃-w⦄] [⦃-v⦄] [⟨source⟩...]\n\n⌊Options⌋:\n⦃?⦄|⦃-?⦄|⦃-h⦄  ¬Show this message\n⦃??⦄       ¬Show advanced usage notes\n⦃-a+⦄|⦃-a-⦄  ¬Enable (default) or disable glob autorecurse mode (see extended usage)\n⦃-w⦄       ¬Only show warnings\n⦃-v⦄       ¬Verbose output\n⟨source⟩   ¬File or directory to search, glob syntax supported. Without source, search stdin."

	MyMarkup.RenderMarkup(strings.ReplaceAll(text, "{APP_NAME}", APP_NAME))
}

func extendedUsage() {
	header()
	text := "Copyright ©2025 Pierre Violent\n\n⟪⌊Advanced usage notes⌋⟫\n\nCounts include with and without BOM variants.\n8-bit text files are likely Windows 1252/Latin-1/ANSI or OEM 850/OEM 437, there is no detailed analysis.\n\n⌊EOL styles:⌋\n- ¬⟪Windows⟫: ⟦\\r\\n⟧\n- ¬⟪Unix⟫: ⟦\\n⟧\n- ¬⟪Mac⟫: ⟦\\r⟧\n\n⌊Warnings report:⌋\n- ¬Empty files\n- ¬Source text files (based on extension) that should contain text, but with unrecognized content\n- ¬UTF-8 files with BOM\n- ¬UTF-16 files without BOM\n- ¬Different encodings for a given file type (extension) in a directory\n- ¬Mixed EOL styles in a file\n- ¬Different EOL styles for a given file type (extension) in a directory"

	MyMarkup.RenderMarkup(strings.ReplaceAll(text, "{APP_NAME}", APP_NAME))
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
	flag.BoolVar(&options.ShowOnlyWarnings, "w", false, "Only show warnings")
	flag.BoolVar(&options.Verbose, "v", false, "Verbose output")

	flag.Parse()

	if *showHelp || *showHelp2 {
		usage()
		return nil, fmt.Errorf("")
	}

	if *showExtendedHelp {
		extendedUsage()
		return nil, fmt.Errorf("")
	}

	switch *autorecurseStr {
	case "+":
		options.Autorecurse = true
	case "-":
		options.Autorecurse = false
	default:
		return nil, fmt.Errorf("Only -a+ and -a- (enable/disable autorecurse) are supported")
	}

	options.Sources = flag.Args()

	return options, nil
}
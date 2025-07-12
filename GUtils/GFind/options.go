// gfind options

// 2025-07-12 	PV 		First version

package main

import (
	"fmt"
	"strings"

	"github.com/PieVio/MyGlob"
	"github.com/PieVio/MyMarkup"
	"github.com/PieVio/TextAutoDecode"
)

type Options struct {
	sources       []string
	Actions_names map[string]bool
	search_files  bool
	search_dirs   bool
	names         []string
	isempty       bool
	recycle       bool
	autorecurse   bool
	noaction      bool
	verbose       bool
}

func header() {
	fmt.Printf("%s %s\n", APP_NAME, APP_VERSION)
	fmt.Println(APP_DESCRIPTION)
}

func usage() {
	header()
	fmt.Println()
	text := `⌊Usage⌋: {APP_NAME} ¬[⦃?⦄|⦃-?⦄|⦃-h⦄|⦃??⦄] [⦃-v⦄] [⦃-n⦄] [⦃-f⦄|⦃-type f⦄|⦃-d⦄|⦃-type d⦄] [⦃-e⦄|⦃-empty⦄] [⦃-r+⦄|⦃-r-⦄] [⦃-a+⦄|⦃-a-⦄] [⟨action⟩...] [⦃-name⦄ ⟨name⟩] ⟨source⟩...

⌊Options⌋:
⦃?⦄|⦃-?⦄|⦃-h⦄          ¬Show this message
⦃??⦄               ¬Show advanced usage notes
⦃-v⦄               ¬Verbose output
⦃-n⦄               ¬No action: display actions, but don't execute them
⦃-f⦄|⦃-type f⦄       ¬Search for files
⦃-d⦄|⦃-type d⦄       ¬Search for directories
⦃-e⦄|⦃-empty⦄        ¬Only find empty files or directories
⦃-r+⦄|⦃-r-⦄          ¬Delete to recycle bin (default) or delete forever; Recycle bin is not allowed on network sources
⦃-a+⦄|⦃-a-⦄          ¬Enable (default) or disable glob autorecurse mode (see extended usage)
⦃-name⦄ ⟨name⟩       ¬Append ⟦**/⟧⟨name⟩ to each source directory (compatibility with XFind/Search)
⟨source⟩           ¬File or directory to search

⌊Actions⌋:
⦃-print⦄           ¬Default, print matching files names and dir names
⦃-dir⦄             ¬Variant of ⦃-print⦄, with last modification date and size
⦃-delete⦄          ¬Delete matching files
⦃-rmdir⦄           ¬Delete matching directories, whether empty or not`

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

⌊Compatibility with XFind⌋:
- ¬Option ⦃-norecycle⦄ can be used instead of ⦃-r-⦄ to indicate to delete forever.
- ¬Option -name can be used to indicate a specific file name to search.`

	MyMarkup.RenderMarkup(strings.ReplaceAll(text, "{APP_NAME}", APP_NAME))
	fmt.Println()
	MyMarkup.RenderMarkup(MyGlob.GlobSyntax())
}	

func NewOptions() (*Options, error) {

	return &Options{}, nil
}
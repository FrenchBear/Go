// gfind options

// 2025-07-12 	PV 		First version
// 2025-07-13 	PV 		Option -nop

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/PieVio/MyGlob"
	"github.com/PieVio/MyMarkup"
	"github.com/PieVio/TextAutoDecode"
)

type Options struct {
	sources       []string
	actions_names map[string]bool
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
⦃-r+⦄|⦃-r-⦄          ¬Delete to recycle bin or delete forever (default); Recycle bin is not allowed on network sources
⦃-a+⦄|⦃-a-⦄          ¬Enable (default) or disable glob autorecurse mode (see extended usage)
⦃-name⦄ ⟨name⟩       ¬Append ⟦**/⟧⟨name⟩ to each source directory (compatibility with XFind/Search)
⟨source⟩           ¬File or directory to search

⌊Actions⌋:
⦃-print⦄           ¬Default, print matching files names and dir names
⦃-dir⦄             ¬Variant of ⦃-print⦄, with last modification date and size
⦃-nop[rint]⦄       ¬Do nothing, useful to replace default action ⦃-print⦄ to count files and folders with option ⦃-v⦄
⦃-delete⦄          ¬Delete matching files. ⚠ In this go version, all files are permanently deleted
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
- ¬Option ⦃-name⦄ can be used to indicate a specific file name or pattern to search.`

	MyMarkup.RenderMarkup(strings.ReplaceAll(text, "{APP_NAME}", APP_NAME))
	fmt.Println()
	MyMarkup.RenderMarkup(MyGlob.GlobSyntax())
}

func NewOptions() (*Options, error) {
	opt := Options{autorecurse: true, recycle: true}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if strings.HasPrefix(arg, "-") {
			// Options are case insensitive
			argls := strings.ToLower(arg[1:])

			switch argls {
			case "?", "h", "help", "-help":
				usage()
				return nil, fmt.Errorf("")
			case "v":
				opt.verbose = true
			case "n", "noaction":
				opt.noaction = true
			case "f":
				opt.search_files = true
			case "d":
				opt.search_dirs = true
			case "type":
				if i == len(os.Args)-1 {
					return nil, fmt.Errorf("Option -type requires an argument f or d")
				}
				i++
				argopt := os.Args[i]
				switch argopt {
				case "f":
					opt.search_files = true
				case "d":
					opt.search_dirs = true
				default:
					return nil, fmt.Errorf("Option -type requires an argument f or d")
				}

			case "name":
				if i == len(os.Args)-1 {
					return nil, fmt.Errorf("Option -type requires an argument f or d")
				}
				i++
				argopt := os.Args[i]
				opt.names = append(opt.names, argopt)

			case "e", "empty":
				opt.isempty = true

			case "r+", "recycle":
				opt.recycle = true
			case "r-", "norecycle":
				opt.recycle = false

			case "a+":
				opt.autorecurse = true
			case "a-":
				opt.autorecurse = false

			case "print":
				opt.actions_names = map[string]bool{"print": true}
			case "dir":
				opt.actions_names = map[string]bool{"dir": true}
			case "nop", "noprint":
				opt.actions_names = map[string]bool{"nop": true}
			case "rm", "del", "delete":
				opt.actions_names = map[string]bool{"delete": true}
			case "rd", "rmdir":
				opt.actions_names = map[string]bool{"rmdir": true}

			default:
				return nil, fmt.Errorf("Invalid/unsupported option %s", arg)

			}
		} else {
			switch strings.ToLower(arg) {
			case "?", "h", "help", "-help":
				usage()
				return nil, fmt.Errorf("")
			case "??":
				extendedUsage()
				return nil, fmt.Errorf("")
			default:
				opt.sources = append(opt.sources, arg)
			}
		}
	}

	// If neither filtering files or dirs has been requested, then we search for both
	if !opt.search_dirs && !opt.search_files {
		opt.search_dirs = true
		opt.search_files = true
	}

	// If no action is specified, then print action is default
	if len(opt.actions_names) == 0 {
		opt.actions_names = map[string]bool{"print": true}
	}

	return &opt, nil
}

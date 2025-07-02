// mymarkup_main.go
// Main package (simple testing) for package MyMarkup
//
// go mod edit -replace github.com/PieVio/MyMarkup=C:\Development\GitHub\Go\Packages\MyMarkup
// go mod tidy
//
// 2025-06-23	PV		First version

package main

import (
	"fmt"

	MyMarkup "github.com/PieVio/MyMarkup"
)

func main() {
	fmt.Printf("MyMarkup lib version: %s\n\n", MyMarkup.Version())

	MyMarkup.RenderMarkup("⌊Hello⌋, ⟪world⟫⦃!⦄")

	text := `⌊Usage⌋: rgrep ¬[⦃?⦄|⦃-?⦄|⦃-h⦄|⦃??⦄] [⦃-i⦄] [⦃-w⦄] [⦃-F⦄] [⦃-r⦄] [⦃-v⦄] [⦃-c⦄] [⦃-l⦄] pattern [source...]
⦃?⦄|⦃-?⦄|⦃-h⦄  ¬Show this message
⦃??⦄       ¬Show advanced usage notes
⦃-i⦄       ¬Ignore case during search
⦃-w⦄       ¬Whole word search
⦃-F⦄       ¬Fixed string search (no regexp interpretation), also for patterns starting with - ? or help
⦃-r⦄       ¬Use autorecurse, see advanced help
⦃-c⦄       ¬Suppress normal output, show count of matching lines for each file
⦃-l⦄       ¬Suppress normal output, show matching file names only
⦃-v⦄       ¬Verbose output
pattern  ¬Regular expression to search
source   ¬File or directory search, glob syntax supported. Without source, search stdin

⟪⌊Advanced usage notes⌋⟫

Counts include with and without BOM variants.
8-bit text files are likely Windows 1252/Latin-1/ANSI or OEM 850/OEM 437, there is no detailed analysis.
Files without BOM must be more than 10 characters for auto-detection of UTF-8 or UTF-16.

⌊EOL Styles⌋
- ¬⟪Windows⟫: \\r\\n
- ¬⟪Unix⟫: \\n
- ¬⟪Mac⟫: \\r

⌊Warnings reported⌋
• ¬Empty files.
• ¬Source text files (based on extension) that should contain text, but with unrecognized content.
• ¬UTF-8 files with BOM.
• ¬UTF-16 files without BOM.
• ¬Different encodings for a given file type (extension) in a directory.
• ¬Mixed EOL styles in a file.
• ¬Different EOL styles for a given file type (extension) in a directory.

⌊Glob pattern rules⌋
• ¬⟦?⟧ matches any single character.
• ¬⟦*⟧ matches any (possibly empty) sequence of characters.
• ¬⟦**⟧ matches the current directory and arbitrary subdirectories. To match files in arbitrary subdirectories, use ⟦**\\*⟧. This sequence must form a single path component, so both **a and b** are invalid and will result in an error.
• ¬⟦[...]⟧ matches any character inside the brackets. Character sequences can also specify ranges of characters, as ordered by Unicode, so e.g. ⟦[0-9]⟧ specifies any character between 0 and 9 inclusive. Special cases: ⟦[[]⟧ represents an opening bracket, ⟦[]]⟧ represents a closing bracket. 
• ¬⟦[!...]⟧ is the negation of ⟦[...]⟧, i.e. it matches any characters not in the brackets.
• ¬The metacharacters ⟦?⟧, ⟦*⟧, ⟦[⟧, ⟦]⟧ can be matched by using brackets (e.g. ⟦[?]⟧). When a ⟦]⟧ occurs immediately following ⟦[⟧ or ⟦[!⟧ then it is interpreted as being part of, rather then ending, the character set, so ⟦]⟧ and NOT ⟦]⟧ can be matched by ⟦[]]⟧ and ⟦[!]]⟧ respectively. The ⟦-⟧ character can be specified inside a character sequence pattern by placing it at the start or the end, e.g. ⟦[abc-]⟧.
• ¬⟦{choice1,choice2...}⟧  match any of the comma-separated choices between braces. Can be nested, and include ⟦?⟧, ⟦*⟧ and character classes.
• ¬Character classes ⟦[ ]⟧ accept regex syntax such as ⟦[\\d]⟧ to match a single digit, see https://docs.rs/regex/latest/regex/#character-classes for character classes and escape sequences supported.

⌊Autorecurse glob pattern transformation⌋
• ¬⟪Constant pattern (no filter, no **⟧) pointing to a directory⟫: ⟦\\**\\*⟧ is appended at the end to search all files of all subdirectories.
• ¬⟪Patterns without ⟦**⟧ and ending with a filter⟫: ⟦\\**⟧ is inserted before final filter to find all matching files of all subdirectories.
`
	MyMarkup.RenderMarkup(text)
	test_own()
}

func test_own() {
	fmt.Println("Style Default")
	fmt.Printf("Style %sBold%s, and default\n", MyMarkup.STYLE_BOLD_ON, MyMarkup.STYLE_BOLD_OFF)
	fmt.Printf("Style %sUnderline%s, and default\n", MyMarkup.STYLE_UNDERLINE_ON, MyMarkup.STYLE_UNDERLINE_OFF)
	fmt.Printf("Style %sDimmed%s, and default\n", MyMarkup.STYLE_DIM_ON, MyMarkup.STYLE_DIM_OFF)
	fmt.Printf("Style %sItalic%s, and default\n", MyMarkup.STYLE_ITALIC_ON, MyMarkup.STYLE_ITALIC_OFF)
	fmt.Printf("Style %sUnderline%s, and default\n", MyMarkup.STYLE_UNDERLINE_ON, MyMarkup.STYLE_UNDERLINE_OFF)
	fmt.Printf("Style %sBlink%s, and default\n", MyMarkup.STYLE_BLINK_ON, MyMarkup.STYLE_BLINK_OFF)
	fmt.Printf("Style %sReverse%s, and default\n", MyMarkup.STYLE_REVERSE_ON, MyMarkup.STYLE_REVERSE_OFF)
	fmt.Printf("Style %sHidden%s, and default\n", MyMarkup.STYLE_HIDDEN_ON, MyMarkup.STYLE_HIDDEN_OFF)
	fmt.Printf("Style %sStrikethrough%s, and default\n", MyMarkup.STYLE_STRIKETHROUGH_ON, MyMarkup.STYLE_STRIKETHROUGH_OFF)

	type NameColor struct {
		name  string
		color string
	}

	fmt.Println("\nColors")
	fga := []NameColor{
		{"Black", MyMarkup.FG_BLACK},
		{"Red", MyMarkup.FG_RED},
		{"Green", MyMarkup.FG_GREEN},
		{"Yellow", MyMarkup.FG_YELLOW},
		{"Blue", MyMarkup.FG_BLUE},
		{"Magenta", MyMarkup.FG_MAGENTA},
		{"Cyan", MyMarkup.FG_CYAN},
		{"White", MyMarkup.FG_WHITE},
		{"Default", MyMarkup.FG_DEFAULT},
		{"Bright Black", MyMarkup.FG_BRIGHT_BLACK},
		{"Bright Red", MyMarkup.FG_BRIGHT_RED},
		{"Bright Green", MyMarkup.FG_BRIGHT_GREEN},
		{"Bright Yellow", MyMarkup.FG_BRIGHT_YELLOW},
		{"Bright Blue", MyMarkup.FG_BRIGHT_BLUE},
		{"Bright Magenta", MyMarkup.FG_BRIGHT_MAGENTA},
		{"Bright Cyan", MyMarkup.FG_BRIGHT_CYAN},
		{"Bright White", MyMarkup.FG_BRIGHT_WHITE},
	}

	for _, nc := range fga {
		fmt.Printf("%s%s%s\n", nc.color, nc.name, MyMarkup.FG_DEFAULT)
	}
}

// MyGlob.go
// MyGlob package, my own implementation of glob search
//
// 2025-07-01	PV 		Converted from Rust by Gemini

package MyGlob

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	LIB_VERSION = "1.0.0"
)

// Segment is an interface for a segment of a glob pattern.
type Segment interface {
	isSegment()
}

// ConstantSegment is a constant string segment.
type ConstantSegment struct {
	Value string
}

func (c ConstantSegment) isSegment() {}

// RecurseSegment is a recurse (**) segment.
type RecurseSegment struct{}

func (r RecurseSegment) isSegment() {}

// FilterSegment is a glob filter segment, converted into a Regex.
type FilterSegment struct {
	Regex *regexp.Regexp
}

func (f FilterSegment) isSegment() {}

// MyGlobSearch is the main struct of MyGlob.
type MyGlobSearch struct {
	root        string
	segments    []Segment
	ignoreDirs  []string
	isConstant  bool
}

// MyGlobBuilder is used to build a MyGlobSearch object.
type MyGlobBuilder struct {
	globPattern string
	ignoreDirs  []string
	autorecurse bool
}

// MyGlobError represents an error returned by MyGlob.
type MyGlobError struct {
	Message string
}

func (e MyGlobError) Error() string {
	return e.Message
}

// Version returns the library version.
func Version() string {
	return LIB_VERSION
}

// GlobSyntax returns the glob pattern syntax documentation.
func GlobSyntax() string {
	return `⌊Glob pattern rules⌋:
- ¬⟦?⟧ matches any single character.
- ¬⟦*⟧ matches any (possibly empty) sequence of characters.
- ¬⟦**⟧ matches the current directory and arbitrary subdirectories. To match files in arbitrary subdirectories, use ⟦**/*⟧. This sequence must form a single path component, so both ⟦**a⟧ and ⟦b**⟧ are invalid and will result in an error.
- ¬⟦[...]⟧ matches any character inside the brackets. Character sequences can also specify ranges of characters (Unicode order), so ⟦[0-9]⟧ specifies any character between 0 and 9 inclusive. Special cases: ⟦[[]⟧ represents an opening bracket, ⟦[]]⟧ represents a closing bracket. 
- ¬⟦[!...]⟧ is the negation of ⟦[...]⟧, it matches any characters not in the brackets.
- ¬The metacharacters ⟦?⟧, ⟦*⟧, ⟦[⟧, ⟦]⟧ can be matched by escaping them between brackets such as ⟦[\?]⟧ or ⟦[[\]]⟧. When a ⟦]⟧ occurs immediately following ⟦[⟧ or ⟦[!⟧ then it is interpreted as being part of, rather then ending, the character set, so ⟦]⟧ and NOT ⟦]⟧ can be matched by ⟦[]]⟧ and ⟦[!]]⟧ respectively. The ⟦-⟧ character can be specified inside a character sequence pattern by placing it at the start or the end, e.g. ⟦[abc-]⟧.
- ¬⟦{choice1,choice2...}⟧  match any of the comma-separated choices between braces. Can be nested, and include ⟦?⟧, ⟦*⟧ and character classes.
- ¬Character classes ⟦[ ]⟧ accept regex syntax such as ⟦[\d]⟧ to match a single digit, see https://docs.rs/regex/latest/regex/#character-classes for character classes and escape sequences supported.

⌊Autorecurse glob pattern transformation⌋:
- ¬⟪Constant pattern⟫ (no filter, no ⟦**⟧) pointing to a directory: ⟦/**/*⟧ is appended at the end to search all files of all subdirectories.
- ¬⟪Patterns without ⟦**⟧ and ending with a filter⟫: ⟦/**⟧ is inserted before final filter to find all matching files of all subdirectories.`
}

// New creates a new MyGlobBuilder.
func New(globPattern string) *MyGlobBuilder {
	return &MyGlobBuilder{
		globPattern: globPattern,
		ignoreDirs: []string{
			"$recycle.bin",
			"system volume information",
			".git",
		},
	}
}

// AddIgnoreDir adds a directory to the ignore list.
func (b *MyGlobBuilder) AddIgnoreDir(dir string) *MyGlobBuilder {
	b.ignoreDirs = append(b.ignoreDirs, strings.ToLower(dir))
	return b
}

// Autorecurse sets the autorecurse flag.
func (b *MyGlobBuilder) Autorecurse(active bool) *MyGlobBuilder {
	b.autorecurse = active
	return b
}

// Compile builds a new MyGlobSearch from the builder.
func (b *MyGlobBuilder) Compile() (*MyGlobSearch, error) {
	if b.globPattern == "" {
		return nil, MyGlobError{"Glob pattern can't be empty"}
	}
	if strings.HasSuffix(b.globPattern, "\\") || strings.HasSuffix(b.globPattern, "/") {
		return nil, MyGlobError{"Glob pattern can't end with \\ or /"}
	}

	dirSep := string(os.PathSeparator)
	glob := b.globPattern + dirSep

	var cut, pos int
	for i, c := range glob {
		if strings.ContainsRune("*?[{", c) {
			break
		}
		if c == '/' || c == '\\' {
			cut = i
		}
		pos = i + 1
	}

	root := glob[:cut]
	if root == "" {
		root = "."
	}

	var segments []Segment
	var err error
	if pos < len(glob) {
		var globPart string
		if cut == 0 {
			globPart = glob
		} else {
			globPart = glob[cut+1:]
		}
		segments, err = globToSegments(globPart)
		if err != nil {
			return nil, err
		}
	}

	if b.autorecurse {
		if len(segments) == 0 {
			if fi, err := os.Stat(root); err == nil && fi.IsDir() {
				segments = append(segments, RecurseSegment{})
				re, _ := regexp.Compile("(?i)^.*$")
				segments = append(segments, FilterSegment{re})
			}
		} else {
			hasRecurse := false
			for _, s := range segments {
				if _, ok := s.(RecurseSegment); ok {
					hasRecurse = true
					break
				}
			}
			if !hasRecurse {
				if _, ok := segments[len(segments)-1].(FilterSegment); ok {
					insertIndex := len(segments) - 1
					segments = append(segments, nil)
					copy(segments[insertIndex+1:], segments[insertIndex:])
					segments[insertIndex] = RecurseSegment{}
				}
			}
		}
	}

	return &MyGlobSearch{
		root:        root,
		segments:    segments,
		ignoreDirs:  b.ignoreDirs,
		isConstant:  len(segments) == 0,
	}, nil
}

func globToSegments(globPattern string) ([]Segment, error) {
	var segments []Segment
	regexBuffer := ""
	constantBuffer := ""
	braceDepth := 0
	iter := []rune(globPattern)
	i := 0

	for i < len(iter) {
		c := iter[i]
		i++

		if c != '\\' && c != '/' {
			constantBuffer += string(c)
		}

		switch c {
		case '*':
			regexBuffer += ".*"
		case '?':
			regexBuffer += "."
		case '{':
			braceDepth++
			regexBuffer += "("
		case ',':
			if braceDepth > 0 {
				regexBuffer += "|"
			} else {
				regexBuffer += string(c)
			}
		case '}':
			braceDepth--
			if braceDepth < 0 {
				return nil, MyGlobError{"Extra closing }"}
			}
			regexBuffer += ")"
		case '\\', '/':
			if braceDepth > 0 {
				return nil, MyGlobError{fmt.Sprintf("Invalid %c between { }", c)}
			}

			if constantBuffer == "**" {
				segments = append(segments, RecurseSegment{})
			} else if strings.Contains(constantBuffer, "**") {
				return nil, MyGlobError{fmt.Sprintf("Glob pattern ** must be alone between %c", c)}
			} else if strings.ContainsAny(constantBuffer, "*?[{") {
				if braceDepth > 0 {
					return nil, MyGlobError{"Unclosed {"}
				}
				re, err := regexp.Compile("(?i)^" + regexBuffer + "$")
				if err != nil {
					return nil, err
				}
				segments = append(segments, FilterSegment{re})
			} else {
				segments = append(segments, ConstantSegment{constantBuffer})
			}
			regexBuffer = ""
			constantBuffer = ""
		case '[':
			regexBuffer += "["
			depth := 1
			if i < len(iter) && iter[i] == '!' {
				i++
				regexBuffer += "^"
			}
		bracketLoop:
			for i < len(iter) {
				innerC := iter[i]
				i++
				switch innerC {
				case ']':
					regexBuffer += "]"
					depth--
					if depth == 0 {
						break bracketLoop
					}
				case '\\':
					if i < len(iter) {
						regexBuffer += "\\" + string(iter[i])
						i++
					} else {
						regexBuffer += "\\"
					}
				default:
					regexBuffer += string(innerC)
				}
			}
		case '.', '+', '(', ')', '|', '^', '$':
			regexBuffer += "\\" + string(c)
		default:
			regexBuffer += string(c)
		}
	}

	if regexBuffer != "" {
		return nil, MyGlobError{"Invalid glob pattern"}
	}

	if len(segments) > 0 {
		if _, ok := segments[len(segments)-1].(RecurseSegment); ok {
			re, _ := regexp.Compile("(?i)^.*$")
			segments = append(segments, FilterSegment{re})
		}
	}

	return segments, nil
}

// MyGlobMatch represents a match from a glob search.
type MyGlobMatch struct {
	Path string
	Err  error
	IsDir bool
}

// Explore returns a channel of matches.
func (gs *MyGlobSearch) Explore() <-chan MyGlobMatch {
	ch := make(chan MyGlobMatch)
	go func() {
		defer close(ch)
		var stack []searchPendingData
		if len(gs.segments) == 0 {
			fi, err := os.Stat(gs.root)
			if err != nil {
				ch <- MyGlobMatch{Err: err}
				return
			}
			if fi.IsDir() {
				ch <- MyGlobMatch{Path: gs.root, IsDir: true}
			} else {
				ch <- MyGlobMatch{Path: gs.root}
			}
			return
		}

		stack = append(stack, searchPendingData{path: gs.root, depth: 0})
		for len(stack) > 0 {
			item := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if item.depth >= len(gs.segments) {
				continue
			}

			segment := gs.segments[item.depth]
			switch s := segment.(type) {
			case ConstantSegment:
				newPath := filepath.Join(item.path, s.Value)
				fi, err := os.Stat(newPath)
				if err == nil {
					if item.depth == len(gs.segments)-1 {
						if fi.IsDir() {
							ch <- MyGlobMatch{Path: newPath, IsDir: true}
						} else {
							ch <- MyGlobMatch{Path: newPath}
						}
					} else {
						if fi.IsDir() {
							stack = append(stack, searchPendingData{path: newPath, depth: item.depth + 1})
						}
					}
				}
				if item.recurse {
					entries, err := os.ReadDir(item.path)
					if err != nil {
						ch <- MyGlobMatch{Err: err}
						continue
					}
					for _, entry := range entries {
						if entry.IsDir() {
							p := filepath.Join(item.path, entry.Name())
							fnlc := strings.ToLower(entry.Name())
							isIgnored := false
							for _, ignored := range gs.ignoreDirs {
								if ignored == fnlc {
									isIgnored = true
									break
								}
							}
							if !isIgnored {
								stack = append(stack, searchPendingData{path: p, depth: item.depth, recurse: true})
							}
						}
					}
				}
			case RecurseSegment:
				stack = append(stack, searchPendingData{path: item.path, depth: item.depth + 1, recurse: true})
			case FilterSegment:
				entries, err := os.ReadDir(item.path)
				if err != nil {
					ch <- MyGlobMatch{Err: err}
					continue
				}
				var dirs []string
				for _, entry := range entries {
					fname := entry.Name()
					if entry.IsDir() {
						flnc := strings.ToLower(fname)
						isIgnored := false
						for _, ignored := range gs.ignoreDirs {
							if ignored == flnc {
								isIgnored = true
								break
							}
						}
						if !isIgnored {
							if s.Regex.MatchString(fname) {
								newPath := filepath.Join(item.path, fname)
								if item.depth == len(gs.segments)-1 {
									ch <- MyGlobMatch{Path: newPath, IsDir: true}
								} else {
									stack = append(stack, searchPendingData{path: newPath, depth: item.depth + 1})
								}
							}
							dirs = append(dirs, filepath.Join(item.path, fname))
						}
					} else {
						if item.depth == len(gs.segments)-1 && s.Regex.MatchString(fname) {
							ch <- MyGlobMatch{Path: filepath.Join(item.path, fname)}
						}
					}
				}
				if item.recurse {
					for _, dir := range dirs {
						stack = append(stack, searchPendingData{path: dir, depth: item.depth, recurse: true})
					}
				}
			}
		}
	}()
	return ch
}

type searchPendingData struct {
	path    string
	depth   int
	recurse bool
}

func (gs *MyGlobSearch) IsConstant() bool {
	return gs.isConstant
}
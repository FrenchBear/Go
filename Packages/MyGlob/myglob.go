// MyGlob.go
// MyGlob package, my own implementation of glob search
//
// 2025-07-01	PV 		Converted from Rust by Gemini
// 2025-07-12	PV 		1.1.0 Accepts tapperns ending with / or \, and special case for "?:\"
// 2025-08-11	PV 		1.2.0 Use getRoot function to separate constant root prefix from segments
// 2025-08-18	PV 		1.3.0 SetChannelSize method
// 2025-09-07	PV 		1.4.0 MaxDepth; IsConstant removed
// 2025-09-08	PV 		1.5.0 Replaced stack by a queue for more natural output order

package MyGlob

import (
	"container/list"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	LIB_VERSION = "1.5.0"
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

// FilterSegment is a glob filter segment, converted into a Regexp.
type FilterSegment struct {
	Regexp *regexp.Regexp
}

func (f FilterSegment) isSegment() {}

// MyGlobSearch is the main struct of MyGlob.
type MyGlobSearch struct {
	root       string
	segments   []Segment
	ignoreDirs []string
	maxDepth   int
	//	isConstant  bool
	channelSize int
}

// MyGlobBuilder is used to build a MyGlobSearch object.
type MyGlobBuilder struct {
	globPattern string
	ignoreDirs  []string
	maxDepth    int
	autoRecurse bool
	channelSize int
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
- ¬The metacharacters ⟦?⟧, ⟦*⟧, ⟦[⟧, ⟦]⟧ can be matched by escaping them between brackets such as ⟦[\?]⟧ or ⟦[\[]⟧. When a ⟦]⟧ occurs immediately following ⟦[⟧ or ⟦[!⟧ then it is interpreted as being part of, rather than ending the character set, so ⟦]⟧ and NOT ⟦]⟧ can be matched by ⟦[]]⟧ and ⟦[!]]⟧ respectively. The ⟦-⟧ character can be specified inside a character sequence pattern by placing it at the start or the end, e.g. ⟦[abc-]⟧.
- ¬⟦{choice1,choice2...}⟧  match any of the comma-separated choices between braces. Can be nested, and include ⟦?⟧, ⟦*⟧ and character classes.
- ¬Character classes ⟦[ ]⟧ accept regexp syntax such as ⟦[\d]⟧ to match a single digit, see https://pkg.go.dev/regexp/syntax for character classes and escape sequences supported.

⌊Autorecurse glob pattern transformation⌋:
- ¬⟪Constant pattern⟫ (no filter, no ⟦**⟧) pointing to a directory: ⟦/**/*⟧ is appended at the end to search all files of all subdirectories.
- ¬⟪Patterns without ⟦**⟧ and ending with a filter⟫: ⟦/**⟧ is inserted before the final filter to find all matching files of all subdirectories.`
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
		channelSize: 1, // Default buffer size
	}
}

// AddIgnoreDir adds a directory to the ignore list.
func (b *MyGlobBuilder) AddIgnoreDir(dir string) *MyGlobBuilder {
	b.ignoreDirs = append(b.ignoreDirs, strings.ToLower(dir))
	return b
}

// Set maxdepth, counted from ** segment, 0 means no limit (default)
func (b *MyGlobBuilder) MaxDepth(depth int) *MyGlobBuilder {
	b.maxDepth = depth
	return b
}

// Autorecurse sets the autorecurse flag.
func (b *MyGlobBuilder) Autorecurse(active bool) *MyGlobBuilder {
	b.autoRecurse = active
	return b
}

// Autorecurse sets the autorecurse flag.
func (b *MyGlobBuilder) ChannelSize(size int) *MyGlobBuilder {
	if size <= 0 {
		size = 1 // Default buffer size
	}
	b.channelSize = size
	return b
}

// getRoot separates a constant root prefix from the rest of a glob pattern.
// This is a direct translation of the provided Rust function's logic.
func getRoot(globPattern string) (root, remainder string) {
	glob := globPattern
	// Instead of an error, treat an empty pattern as "*", similar to shell command behavior.
	if glob == "" {
		glob = "*"
	}

	// Find the end of the constant prefix, which is the position of the first glob metacharacter.
	specialCharIdx := strings.IndexAny(glob, "*?[{")

	// Case 1: The pattern contains no special characters.
	// The entire string is the root, and there is no remainder.
	if specialCharIdx == -1 {
		return glob, ""
	}

	// Case 2: The pattern contains special characters.
	// We search for a path separator only within the constant part of the pattern.
	prefix := glob[:specialCharIdx]
	lastSeparatorIdx := strings.LastIndexAny(prefix, "/\\")

	if lastSeparatorIdx == -1 {
		// No path separator was found in the constant prefix.
		// The root is the current directory ".", and the remainder is the entire pattern.
		root = "."
		remainder = glob
	} else {
		// A path separator was found.
		// The root is everything up to and including that last separator.
		cutPoint := lastSeparatorIdx + 1
		root = glob[:cutPoint]
		remainder = glob[cutPoint:]
	}

	return root, remainder
}

// Compile builds a new MyGlobSearch from the builder.
func (b *MyGlobBuilder) Compile() (*MyGlobSearch, error) {
	root, rem := getRoot(b.globPattern)

	var segments []Segment
	var err error
	if rem != "" {
		segments, err = globToSegments(rem)
		if err != nil {
			return nil, err
		}
	}

	if b.autoRecurse {
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
		root:       root,
		segments:   segments,
		ignoreDirs: b.ignoreDirs,
		maxDepth:   b.maxDepth,
		//		isConstant:  len(segments) == 0,
		channelSize: b.channelSize,
	}, nil
}

func globToSegments(globPattern string) ([]Segment, error) {
	// Make sure that pattern ends with path separator to simplyfy code
	dirSep := string(os.PathSeparator)
	if !strings.HasSuffix(globPattern, "/") && !strings.HasSuffix(globPattern, "\\") {
		// If not, append the OS-specific separator.
		globPattern += dirSep
	}

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
	Path  string
	Err   error
	IsDir bool
}

type searchPendingData struct {
	path          string
	depth         int
	recurse       bool
	recurse_depth int
}

// Explore returns a channel of matches.
func (gs *MyGlobSearch) Explore() <-chan MyGlobMatch {
	ch := make(chan MyGlobMatch, gs.channelSize)
	go func() {
		defer close(ch)

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

		queue := list.New()
		queue.PushBack(searchPendingData{path: gs.root, depth: 0})
		for queue.Len() > 0 {
			item := queue.Front().Value.(searchPendingData)
			queue.Remove(queue.Front())		// Need to call remove, there is no PopFront
			// It's a O(1) operation since queue is actially a dequeue. Remove takes an element pointer, so it just
			// needs to update next/previous pointers of previous/next elements

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
							queue.PushBack(searchPendingData{path: newPath, depth: item.depth + 1})
						}
					}
				}
				if item.recurse && (gs.maxDepth == 0 || item.recurse_depth < gs.maxDepth) {
					for direntry := range readDirStream(item.path, true) {
						if direntry.Err != nil {
							ch <- MyGlobMatch{Err: direntry.Err}
							continue
						}
						entry := direntry.Entry

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
								queue.PushBack(searchPendingData{path: p, depth: item.depth, recurse: true, recurse_depth: item.recurse_depth + 1})
							}
						}
					}
				}

			case RecurseSegment:
				queue.PushBack(searchPendingData{path: item.path, depth: item.depth + 1, recurse: true, recurse_depth: 0})

			case FilterSegment:
				var dirs []string
				for direntry := range readDirStream(item.path, false) {
					if direntry.Err != nil {
						ch <- MyGlobMatch{Err: direntry.Err}
						continue
					}
					entry := direntry.Entry

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
							if s.Regexp.MatchString(fname) {
								if gs.maxDepth == 0 || item.recurse_depth < gs.maxDepth {
									newPath := filepath.Join(item.path, fname)
									if item.depth == len(gs.segments)-1 {
										ch <- MyGlobMatch{Path: newPath, IsDir: true}
									} else {
										queue.PushBack(searchPendingData{path: newPath, depth: item.depth + 1})
									}
								}
							}
							dirs = append(dirs, filepath.Join(item.path, fname))
						}
					} else {  // File
						if item.depth == len(gs.segments)-1 && s.Regexp.MatchString(fname) {
							ch <- MyGlobMatch{Path: filepath.Join(item.path, fname)}
						}
					}
				}
				if item.recurse && (gs.maxDepth==0 || item.recurse_depth < gs.maxDepth) {
					for _, dir := range dirs {
						queue.PushBack(searchPendingData{path: dir, depth: item.depth, recurse: true, recurse_depth: item.recurse_depth + 1})
					}
				}
			}
		}
	}()
	return ch
}


// func (gs *MyGlobSearch) IsConstant() bool {
// 	return gs.isConstant
// }

// grepiterator.go
// Iterates over lines of a text matching some pattern
//
// 2025-08-13   PV (converted from Rust by Gemini)

package main

import (
	"regexp"
	"strings"
)

// GrepMatchRange represents the start and end position of a match
// within a single line of text. It's the Go equivalent of Rust's Range<usize>.
type GrepMatchRange struct {
	Start int
	End   int
}

// GrepLineMatches contains a line of text and all the regex matches found within it.
// This is the item that will be sent over the channel.
type GrepLineMatches struct {
	Line   string
	Ranges []GrepMatchRange
}

// Grep iterates over a text, finds lines matching the given regular expression,
// and returns a channel that yields each matching line along with all its match ranges.
// It is the Go equivalent of the Rust `GrepLineMatches::new` iterator.
func Grep(txt string, re *regexp.Regexp) <-chan GrepLineMatches {
	// Create a channel to send the results back.
	// The caller will read from this channel.
	ch := make(chan GrepLineMatches)

	// Start a goroutine to do the processing. This allows the Grep function
	// to return the channel immediately without blocking.
	go func() {
		// Ensure the channel is closed when the goroutine finishes.
		// This is crucial to signal the end of the stream to the receiver.
		defer close(ch)

		// Find all non-overlapping matches in the text at once.
		// FindAllStringIndex returns a slice of [start, end] byte indices.
		allMatches := re.FindAllStringIndex(txt, -1)
		if len(allMatches) == 0 {
			return // No matches, close the channel and exit.
		}

		// State for the current line being processed.
		var currentLineMatches GrepLineMatches
		// Start index of the line currently being processed.
		// -1 acts as a sentinel value indicating we haven't started the first line yet.
		currentLineStartIx := -1

		// Iterate over each match found.
		for _, match := range allMatches {
			matchStart := match[0]

			// Find the beginning of the line for the current match.
			// Search backwards from the match start for a newline character.
			lineStartIx := 0
			if lastNewlineIx := strings.LastIndexByte(txt[:matchStart], '\n'); lastNewlineIx != -1 {
				lineStartIx = lastNewlineIx + 1
			}

			// Check if this match is on a new line.
			if lineStartIx != currentLineStartIx {
				// If this isn't the very first line we're processing,
				// it means we've just finished a previous line. Send it.
				if currentLineStartIx != -1 {
					ch <- currentLineMatches
				}

				// Now, start processing the new line.
				// Find the end of this new line.
				lineEndIx := len(txt)
				if nextNewlineIx := strings.IndexByte(txt[lineStartIx:], '\n'); nextNewlineIx != -1 {
					lineEndIx = lineStartIx + nextNewlineIx
				}

				// Reset the state for the new line.
				currentLineStartIx = lineStartIx
				currentLineMatches = GrepLineMatches{
					Line:   strings.TrimRight(txt[lineStartIx:lineEndIx], "\r\n"),
					Ranges: []GrepMatchRange{},
				}
			}

			// Add the current match's range to the list for the current line.
			// The range is relative to the start of the line, not the whole text.
			currentLineMatches.Ranges = append(currentLineMatches.Ranges, GrepMatchRange{
				Start: match[0] - currentLineStartIx,
				End:   match[1] - currentLineStartIx,
			})
		}

		// After the loop, the last processed line hasn't been sent yet.
		// Send it now.
		ch <- currentLineMatches
	}()

	// Return the channel to the caller.
	return ch
}

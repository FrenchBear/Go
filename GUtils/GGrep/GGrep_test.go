// Tests for GGrep
//
// 2025-07-10 	PV 		First version

package main

import (
	"regexp"
	"testing"
)

func assert_eq[T comparable](t *testing.T, a, b T) {
	if a != b {
		t.Errorf("Expected %v, got %v", b, a)
	}
}

func TestGrepIterator(t *testing.T) {
	text := `Go is a statically typed, compiled programming language designed at Google.
Its syntax is loosely based on C, but with memory safety, garbage collection,
structural typing, and CSP-style concurrency. The language is often referred to as Golang.
This is another line about go or GO or even gO.
Final line without the word.`

	// Compile the regex. (?i) makes it case-insensitive.
	re := regexp.MustCompile(`(?im)go`)

	ch := Grep(text, re)

	lm := <-ch
	assert_eq(t, lm.Line, "Go is a statically typed, compiled programming language designed at Google.")
	assert_eq(t, len(lm.Ranges), 2)
	assert_eq(t, lm.Ranges[0].Start, 0)
	assert_eq(t, lm.Ranges[0].End, 2)
	assert_eq(t, lm.Ranges[1].Start, 68)
	assert_eq(t, lm.Ranges[1].End, 70)

	lm = <-ch
	assert_eq(t, lm.Line, "structural typing, and CSP-style concurrency. The language is often referred to as Golang.")
	assert_eq(t, len(lm.Ranges), 1)
	assert_eq(t, lm.Ranges[0].Start, 83)
	assert_eq(t, lm.Ranges[0].End, 85)

	lm = <-ch
	assert_eq(t, lm.Line, "This is another line about go or GO or even gO.")
	assert_eq(t, len(lm.Ranges), 3)
	assert_eq(t, lm.Ranges[0].Start, 27)
	assert_eq(t, lm.Ranges[0].End, 29)
	assert_eq(t, lm.Ranges[1].Start, 33)
	assert_eq(t, lm.Ranges[1].End, 35)
	assert_eq(t, lm.Ranges[2].Start, 44)
	assert_eq(t, lm.Ranges[2].End, 46)
}

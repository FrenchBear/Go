//go:build !windows

// osspecific_non_windows.go
// Non-windows specific code
//
// 2026-07-02	PV 		First version, also first example of os-specific compilation

package main

import (
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
)

func is_hidden_folder(path string) (bool, bool) {
	_, d := filepath.Split(path)
	return strings.HasPrefix(d, "."), false
}

// Sortmethod is ignored in Linux, it's always sorted using casefold
func path_comparer(_sortmethod int, s1, s2 string) int {
	// The caser for folding. It's stateless and safe for concurrent use.
	// It's efficient to create it once and reuse it.
	folder := cases.Fold()

	str1 := folder.String(s1)
	str2 := folder.String(s2)
	return strings.Compare(str1, str2)
}


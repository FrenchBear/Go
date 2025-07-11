//go:build !windows

// reparse_non_windows.go
// Alternate support for reparse points, symbolic links dir and junctions on non-windows systems
//
// 2025-07-11	PV		First version

package main

import (
	"os"
	"path/filepath"
)

func IsSymLinkD(path string) (bool, string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, "", err
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return false, "", nil
	}

	temp, err1 := os.Readlink(path)
	target, err2 := filepath.EvalSymlinks(temp)
	if err1 != nil || err2 != nil {
		return false, "", err
	}

	return true, target, nil
}

// IsJunction checks if the given path is an NTFS junction point.
func IsJunction(path string) (bool, string, error) {
	return false, "", nil
}


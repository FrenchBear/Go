// osspecific_windows.go
// Windows-specific code
//
// 2026-07-02	PV 		First version, also first example of os-specific compilation

//go:build windows

package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/text/cases"
)

func is_hidden_folder(path string) (bool, bool) {
	// Convert the path to a UTF-16 pointer for Windows API.
	// Use windows.UTF16PtrFromString for safer conversion if available (golang.org/x/sys/windows)
	// Otherwise, syscall.UTF16PtrFromString.
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, false
	}

	// Get file attributes using GetFileAttributesW.
	attributes, err := syscall.GetFileAttributes(pathPtr)
	if err != nil {
		return false, false
	}

	_, d := filepath.Split(path)

	// Check the HIDDEN and SYSTEM bits.
	isHidden := (attributes & syscall.FILE_ATTRIBUTE_HIDDEN) != 0 || strings.HasPrefix(d, ".") 
	isSystem := (attributes & syscall.FILE_ATTRIBUTE_SYSTEM) != 0 || ( (attributes & syscall.FILE_ATTRIBUTE_HIDDEN) != 0 && strings.HasPrefix(d, "$"))

	return isHidden, isHidden && isSystem
}

// Define the signature of StrCmpLogicalW
// int StrCmpLogicalW(LPCWSTR psz1, LPCWSTR psz2);
// Returns:
//   -1 if psz1 comes before psz2
//    0 if psz1 is identical to psz2
//    1 if psz1 comes after psz2
var (
	shlwapi        = syscall.NewLazyDLL("shlwapi.dll")
	strCmpLogicalW = shlwapi.NewProc("StrCmpLogicalW")
)

func path_comparer(sortmethod int, s1, s2 string) int {
	switch sortmethod {
	case 2: return path_comparer_CaseFold(s1, s2)
	default: return path_comparer_StrCmpLogicalW(s1, s2)
	}
}

// StrCmpLogicalWGo is a Go wrapper for the Windows API StrCmpLogicalW function.
// It compares two strings using the natural sort algorithm, similar to Windows File Explorer.
func path_comparer_StrCmpLogicalW(s1, s2 string) int {
	// Convert Go strings to null-terminated UTF-16 pointers for Windows API
	// syscall.UTF16PtrFromString allocates memory that needs to be freed
	// It's safer to use the x/sys/windows package for this.
	// For simplicity in this example, we'll use syscall directly,
	// but be aware of potential memory considerations in very high-performance loops.
	// In most cases, the garbage collector will handle it.

	// Create UTF-16 pointers
	p1, err := syscall.UTF16PtrFromString(s1)
	if err != nil {
		// Handle error: perhaps log or panic, depending on your application's needs.
		// For a comparison function, panicking might be acceptable if input is guaranteed valid.
		panic(fmt.Sprintf("Failed to convert string 1 to UTF16: %v", err))
	}
	p2, err := syscall.UTF16PtrFromString(s2)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert string 2 to UTF16: %v", err))
	}

	// Call the underlying Windows API function
	// The uintptr(unsafe.Pointer(p1)) converts the pointer to a uintptr,
	// which is what NewProc.Call expects for arguments.
	ret, _, _ := strCmpLogicalW.Call(uintptr(unsafe.Pointer(p1)), uintptr(unsafe.Pointer(p2)))

	// StrCmpLogicalW returns:
	// < 0 if psz1 comes before psz2
	//   0 if psz1 is identical to psz2
	// > 0 if psz1 comes after psz2
	// We cast the result to int64 first to handle potential negative return values correctly.
	return int(int32(ret))
}

func path_comparer_CaseFold(s1, s2 string) int {
	// The caser for folding. It's stateless and safe for concurrent use.
	// It's efficient to create it once and reuse it.
	folder := cases.Fold()

	str1 := folder.String(s1)
	str2 := folder.String(s2)
	return strings.Compare(str1, str2)
}

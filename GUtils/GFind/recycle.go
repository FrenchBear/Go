// recycle.go
// Support for Windows Recycle Bin
// Works with local files and folders, remote files are permanently deleted
//
// 2025-08-13	PV		First version

package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Constants from the Windows API needed for SHFileOperation
const (
	FO_DELETE          = 0x0003 // Operation: Delete the file
	FOF_ALLOWUNDO      = 0x0040 // Flag: Allow the operation to be undone (sends to Recycle Bin)
	FOF_NOCONFIRMATION = 0x0010 // Flag: Don't ask for confirmation
)

// SHFILEOPSTRUCT is the structure required by the SHFileOperation function.
// It defines the operation to be performed.
type SHFILEOPSTRUCT struct {
	hwnd   windows.HWND
	wFunc  uint32
	pFrom  *uint16
	pTo    *uint16
	fFlags uint16
	// The rest of the struct fields are not needed for this operation.
	_ bool
	_ uintptr
	_ *uint16
}

// recycleFile moves the specified file to the Windows Recycle Bin.
func recycleFile(path string) error {
	// The path must be double-null-terminated.
	// We convert the Go string to a UTF-16 pointer and add the extra null terminator.
	absPath, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	// Create a buffer that is large enough to hold the path and two null terminators.
	// The second null terminator is what makes it a "double-null-terminated" string.
	buffer := make([]uint16, len(path)+2)
	copy(buffer, unsafe.Slice(absPath, len(path)))
	// The slice is already zero-initialized, so the extra nulls are there.

	// Load the shell32.dll library and get the SHFileOperation function.
	shell32 := windows.NewLazySystemDLL("shell32.dll")
	shFileOp := shell32.NewProc("SHFileOperationW") // 'W' for wide-character (UTF-16) version

	// Prepare the SHFILEOPSTRUCT structure.
	op := &SHFILEOPSTRUCT{
		hwnd:   0, // No owner window
		wFunc:  FO_DELETE,
		pFrom:  &buffer[0],
		pTo:    nil, // Not needed for a delete operation
		fFlags: FOF_ALLOWUNDO | FOF_NOCONFIRMATION,
	}

	// Call the function.
	// The return value indicates success (0) or an error code.
	ret, _, err := shFileOp.Call(uintptr(unsafe.Pointer(op)))
	if ret != 0 {
		// If the function call itself returns an error status from the OS.
		return fmt.Errorf("SHFileOperation failed with code %d: %w", ret, err)
	}

	return nil
}
// reparse_windows.go
// Support for NTFS reparse points, symbolic links dir and junctions
//
// 2025-07-03	PV		First version, refactored gemini code

//go:build windows

package main

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// reparseDataBuffer is the structure that receives the reparse point data.
// We define it ourselves to match the Windows API structure.
// Note: golang.org/x/sys/windows.REPARSE_DATA_BUFFER can also be used,
// but defining it helps understand the layout.
type reparseDataBuffer struct {
	ReparseTag        uint32
	ReparseDataLength uint16
	Reserved          uint16
	// The following fields are part of a union in C, so we access them
	// via a buffer starting at this point.
	PathBuffer [1]uint16
}

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
	pathUTF16Ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, "", err
	}

	// Get file attributes
	attr, err := windows.GetFileAttributes(pathUTF16Ptr)
	if err != nil {
		return false, "", err
	}

	// Check if the path is a reparse point. If not, it can't be a junction.
	if attr&windows.FILE_ATTRIBUTE_REPARSE_POINT == 0 {
		return false, "", nil
	}

	// It's a reparse point, now we need to check if it's a junction.
	// Open the file/directory handle with flags to access the reparse point data.
	fd, err := windows.CreateFile(
		pathUTF16Ptr,
		0, 	// windows.GENERIC_READ,	// If this parameter is zero, the application can query certain metadata such as file, directory, or device attributes without accessing that file or device, even if GENERIC_READ access would have been denied
		windows.FILE_SHARE_READ,
		nil,
		windows.OPEN_EXISTING,
		// These flags are crucial:
		// FILE_FLAG_BACKUP_SEMANTICS is needed for directories.
		// FILE_FLAG_OPEN_REPARSE_POINT ensures we open the link itself, not the target.
		windows.FILE_FLAG_BACKUP_SEMANTICS|windows.FILE_FLAG_OPEN_REPARSE_POINT,
		0,
	)
	if err != nil {
		return false, "", err
	}
	defer windows.CloseHandle(fd)

	// Create a buffer to hold the reparse data.
	// The buffer needs to be large enough for the reparse data structure and the path.
	// MAX_PATH is 260, times 2 for UTF-16, plus the struct overhead. 1024 is safe.
	buffer := make([]byte, windows.MAXIMUM_REPARSE_DATA_BUFFER_SIZE)
	var bytesReturned uint32

	// Use DeviceIoControl to get the reparse point data.
	err = windows.DeviceIoControl(
		fd,
		windows.FSCTL_GET_REPARSE_POINT,
		nil,
		0,
		&buffer[0],
		uint32(len(buffer)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return false, "", err
	}

	// Interpret the buffer as our reparseDataBuffer struct.
	// unsafe.Pointer is needed for this type of low-level conversion.
	rdb := (*reparseDataBuffer)(unsafe.Pointer(&buffer[0]))

	// Check the reparse tag. For junctions, it's IO_REPARSE_TAG_MOUNT_POINT.
	if rdb.ReparseTag != windows.IO_REPARSE_TAG_MOUNT_POINT {
		return false, "", nil
	}

	// The path information for a junction starts after the header.
	// The structure in C is a union, but for a junction (mount point),
	// it contains SubstituteNameOffset, SubstituteNameLength,
	// PrintNameOffset, and PrintNameLength, followed by the PathBuffer.
	//
	// For simplicity, we can calculate the start of the path buffer.
	// The path starts at an offset inside the generic PathBuffer.
	// Let's find the Substitute Name, which is the actual target.
	// The offset is relative to the start of the PathBuffer field.

	mySubstituteNameOffset := *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(rdb)) + 8))
	mySubstituteNameLength := *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(rdb)) + 10))

	myPathSlice := (*[1024]uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(rdb)) + 16+uintptr(mySubstituteNameOffset)))
	myTarget := syscall.UTF16ToString(myPathSlice[:mySubstituteNameLength/2+4])
	return true, strings.TrimPrefix(myTarget, `\??\`), nil
}

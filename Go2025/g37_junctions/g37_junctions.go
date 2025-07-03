// g37_junctions.go
// Learning go, System programming, Detect NTFS junctions
//
// 2025-07-03	PV		First version (Gemini)

package main

import (
	"fmt"
	"os"
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

// IsJunction checks if the given path is an NTFS junction point.
func IsJunction(path string) (bool, error) {
	pathUTF16Ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, err
	}

	// Get file attributes
	attr, err := windows.GetFileAttributes(pathUTF16Ptr)
	if err != nil {
		return false, err
	}

	// Check if the path is a reparse point. If not, it can't be a junction.
	if attr&windows.FILE_ATTRIBUTE_REPARSE_POINT == 0 {
		return false, nil
	}

	// It's a reparse point, now we need to check if it's a junction.
	// Open the file/directory handle with flags to access the reparse point data.
	fd, err := windows.CreateFile(
		syscall.StringToUTF16Ptr(path),
		windows.GENERIC_READ,
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
		return false, err
	}
	defer windows.CloseHandle(fd)

	// Create a buffer to hold the reparse data.
	// The buffer needs to be large enough for the reparse data structure and the path.
	// MAX_PATH is 260, times 2 for UTF-16, plus the struct overhead. 1024 is safe.
	buffer := make([]byte, 1024)
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
		return false, err
	}

	// Interpret the buffer as our reparseDataBuffer struct.
	// unsafe.Pointer is needed for this type of low-level conversion.
	rdb := (*reparseDataBuffer)(unsafe.Pointer(&buffer[0]))

	// Check the reparse tag. For junctions, it's IO_REPARSE_TAG_MOUNT_POINT.
	if rdb.ReparseTag == windows.IO_REPARSE_TAG_MOUNT_POINT {
		return true, nil
	}

	return false, nil
}

// ReadJunction reads the target path of an NTFS junction.
// It returns the target path and an error if the path is not a valid junction.
func ReadJunction(path string) (string, error) {
	// First, verify it's a junction.
	isJunc, err := IsJunction(path)
	if err != nil {
		return "", fmt.Errorf("error checking if path is a junction: %w", err)
	}
	if !isJunc {
		return "", fmt.Errorf("path is not a junction: %s", path)
	}

	pathUTF16Ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}


	// The logic to get the reparse data is the same as in IsJunction.
	// We repeat it here to be self-contained, but in a real app, you'd refactor.
	fd, err := windows.CreateFile(
		pathUTF16Ptr,
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS|windows.FILE_FLAG_OPEN_REPARSE_POINT,
		0,
	)
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(fd)

	buffer := make([]byte, windows.MAXIMUM_REPARSE_DATA_BUFFER_SIZE)
	var bytesReturned uint32

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
		return "", err
	}

	// Cast the buffer to the reparse data structure.
	rdb := (*reparseDataBuffer)(unsafe.Pointer(&buffer[0]))

	// The path information for a junction starts after the header.
	// The structure in C is a union, but for a junction (mount point),
	// it contains SubstituteNameOffset, SubstituteNameLength,
	// PrintNameOffset, and PrintNameLength, followed by the PathBuffer.
	//
	// For simplicity, we can calculate the start of the path buffer.
	// The path starts at an offset inside the generic PathBuffer.
	// Let's find the Substitute Name, which is the actual target.
	// The offset is relative to the start of the PathBuffer field.
	substituteNameOffset := unsafe.Offsetof(rdb.PathBuffer) + 4 // offset for SubstituteNameOffset/Length
	substituteNameLength := *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(rdb)) + substituteNameOffset + 2))

	// The path itself starts at a further offset
	pathOffset := substituteNameOffset + 4 // an additional 4 for PrintNameOffset/Length

	// Get a slice of the uint16 (UTF-16) characters.
	pathSlice := (*[1024]uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(rdb)) + pathOffset))
	
	// Convert the UTF-16 slice to a Go string.
	target := syscall.UTF16ToString(pathSlice[:substituteNameLength/2+4])

	// The target path is usually prefixed with "\??\". We remove it to get a clean path.
	// Example: "\??\C:\Users\Default" becomes "C:\Users\Default"
	
	cleanTarget := strings.TrimPrefix(target, `\??\`)

	return cleanTarget, nil
}


func main() {
	// Let's test the paths from your DIR output
	pathsToTest := []string{
		`C:\Development`, // A Junction
		`C:\Tmp`,          // A Junction
		`C:\DocumentsOD`,  // A Directory Symlink
		`C:\Program Files`, // A regular directory
	}
    
    // Create dummy files/links for testing if they don't exist
    // This requires admin privileges to create symlinks/junctions
    fmt.Println("NOTE: For accurate testing, these paths should exist as described.")
    fmt.Println("You may need to run 'mklink /J C:\\Tmp C:\\Temp' as an example.")
    fmt.Println("--------------------------------------------------")

	for _, path := range pathsToTest {
		fi, err := os.Lstat(path)
		if err != nil {
			fmt.Printf("Path: %s\n  Error: %v\n\n", path, err)
			continue
		}

		fmt.Printf("Path: %s\n", path)

		// Using standard Go library to check for symlinks
		if fi.Mode()&os.ModeSymlink != 0 {
			target, _ := os.Readlink(path)
			fmt.Printf("  - Type: Symbolic Link (detected by Go's os.ModeSymlink)\n")
			fmt.Printf("  - Target: %s\n", target)
		} else {
			// Check if it's a junction using our custom function
			isJunc, err := IsJunction(path)
			if err != nil {
				fmt.Printf("  - Error checking for junction: %v\n", err)
				continue
			}

			if isJunc {
				target, err := ReadJunction(path)
				if err != nil {
					fmt.Printf("  - Error reading junction target: %v\n", err)
				} else {
					fmt.Printf("  - Type: Junction (detected via Windows API)\n")
					fmt.Printf("  - Target: %s\n", target)
				}
			} else {
				fmt.Println("  - Type: Regular Directory or File")
			}
		}
		fmt.Println()
	}
}

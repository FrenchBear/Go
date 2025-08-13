// g46_recycle_bin.go
// Learning go, System programming, Delete a file to recycle bin on Windows
//
// 2025-08-13	PV		First version

package main

import (
	"fmt"
	"log"
	"os"
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

func main() {
	fmt.Printf("Go Delete file to Recycle Bin on Windows\n\n")

	// test_local_file()
	// test_local_directory()
	test_network_file()
}

func test_network_file() {
	// Network files are parmanently deleted, no warning, no error
	filePath := `\\teraz\temp\temp_file_to_delete.txt`
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create temp network file: %v", err)
	}
	file.WriteString("This is a temporary file.")
	file.Close()
	fmt.Printf("✅ Created a temporary network file: '%s'\n", filePath)

	// 2. Move the file to the Recycle Bin.
	fmt.Println("Attempting to move the network file to the Recycle Bin...")
	err = recycleFile(filePath)
	if err != nil {
		log.Fatalf("❌ Failed to move network file to Recycle Bin: %v", err)
	}
}

func test_local_directory() {
	// Also work for (non-empty) directories
	err := recycleFile(`S:\Search1 - Copy`)
	if err != nil {
		log.Fatalf("❌ Failed to move file to Recycle Bin: %v", err)
	}
}

func test_local_file() {
	// 1. Create a dummy file to delete.
	filePath := "temp_file_to_delete.txt"
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create dummy file: %v", err)
	}
	file.WriteString("This is a temporary file.")
	file.Close()
	fmt.Printf("✅ Created a temporary file: '%s'\n", filePath)

	// 2. Move the file to the Recycle Bin.
	fmt.Println("Attempting to move the file to the Recycle Bin...")
	err = recycleFile(filePath)
	if err != nil {
		log.Fatalf("❌ Failed to move file to Recycle Bin: %v", err)
	}

	// 3. Confirm the result.
	fmt.Println("✅ Successfully moved file to the Recycle Bin!")
	fmt.Println("Please check your Recycle Bin to confirm.")

	// Verify the file no longer exists at the original path.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("File '%s' no longer exists at the original location.\n", filePath)
	} else {
		fmt.Printf("File '%s' still exists. The operation might have failed silently.\n", filePath)
	}

}

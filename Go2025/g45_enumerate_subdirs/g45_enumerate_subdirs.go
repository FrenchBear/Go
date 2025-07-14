// g45_enumerate_subdirs.go
// Learning go, System programming, Enumerating (large) diretories
//
// 2025-07-14	PV		First version

/*
os.ReadDir is really slow on large directories (>10000 entries), since it has to read all entries in
an array first, and call stat on each entry before returning, blocking calling thread for seconds until it's done.
Function readDirStream is an attempted optimization reading entrines by slices of 100, and returning results in
a channel: it works in a secondary thread, start returning results immediately, and use less memory, but since
it also calls os.Stat, at the end it's takes as long to enumerate a directory.
*/

/*
Gemini question:

I have serious performance problems enumerating large directories (>10000 entries) in Go on Windows. I just need to get
a list of subdirectories (just subdirectories, not links, junctions...), nothing else, and I don't need any directory
attribute, metadata, modification time or whatever.

If I use os.ReadDir, I've to wait that all entries have been read, and the function calls os.Stat on each entry, it
takes several seconds to execute.

If I use dir := os.Open on a directory and then dir.ReadDir, I can read entries on chuck of entries at a time, says 100
entries for instance, but it also calls os.Stat on each entry, which means that at the end, reading the whole directory
takes as long as os.ReadDir.

Running either one in a separate goroutine does not improve my performance issue since I've to read the whole directory,
and processing time outside goroutine is negligible.

And finally, dir.ReadDirNames is fast... but it returns together all names without specifying whether it's a file or a
directory, so I've to call os.Stat for each entry, which is way too slow.


Both Rust and C# have methods to get quickly a list of subdirectories on Windows, not calling stat, just based by the
return of Win32 directory enumeration call: this call returns all names including a boolean IsDirectory, it's extremely
fast.


How can I achieve the same result in Go? Go is supposed to be good at system programming, but if a simple directory
enumeration to get a list of subdirectories is 3-5 times slower than C# or Rust, it seriously rules out Go to write
efficient system tools based on directories analysis...



Gemini solution:

You can achieve the same high-speed performance in Go by sidestepping the standard library and directly interfacing with
the Windows API. This is accomplished using the golang.org/x/sys/windows package, which provides the necessary tools to
call low-level system functions.

Tapping into the Windows API for High-Speed Subdirectory Enumeration

The key to efficient directory scanning on Windows lies in using the FindFirstFile and FindNextFile functions from the
kernel32.dll library. These functions retrieve directory entries along with their attributes—including a flag that
identifies an entry as a directory—in a single operation, thus avoiding the overhead of separate stat calls.

Here's a Go function that efficiently lists only the subdirectories within a given path on Windows by using this direct
API approach
*/

// Problem is, while this code looks good and minimalist, calling Win32 FindFirstFile and FindNextFile, it still takes
// more than 4s to enumerate the content of C:\Windows\servicing\Packages while Rust takes 0.3s to do the same thing...
// No idea why it takes so long, folder has 11728 entries, but it's 10 times faster in rust and 5 times faster in C#
// More in Rust,  timethis rfind C:\Windows\servicing\Packages -type d  takes 65ms!!! (it really enumerates folder to
// find subdirectories)

/*
Gemini comment:

The staggering performance difference is not due to flawed Go logic but to a much deeper and more subtle issue: the high
overhead of repeated system calls in Go versus Rust.

The Real Bottleneck: Per-Call Syscall Overhead

The FindNextFile function requires a "kernel transition"—a switch from user mode (your program) to kernel mode (the
OS)—for every single file. When you do this 11,728 times in a tight loop, the cumulative cost of these transitions
becomes the main bottleneck.

- Rust's Advantage: Rust's standard library (and its underlying winapi-rs or windows-sys crates) is designed for
  zero-cost abstractions. Its Foreign Function Interface (FFI) compiles down to something extremely close to a native C
  function call, making the kernel transition as lightweight as possible.

- Go's Overhead: Go's runtime is more complex. It manages goroutines, a scheduler, and a garbage collector. The
  transition from Go code to a C-style system call involves more overhead to safely bridge these two different
  execution models. In a loop with thousands of iterations, this small per-call cost adds up to a significant amount
  of time.

So, while your Go logic is sound, you're fighting the fundamental cost of the FindNextFile loop model in Go.

The Solution: Batch Processing with the Native API

To achieve native speed, we must change the strategy entirely. Instead of asking for one file at a time, we need to ask
for a large batch of files in a single system call. This is what high-performance directory enumeration libraries (and
likely Rust's standard library) do under the hood using the NtQueryDirectoryFile function from Windows' native API
(ntdll.dll).

This function is more complex to use but drastically more efficient. It fills a large buffer with as many directory
entries as will fit, massively reducing the number of expensive kernel transitions.

Here is a practical implementation in Go that uses this high-performance batching approach. It will give you the speed
you're looking for.

*/

// Ok, but three mistakes!
// Missed field ChangeTime in fileDirectoryInformation struct
// Wrong parameters call at the end of ntQueryDirectoryFile call, should be "0, uintptr(0), 0," and not "0, 0"
// STATUS_NO_MORE_FILES is 0x80000006, not 0xC0000006
//
// Execution time: 8ms!!!!!   Excellent!

// Sources:
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/ntifs/nf-ntifs-ntquerydirectoryfile
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/ntifs/nf-ntifs-ntquerydirectoryfileex (FILE_INFORMATION_CLASS values)
// https://learn.microsoft.com/en-us/windows-hardware/drivers/ddi/ntifs/ns-ntifs-_file_directory_information (output struct)

package main

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// FILE_DIRECTORY_INFORMATION is one of the structures that NtQueryDirectoryFile can return.
// We only care about the FileName and FileAttributes.
type fileDirectoryInformation struct {
	NextEntryOffset uint32
	FileIndex       uint32
	CreationTime    int64
	LastAccessTime  int64
	LastWriteTime   int64
	ChangeTime      int64
	EndOfFile       int64
	AllocationSize  int64
	FileAttributes  uint32
	FileNameLength  uint32
	FileName        [1]uint16
}

// getSubdirectoriesBatch is a high-performance implementation using NtQueryDirectoryFile.
func getSubdirectoriesBatch(path string) ([]string, error) {
	// 1. Get a handle to the directory
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}

	handle, err := windows.CreateFile(
		p,
		windows.FILE_LIST_DIRECTORY,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS, // Necessary for opening a directory handle
		0,
	)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(handle)

	var subdirectories []string
	// 2. Create a large buffer to receive file data in batches
	// A 64KB buffer is a reasonable size.
	buffer := make([]byte, 64*1024) 
	
	ntdll := windows.MustLoadDLL("ntdll.dll")
	ntQueryDirectoryFile := ntdll.MustFindProc("NtQueryDirectoryFile")

	ne := 0

	for {
		var ioStatusBlock struct { // IO_STATUS_BLOCK
			Status, Information uintptr
		}

		// 3. Make the single, efficient syscall to get a batch of entries
		status, _, err := ntQueryDirectoryFile.Call(
			uintptr(handle),
			0, // Event
			0, // ApcRoutine
			0, // ApcContext
			uintptr(unsafe.Pointer(&ioStatusBlock)),
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(len(buffer)),
			uintptr(1), 	//windows.FileDirectoryInformation), // Specify the info class we want
			0,          // ReturnSingleEntry (FALSE)
			uintptr(0), // FileName (NULL)
			0,          // RestartScan (FALSE)
		)

		// 0x80000006 is STATUS_NO_MORE_FILES
		if status == 0x80000006 { 
			break
		}
		if status != 0 {
			fmt.Println("$99", status)
			return nil, fmt.Errorf("NtQueryDirectoryFile failed with status: %X, error: %v", status, err)
		}

		// 4. Iterate through the linked list of entries within the buffer
		offset := uintptr(0)
		for {
			ne++

			info := (*fileDirectoryInformation)(unsafe.Pointer(&buffer[offset]))

			if info.FileAttributes&windows.FILE_ATTRIBUTE_DIRECTORY != 0 {
				// We need to slice the flexible array member `FileName`
				nameSlice := unsafe.Slice(&info.FileName[0], info.FileNameLength/2)
				name := windows.UTF16ToString(nameSlice)
				if name != "." && name != ".." {
					subdirectories = append(subdirectories, name)
				}
			}

			if info.NextEntryOffset == 0 {
				break
			}
			offset += uintptr(info.NextEntryOffset)
		}
	}

	//fmt.Println("Entries:", ne)

	return subdirectories, nil
}


func main() {
	targetDir := "C:\\Windows\\servicing\\Packages"
	
	fmt.Printf("Scanning directory with batching: %s\n", targetDir)
	start := time.Now()

	dirs, err := getSubdirectoriesBatch(targetDir)
	if err != nil {
		log.Fatalf("Error listing subdirectories: %v", err)
	}
	
	duration := time.Since(start)

	fmt.Printf("Found %d subdirectories.\n", len(dirs))
	fmt.Printf("🚀 Execution time: %s\n", duration)
}








/* First attempt

package main

import ( 
	"fmt" 
	"log"
	 "path/filepath"
	  "syscall"
	   "time"

    "golang.org/x/sys/windows"
)

func main() { r := `C:\Windows\servicing\Packages`

    start := time.Now() 
    dirs, err := ListSubdirectories(r)
    if err != nil {
        log.Fatalf("Error listing subdirectories: %v", err)
    }
    elapsed := time.Since(start)

    fmt.Printf("Found %d subdirectories:\n", len(dirs))
    for _, dir := range dirs {
        fmt.Println(dir)
    }
    fmt.Printf("Elapsed time: %s\n", elapsed)
}


// ListSubdirectories returns a slice of subdirectory names in a given path. // It uses the Windows API for high
performance with large directories. func ListSubdirectories(path string) ([]string, error) { var subdirectories []string

    // Append `\*` to the path to search for all entries in the directory.
    searchPath, err := syscall.UTF16PtrFromString(filepath.Join(path, `*`))
    if err != nil {
        return nil, fmt.Errorf("error creating search path: %w", err)
    }

    var findData windows.Win32finddata
    handle, err := windows.FindFirstFile(searchPath, &findData)
    if err != nil {
        // If the directory is empty, FindFirstFile can return ERROR_FILE_NOT_FOUND.
        // This is not an error in our case, so we return an empty slice.
        if err == windows.ERROR_FILE_NOT_FOUND {
            return subdirectories, nil
        }
        return nil, fmt.Errorf("error finding first file: %w", err)
    }
    defer windows.FindClose(handle)

    ne := 0
    for {
        ne++

        // Check if the found entry is a directory.
        if findData.FileAttributes&windows.FILE_ATTRIBUTE_DIRECTORY != 0 {
            // Convert the UTF-16 file name to a Go string.
            name := windows.UTF16ToString(findData.FileName[:])

            // Ignore the "." and ".." entries.
            if name != "." && name != ".." {
                subdirectories = append(subdirectories, name)
            }
        }

        // Move to the next entry in the directory.
        err = windows.FindNextFile(handle, &findData)
        if err != nil {
            // ERROR_NO_MORE_FILES indicates the end of the directory, which is expected.
            if err == windows.ERROR_NO_MORE_FILES {
                break
            }
            return nil, fmt.Errorf("error finding next file: %w", err)
        }
    }

    fmt.Println("Entries:", ne)

    return subdirectories, nil
}

*/
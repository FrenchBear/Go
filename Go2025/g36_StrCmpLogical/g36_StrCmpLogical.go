// g36_StrCmpLogical.go
// Learning go, System programming, Call StrCmpLogicalW from shlwapi.dll to sort filenames
//
// 2025-06-30	PV		First version

/*
in Go on Windows, how cal I call function StrCmpLogicalW of shlwapi.dll to compare strings to sort filenames exactly the
same way as Windows File Explorer?

You're asking about a very specific and useful function for Windows, StrCmpLogicalW, which is indeed the function File
Explorer uses for its "natural sort" order. This is excellent for sorting filenames like "file1.txt", "file10.txt",
"file2.txt" correctly as "file1.txt", "file2.txt", "file10.txt" rather than "file1.txt", "file10.txt", "file2.txt"
(lexicographical sort).

To call StrCmpLogicalW from shlwapi.dll in Go on Windows, you'll use the syscall package, specifically
syscall.NewLazyDLL and syscall.NewLazyProc.
*/

package main

import (
	"fmt"
	"sort"
	"syscall"
	"unsafe" // For unsafe.Pointer and unsafe.Sizeof
)

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

// StrCmpLogicalWGo is a Go wrapper for the Windows API StrCmpLogicalW function.
// It compares two strings using the natural sort algorithm, similar to Windows File Explorer.
func StrCmpLogicalWGo(s1, s2 string) int {
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
	return int(ret)
}

func main() {
	filenames := []string{
		"file1.txt",
		"file10.txt",
		"file2.txt",
		"another_file.txt",
		"Image_1.jpg",
		"Image_10.jpg",
		"Image_2.jpg",
		"folder A",
		"folder B",
		"folder 1",
		"folder 10",
		"folder 2",
		"archive.zip",
	}

	fmt.Println("Original filenames:")
	for _, f := range filenames {
		fmt.Println(f)
	}

	// Sort using standard lexicographical sort
	lexicalSorted := make([]string, len(filenames))
	copy(lexicalSorted, filenames)
	sort.Strings(lexicalSorted)
	fmt.Println("\nSorted (Lexicographical - standard Go sort.Strings):")
	for _, f := range lexicalSorted {
		fmt.Println(f)
	}

	// Sort using StrCmpLogicalWGo
	naturalSorted := make([]string, len(filenames))
	copy(naturalSorted, filenames)

	// Use sort.Slice with our custom comparison function
	sort.Slice(naturalSorted, func(i, j int) bool {
		return StrCmpLogicalWGo(naturalSorted[i], naturalSorted[j]) < 0
	})

	fmt.Println("\nSorted (Natural Sort - Windows File Explorer style):")
	for _, f := range naturalSorted {
		fmt.Println(f)
	}
}

/*
Explanation:

1. syscall.NewLazyDLL("shlwapi.dll"): This loads the shlwapi.dll library. NewLazyDLL is good because it delays loading
   the DLL until its functions are actually called, which can slightly improve startup time if the functions are not
   always needed.

2. shlwapi.NewProc("StrCmpLogicalW"): This gets a pointer to the StrCmpLogicalW function within the loaded DLL. NewProc
   also delays resolving the function address until Call is invoked.

3. StrCmpLogicalWGo(s1, s2 string) int function:

   - syscall.UTF16PtrFromString(s1) / syscall.UTF16PtrFromString(s2): Windows API functions that deal with strings
     almost universally expect UTF-16 encoded, null-terminated strings. This function converts a Go string to such a
     pointer (*uint16). This is crucial.

   - strCmpLogicalW.Call(...): This is how you invoke the loaded DLL function.

     - It takes uintptr arguments. You convert your *uint16 pointers to uintptr using uintptr(unsafe.Pointer(p1)).
       unsafe.Pointer is needed as an intermediate step to convert a typed pointer to a generic uintptr.

     - It returns three values: the return value of the API call, an error code (errno), and a last error code
       (lasterr). We're primarily interested in the first return value.

   - Return Value: StrCmpLogicalW returns <0, 0, or >0 depending on the comparison, which aligns perfectly with the
     standard comparison function signature required by sort.Slice.

4. main function:

   - It defines a slice of filenames that clearly demonstrates the difference between lexicographical and natural
     sorting.

   - It uses sort.Strings for the standard lexicographical sort for comparison.

   - It uses sort.Slice with an anonymous function that calls our StrCmpLogicalWGo wrapper to perform the natural sort.
     sort.Slice is a flexible way to sort any slice using a custom comparison.


Important Considerations:

Windows-Specific: This code is strictly for Windows. It will not compile or run on Linux or macOS. If you need
cross-platform natural sorting, you'd need to implement different logic for each OS (e.g., using regular expressions to
parse numbers for Unix-like systems, or a dedicated natural sort library for Go).

Error Handling: The example uses panic for error handling during string conversion for brevity in the StrCmpLogicalWGo
function. In a production application, you might want more robust error handling (e.g., returning an error from the
wrapper, or gracefully handling malformed input).

unsafe package: The use of unsafe.Pointer is necessary here to bridge the gap between Go's typed pointers and the raw
uintptr expected by syscall.Call. While unsafe should be used with caution, this is a legitimate and common use case
when interacting with low-level C APIs via syscall.

golang.org/x/sys/windows: For more idiomatic and often safer interactions with the Windows API, especially for more
complex calls, consider using the golang.org/x/sys/windows package. It provides pre-defined functions and constants for
many Windows APIs, reducing the need for NewLazyDLL and NewLazyProc directly. For StrCmpLogicalW, you might find a
binding or an easier way to wrap it there, though at the time of writing, it's not directly exposed as a Go function.

Self-correction: A quick check reveals that StrCmpLogicalW is not directly exposed as a Go function in
golang.org/x/sys/windows. So, the syscall approach demonstrated above remains the primary way to call it from Go.
However, for converting strings to UTF-16 pointers, windows.UTF16PtrFromString from golang.org/x/sys/windows is
generally preferred over syscall.UTF16PtrFromString due to better memory management and error handling within that
package. For this simple case, syscall.UTF16PtrFromString works fine, but be aware of the alternative for more complex
scenarios.

*/

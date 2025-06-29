// g31_transmute.go
// Learning go, System programming, transform []uint8 into []uint16
//
// 2025-06-29	PV		First version

/*
In go, how can I transmute a []uint8 into a uint16[] ?

Transmuting a []uint8 into a []uint16 in Go (meaning, treating the same underlying memory as a slice of a different
type) involves using the unsafe package. This is a powerful but dangerous operation, and you should only use it when you
fully understand the implications and risks.

Why it's "unsafe":

-   Endianness: uint16 values are two bytes. The order in which those bytes are stored in memory (little-endian vs.
    big-endian) matters. If your []uint8 is in one endianness and you interpret it as []uint16 with the wrong
    assumption, your numbers will be incorrect. Go uses the system's native endianness.

-   Alignment: uint16 values generally need to be aligned on 2-byte boundaries in memory. If your []uint8 doesn't start
    at an address that's properly aligned for uint16, accessing the uint16 values can lead to crashes or performance
    issues, depending on the architecture.

-   Garbage Collector (GC): The Go garbage collector isn't aware of these "transmuted" pointers. If the original []uint8
    slice is modified or goes out of scope and gets garbage collected, your []uint16 slice will be pointing to invalid
    memory, leading to crashes.

-   Mutability: If you modify the []uint16 slice, you are directly modifying the underlying bytes of the []uint8 slice,
    and vice-versa. This can lead to unexpected behavior if not carefully managed.
*/

package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	fmt.Println("==== TRANSMUTATION TESTS")
	transmutation_tests()

	fmt.Println("==== ENDIANNESS TESTS")
	CheckEndiannessUsingBinaryPackage()
	CheckEndiannessAtRuntimeUnsafe()
}

// BytesToUint16SliceUnsafe converts a []uint8 to a []uint16 using unsafe operations.
// DANGER: Use with extreme caution. This function does not handle endianness
// and assumes proper memory alignment. Modifying the resulting slice
// directly modifies the underlying byte slice.
func BytesToUint16SliceUnsafe(b []uint8) ([]uint16, error) {
	if len(b)%2 != 0 {
		return nil, fmt.Errorf("length of byte slice must be even for uint16 conversion")
	}

	// Get the address of the underlying array of the byte slice
	bytePtr := unsafe.Pointer(&b[0])

	// Check for alignment (optional but highly recommended)
	// uint16 requires 2-byte alignment.
	if uintptr(bytePtr)%unsafe.Alignof(uint16(0)) != 0 {
		return nil, fmt.Errorf("byte slice not aligned for uint16 conversion")
	}

	// Create a new SliceHeader for uint16
	var uint16SliceHeader reflect.SliceHeader
	uint16SliceHeader.Data = uintptr(bytePtr)
	uint16SliceHeader.Len = len(b) / 2
	uint16SliceHeader.Cap = cap(b) / 2

	// Convert the SliceHeader back to a []uint16
	return *(*[]uint16)(unsafe.Pointer(&uint16SliceHeader)), nil
}

// Safer alternative: Manually convert with endianness consideration
func BytesToUint16Slice(b []uint8, byteOrder binary.ByteOrder) ([]uint16, error) {
	if len(b)%2 != 0 {
		return nil, fmt.Errorf("length of byte slice must be even for uint16 conversion")
	}

	result := make([]uint16, len(b)/2)
	for i := 0; i < len(b)/2; i++ {
		result[i] = byteOrder.Uint16(b[i*2 : (i*2)+2])
	}
	return result, nil
}

func transmutation_tests() {
	data := []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06} // Example byte slice

	// *** UNSAFE METHOD ***
	fmt.Println("--- Unsafe Conversion ---")
	uint16sUnsafe, err := BytesToUint16SliceUnsafe(data)
	if err != nil {
		fmt.Println("Error (unsafe):", err)
	} else {
		fmt.Printf("Original bytes: %x\n", data)
		fmt.Printf("Uint16s (unsafe): %x (raw hex interpretation)\n", uint16sUnsafe)

		// Demonstrate mutability
		uint16sUnsafe[0] = 0xFFFF
		fmt.Printf("Modified bytes (via unsafe uint16s): %x\n", data) // data is also modified!
	}

	// *** SAFER METHOD (Recommended) ***
	fmt.Println("\n--- Safer Conversion (with Endianness) ---")
	// For Little Endian:
	data2 := []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	uint16sLittleEndian, err := BytesToUint16Slice(data2, binary.LittleEndian)
	if err != nil {
		fmt.Println("Error (Little Endian):", err)
	} else {
		fmt.Printf("Original bytes: %x\n", data2)
		fmt.Printf("Uint16s (Little Endian): %x\n", uint16sLittleEndian)
	}
	fmt.Println()

	// For Big Endian:
	data3 := []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	uint16sBigEndian, err := BytesToUint16Slice(data3, binary.BigEndian)
	if err != nil {
		fmt.Println("Error (Big Endian):", err)
	} else {
		fmt.Printf("Original bytes: %x\n", data3)
		fmt.Printf("Uint16s (Big Endian): %x\n", uint16sBigEndian)
	}
	fmt.Println()

	// Example of an odd length slice for the safe method
	_, err = BytesToUint16Slice([]uint8{0x01, 0x02, 0x03}, binary.LittleEndian)
	fmt.Println("Error for odd length slice (safe):", err)
	fmt.Println()
}

// ----------------------------------------------
// CHeck endianness

// DetermineNativeEndianness returns true if the system is little-endian, false if big-endian.
// This is the most reliable and idiomatic way using encoding/binary.
func DetermineNativeEndianness() bool {
	var x uint16 = 1
	buf := make([]byte, 2)
	binary.NativeEndian.PutUint16(buf, x) // Write 1 into buf using native endianness

	// If the first byte is 1, it's little-endian (least significant byte first)
	return buf[0] == 1
}

// CheckEndiannessUsingBinaryPackage demonstrates the idiomatic way
// to work with endianness using the encoding/binary package.
func CheckEndiannessUsingBinaryPackage() {
	fmt.Println("--- Using encoding/binary package (Corrected & Reliable) ---")

	// Determine and print the native endianness
	isLittleEndian := DetermineNativeEndianness()
	if isLittleEndian {
		fmt.Println("Current system is Little Endian (via reliable test value method).")
	} else {
		fmt.Println("Current system is Big Endian (via reliable test value method).")
	}

	// Example usage with explicit endianness
	var val uint16 = 0xABCD
	buf2 := make([]byte, 2)

	fmt.Printf("Original uint16 value: 0x%X (%d)\n", val, val)

	// Write in Little Endian
	binary.LittleEndian.PutUint16(buf2, val)
	fmt.Printf("Little Endian bytes: %x\n", buf2)

	// Write in Big Endian
	binary.BigEndian.PutUint16(buf2, val)
	fmt.Printf("Big Endian bytes: %x\n", buf2)

	// Read using NativeEndian (this will adapt to the system's endianness)
	// You still use the NativeEndian object to read/write based on native order.
	fmt.Printf("Reading 0x%x using binary.NativeEndian: 0x%X\n", buf2, binary.NativeEndian.Uint16(buf2))
	fmt.Println()
}

// CheckEndiannessAtRuntimeUnsafe is a low-level way to determine endianness
// at runtime by inspecting memory. This is generally discouraged unless
// you have a very specific reason and understand the implications of unsafe.
func CheckEndiannessAtRuntimeUnsafe() {
	fmt.Println("\n--- Unsafe Runtime Check (Discouraged for general use, but works) ---")

	var i int32 = 0x01020304 // A multi-byte value
	// Take a pointer to the integer and cast it to a byte pointer
	p := (*byte)(unsafe.Pointer(&i))

	// If the first byte pointed to is 0x04, it's Little Endian (least significant byte first)
	// If the first byte pointed to is 0x01, it's Big Endian (most significant byte first)
	switch *p {
	case 0x04:
		fmt.Println("Current system is Little Endian (via unsafe check).")
	case 0x01:
		fmt.Println("Current system is Big Endian (via unsafe check).")
	default:
		fmt.Println("Could not determine endianness (via unsafe check).")
	}
	fmt.Println()
}

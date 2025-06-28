// g31_dir.go
// Learning go, System programming, Explore directories
//
// 2025-06-28	PV		First version

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Go Directories")

	test_sum_files_size()
}

func test_sum_files_size() {
	dir := `C:\Utils`
	s, err := sum_files_size(dir)
	if err != nil {
		fmt.Println("Error reading", dir, ": ", err)
		return
	}

	fmt.Printf("Sum of files sizes in %s: %v bytes\n", dir, formatLongWithThousandsSeparator(s))
	sK := float64(s)/1024.0
	sM := sK/1024.0
	sG := sM/1024.0
	fmt.Printf("= %s KB\n", formatFloatWithThousandsSeparator(sK, 1))
	fmt.Printf("= %s MB\n", formatFloatWithThousandsSeparator(sM, 1))
	fmt.Printf("= %s GB\n", formatFloatWithThousandsSeparator(sG, 1))
}

func sum_files_size(path string) (int64, error) {
	contents, err := os.ReadDir(path)
	if err != nil {
		return -1, err
	}

	var total int64
	for _, entry := range contents {
		// Visit directory entries
		if entry.IsDir() {
			// If we are processing a directory, we need to keep digging.
			temp, err := sum_files_size(filepath.Join(path, entry.Name()))
			if err != nil {
				return -1, err
			}
			total += temp
			// Get size of each non-directory entry
		} else {
			info, err := entry.Info()
			if err != nil {
				return -1, err
			}
			// Returns an int64 value
			total += info.Size()
		}
	}
	return total, nil
}

// formatLongWithThousandsSeparator formats an int64 with thousands separators.
func formatLongWithThousandsSeparator(n int64) string {
	s := strconv.FormatInt(n, 10)
	isNegative := false
	if n < 0 {
		isNegative = true
		s = s[1:] // Remove the '-' for processing
	}

	parts := []string{}
	for i := len(s); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{s[start:i]}, parts...)
	}

	formatted := strings.Join(parts, " ")	// Ordinary space, but could be non-breaking space
	if isNegative {
		return "-" + formatted
	}
	return formatted
}

// Use 64-bit integers
func formatFloatWithThousandsSeparator(f float64, precision int) string {
	if precision < 0 {
		precision = 0 // Default to no fractional part if precision is negative
	}

	isNegative := false
	if f < 0 {
		isNegative = true
		f = -f // Work with the absolute value
	}

	// Separate integer and fractional parts as strings
	// strconv.FormatFloat for the whole number, then split
	s := strconv.FormatFloat(f, 'f', precision, 64)
	parts := strings.Split(s, ".")

	integerPartStr := parts[0]
	var formattedInt string

	// Convert integerPartStr to an int to use formatIntWithThousandsSeparator
	// This ensures proper grouping even if the integer part itself is very large.
	int64Val, err := strconv.ParseInt(integerPartStr, 10, 64)
	if err != nil {
		// Fallback if Atoi fails (e.g., extremely large number not fitting in int)
		// This is less common but good to consider for robustness.
		formattedInt = integerPartStr // Just use the raw string
	} else {
		formattedInt = formatLongWithThousandsSeparator(int64Val)
	}

	var fractionalPartStr string
	if len(parts) > 1 {
		fractionalPartStr = "." + parts[1]
	}

	result := formattedInt + fractionalPartStr
	if isNegative {
		return "-" + result
	}
	return result
}
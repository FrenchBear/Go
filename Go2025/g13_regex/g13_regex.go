// g13_regex.go
// Learning go, Regular expressions
//
// 2025-06-11	PV		First version
// 2025-06-11	PV		Huge complement (orginal was only match_basic())

package main

import (
	"fmt"
	"regexp"
	"strconv"
)


func main() {
	fmt.Println("Regex in Go")
	fmt.Println()

	match_basic()
	match_string()
	find_first_match()
	find_all_matches()
	capturing_groups()
	find_all_submatches()
	replace_with_fixed_string()
	replace_with_function()
	splitting_strings()
	validating_ip_address()
}

func matchNameSur(s string) bool {
	t := []byte(s)
	re := regexp.MustCompile(`^[A-Z][a-z]*$`)
	return re.Match(t)
}

func matchInt(s string) bool {
t := []byte(s)
re := regexp.MustCompile(`^[-+]?\d+$`)
return re.Match(t)
}

func match_basic() {
	fmt.Println("Basic matching")

	fmt.Println(matchNameSur("Pierre"))		// true
	fmt.Println(matchNameSur("Jean-Paul"))	// false

	fmt.Println(matchInt("-355"))		// true
	fmt.Println(matchInt("12 345"))		// false

	fmt.Println()
}

func match_string() {
	fmt.Println("Matching strings")

	text := "Hello, world!"
	pattern := `Hello, (world|Go)!`

	// Compile the regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}

	// Check for a match
	if re.MatchString(text) {
		fmt.Println("String matches the pattern.")
	} else {
		fmt.Println("String does not match the pattern.")
	}

	text2 := "Hello, Go!"
	if re.MatchString(text2) {
		fmt.Println("String 2 matches the pattern.")
	} else {
		fmt.Println("String 2 does not match the pattern.")
	}

	fmt.Println()
}

func find_first_match() {
	fmt.Println("Find first match")

	text := "My email is test@example.com, and my other email is another@domain.org."
	pattern := `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`

	re := regexp.MustCompile(pattern) // MustCompile panics on error, useful for global regexes

	match := re.FindString(text)
	if match != "" {
		fmt.Println("First email found:", match)
	} else {
		fmt.Println("No email found.")
	}

	fmt.Println()
}

func find_all_matches() {
	fmt.Println("Find all matches")

	text := "Price: $10.99, Discount: $2.50, Total: $8.49"
	pattern := `\$(\d+\.\d{2})` // Matches dollar amounts

	re := regexp.MustCompile(pattern)

	matches := re.FindAllString(text, -1) // -1 means find all occurrences
	if len(matches) > 0 {
		fmt.Println("All prices found:", matches)
	} else {
		fmt.Println("No prices found.")
	}

	// Example with a limited number of matches
	limitedMatches := re.FindAllString(text, 2) // Find at most 2 occurrences
	fmt.Println("Limited matches (2):", limitedMatches)

	fmt.Println()
}

func capturing_groups() {
	fmt.Println("Capturing groups")

	text := "Name: John Doe, Age: 30"
	pattern := `Name: (.*), Age: (\d+)`

	re := regexp.MustCompile(pattern)

	// FindStringSubmatch returns a slice of strings
	// The first element is the full match, subsequent elements are capturing groups.
	match := re.FindStringSubmatch(text)

	if len(match) > 0 {
		fmt.Println("Full match:", match[0])
		if len(match) > 1 {
			fmt.Println("Name:", match[1]) // First capturing group
		}
		if len(match) > 2 {
			fmt.Println("Age:", match[2])  // Second capturing group
		}
	} else {
		fmt.Println("No match found.")
	}

	text2 := "No match here"
	match2 := re.FindStringSubmatch(text2)
	fmt.Println("Match for text2:", match2) // Will be nil

	fmt.Println()
}

func find_all_submatches() {
	fmt.Println("Find all submatches")

	text := `
		User: alice, ID: 101
		User: bob, ID: 102
		User: charlie, ID: 103
	`
	pattern := `User: (\w+), ID: (\d+)`

	re := regexp.MustCompile(pattern)

	// FindAllStringSubmatch returns a slice of slices of strings
	allMatches := re.FindAllStringSubmatch(text, -1)

	if len(allMatches) > 0 {
		for _, match := range allMatches {
			fmt.Printf("User: %s, ID: %s\n", match[1], match[2])
		}
	} else {
		fmt.Println("No matches found.")
	}

	fmt.Println()
}

func replace_with_fixed_string() {
	fmt.Println("Replace with fixed string")

	text := "Hello world, hello Go!"
	pattern := `hello`
	replaceWith := "Hi"

	re := regexp.MustCompile(pattern)

	newText := re.ReplaceAllString(text, replaceWith)
	fmt.Println("Original:", text)
	fmt.Println("Replaced:", newText)

	// Replacing only the first occurrence
	newTextOnce := re.ReplaceAllStringFunc(text, func(s string) string {
		return replaceWith
	})
	fmt.Println("Replaced once:", newTextOnce)

	fmt.Println()
}

func replace_with_function() {
	fmt.Println("Replace with function")

	text := "ItemA: 10.50, ItemB: 20.00, ItemC: 5.25"
	pattern := `Item(\w+): (\d+\.\d{2})`

	re := regexp.MustCompile(pattern)

	// Replace all prices by doubling them
	newText := re.ReplaceAllStringFunc(text, func(match string) string {
		// match is the full string matched by the regex
		// We need to use FindStringSubmatch again on 'match' if we want capturing groups
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match // Should not happen if pattern is correct
		}
		item := submatches[1]
		priceStr := submatches[2]

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return match // In case of parsing error, return original match
		}
		doubledPrice := price * 2
		return fmt.Sprintf("Item%s: %.2f", item, doubledPrice)
	})
	fmt.Println("Original:", text)
	fmt.Println("Modified (doubled prices):", newText)

	// Another example: anonymizing email addresses
	emailText := "My email is user@example.com and his is another@domain.org."
	emailPattern := `(\w+)@(\w+\.\w+)`
	emailRe := regexp.MustCompile(emailPattern)

	anonymizedEmailText := emailRe.ReplaceAllStringFunc(emailText, func(match string) string {
		submatches := emailRe.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}
		// Replace username with "******"
		return fmt.Sprintf("******@%s", submatches[2])
	})
	fmt.Println("Anonymized emails:", anonymizedEmailText)

	fmt.Println()
}

func splitting_strings() {
	fmt.Println("Splitting strings")

	text := "apple,banana;orange go,grape"
	pattern := `[,;\s]+` // Split by comma, semicolon, or whitespace

	re := regexp.MustCompile(pattern)

	parts := re.Split(text, -1) // -1 means split into all possible parts
	fmt.Println("Split parts:", parts)

	// Example with a limited number of splits
	limitedParts := re.Split(text, 2) // Split into at most 2 parts
	fmt.Println("Limited split (2 parts):", limitedParts)

	fmt.Println()
}

func validating_ip_address() {
	fmt.Println("Validating IP address")

	ip1 := "192.168.1.1"
	ip2 := "256.0.0.1" // Invalid
	ip3 := "10.0.0"   // Invalid

	// A simplified regex for IPv4 (not fully exhaustive for all edge cases like leading zeros, but good for common validation)
	// For production, consider a dedicated IP parsing library or a more robust regex.
	ipPattern := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	ipRe := regexp.MustCompile(ipPattern)

	fmt.Printf("%s is valid IP: %t\n", ip1, ipRe.MatchString(ip1))
	fmt.Printf("%s is valid IP: %t\n", ip2, ipRe.MatchString(ip2))
	fmt.Printf("%s is valid IP: %t\n", ip3, ipRe.MatchString(ip3))

	fmt.Println()
}
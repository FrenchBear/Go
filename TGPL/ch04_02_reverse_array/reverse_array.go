package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	j := [...]string{"Monday", "Tuesday", "wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	reverseArray7Str(&j)
	fmt.Println(j)

	// ds is directly a slice
	ds := []string{"Pom", "Pom", "Pom", "Pam", "Pim", "Pim", "Pum", "Pem", "Pem", "Pym"}
	fmt.Println(removeAdjacentsDuplicates(ds))

	s := "Aé♫山𝄞🐗"
	t := []byte(s)
	u := reverseUTF8ByteSlice(t)
	fmt.Printf("<%s>\n<%s>\n<%s>\n", s, string(t), string(u))
}

// Exercise 4.3: Rewrite reverse to use an array pointer instead of a slice
// Pb, array size need to be specified since it's part of the type...
func reverseArray7Str(a *[7]string) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

// Exercise 4.4: Write an in-place function to eliminate adjacent duplicates in a []string slice
func removeAdjacentsDuplicates(s []string) []string {
	j := 0
	for i, v := range s {
		if i == 0 || v != s[i-1] {
			s[j] = v
			j++
		}
	}
	return s[:j]
}

// Exercise 4.7
// Modify reverse to reverse the characters of a []byte slice that represents a UTF-8-encoded string, in place.
// Can you do it without allocating new memory?
func reverseUTF8ByteSlice(s []byte) []byte {
	s2 := make([]byte, len(s))
	copy(s2, s)
	t := len(s2)
	for i := 0; i < len(s); {
		_, size := utf8.DecodeRune(s2[i:])
		t -= size
		copy(s[t:t+size], s2[i:i+size])
		i += size
	}
	return s
}

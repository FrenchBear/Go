// g26_read_size.go
// Learning go, System programming, files, read a limited size from a text file
//
// 2025-06-23	PV		First version

package main

import (
	"fmt"
	"os"
	"unicode/utf8"
)

func main() {
	fmt.Println("Go readSize")

	filename := `C:\DocumentsOD\Doc tech\Encodings\prenoms-utf8bom.txt`

	tad, err := process(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", filename, err)
		os.Exit(1)
	}

	fmt.Println("Result: ", tad.Encoding)
}

const MILLE = 2000

type TextFileEncoding int

const (
	NotText    TextFileEncoding = iota // Binary or unrecognized text (for instance contains chars in 0..31 other than \r \n \t)
	Empty                              // File is empty
	ASCII                              // Only 7-bit characters
	EightBit                           // ANSI/Windows 1525 or other
	UTF8                               // Plain UTF-8 without BOM
	UTF8BOM                            // Starts with EF BB BF
	UTF16LE                            // No BOM but UTF-16 LE detected
	UTF16BE                            // No BOM but UTF-16 BE detected
	UTF16LEBOM                         // Starts with FF FE (Windows)
	UTF16BEBOM                         // Starts with FE FF
	InProgress
)

type TextAutoDecode struct {
	Text     string
	Encoding TextFileEncoding
}

func process(file string) (TextAutoDecode, error) {
	f, err := os.Open(file)
	if err != nil {
		return TextAutoDecode{}, err
	}
	defer f.Close()

	// Empty file?
	buffer_1000 := make([]byte, MILLE)
	n, err := f.Read(buffer_1000)
	if n == 0 {
		return TextAutoDecode{Text: "", Encoding: Empty}, nil
	}

	// UTF-8 BOM?
	// Since we have a BOM, no need to check for ASCII subset
	if n >= 3 && buffer_1000[0] == 0xEF && buffer_1000[1] == 0xBB && buffer_1000[2] == 0xBF {
		s, ok := check_utf8(buffer_1000, n)

		if !ok {
			return TextAutoDecode{Text: "", Encoding: NotText}, nil
		}

		if n < MILLE {
			return TextAutoDecode{Text: s, Encoding: UTF8BOM}, nil
		}

		// ToDo
		// return final_read(&mut buffer_full_read, &mut buffer_full, &mut file, UTF_8, Some(TextFileEncoding::UTF8BOM));
		return TextAutoDecode{Text: s, Encoding: InProgress}, nil
	}

	return TextAutoDecode{Text: "??", Encoding: InProgress}, nil
}

// check_utf8 checks if a small byte buffer of n bytes (max 1000) contains a valid UTF-8 string.
// If valid, it returns the string and true.
// If not valid, it returns an empty string and false.
// Note that if buffer len is exactly 1000, it's possible that the last UTF-8 character is truncated,
// so we reduce the buffer to be safe. Anyway, in this case, we'll reread the whole file and do a global check
func check_utf8(buffer []byte, n int) (string, bool) {
	// Basic input validation
	if buffer == nil || n < 0 {
		return "", false
	}

	nsafe := n
	if n == MILLE {
		for nsafe = MILLE-1; ; nsafe-- {
			// If buffer[nsafe] is a valid beginning for UTF-8 encoding, we can stop here
			if buffer[nsafe]<128 || (buffer[nsafe] & 0b11100000)==0b11000000 || (buffer[nsafe] & 0b11110000)==0b11100000 || (buffer[nsafe] & 0b11111000)==0b11110000 {
				break
			}
			// If it's a continuation character, we can continue, but at most three continuation characters are valid
			if (buffer[nsafe] & 0b11000000)==0b10000000 && n>=MILLE-3 {
				continue
			}
			// Sorry, that's not valid UTF-8...
			return "", false		
		}
	}

	// If last character is <128, it's not truncated and can be kept
	if buffer[nsafe]<128 {
		nsafe++
	}
	buffer_safe := buffer[:nsafe]

	// Use utf8.Valid to check if the byte slice is valid UTF-8
	if utf8.Valid(buffer_safe) {
		// If valid, convert the byte slice to a string and return true
		return string(buffer_safe), true
	}

	// If not valid, return an empty string and false
	return "", false
}

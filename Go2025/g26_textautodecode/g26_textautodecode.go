// g26_textautodecode.go
// Learning go, System programming, files, Detect ancoding of a text file
// Similar to Rust crate TextAutoDecode
//
// 2025-06-23	PV		First version

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

func main() {
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-empty.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-ascii.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf8bom.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16lebom.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16bebom.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf8.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16le.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16be.txt`)
	test(`C:\DocumentsOD\Doc tech\Encodings\prenoms-1252.txt`)
}

func test(filename string) {
	tad, err := ReadTextFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", filename, err)
		os.Exit(1)
	}

	fmt.Printf("%-65.65s %s\n", filename, strEncoding(tad.Encoding))
	if strings.Contains(filename, "empty") {
		return
	}

	var beginning string
	if strings.Contains(filename, "ascii") {
		beginning = "juliette sophie brigitte geraldine"
	} else {
		beginning = "juliette sophie brigitte géraldine"
	}

	if !strings.HasPrefix(tad.Text, beginning) {
		l := min(len(tad.Text), 80)
		fmt.Println("Wrong prefix:", "«"+tad.Text[:l]+"»")
	}
}

const MILLE = 1000

type TextFileEncoding int

const (
	NotText    TextFileEncoding = iota // Binary or unrecognized text (for instance contains chars in 0..31 other than \r \n \t)
	Empty                              // File is empty
	ASCII                              // Only 7-bit characters
	EightBit                           // ANSI/Windows 1252 (only this variant is checked)
	UTF8                               // Plain UTF-8 without BOM
	UTF8BOM                            // Starts with EF BB BF
	UTF16LE                            // No BOM but UTF-16 LE detected
	UTF16BE                            // No BOM but UTF-16 BE detected
	UTF16LEBOM                         // Starts with FF FE (Windows)
	UTF16BEBOM                         // Starts with FE FF
	InProgress
	FileError
)

type TextAutoDecode struct {
	Text     string
	Encoding TextFileEncoding
}

func strEncoding(enc TextFileEncoding) string {
	switch enc {
	case NotText:
		return "NotText"
	case Empty:
		return "Empty"
	case ASCII:
		return "ASCII"
	case EightBit:
		return "EightBit"
	case UTF8:
		return "UTF8"
	case UTF8BOM:
		return "UTF8BOM"
	case UTF16LE:
		return "UTF16LE"
	case UTF16BE:
		return "UTF16BE"
	case UTF16LEBOM:
		return "UTF16LEBOM"
	case UTF16BEBOM:
		return "UTF16BEBOM"
	case InProgress:
		return "InProgress"
	case FileError:
		return "FileError"
	default:
		panic("Error")
	}
}

func ReadTextFile(file string) (TextAutoDecode, error) {
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

	buffer_full := make([]byte, 0)
	is_buffer_full_read := false

	// First we check presence of BOM. If present, then file type is determined,
	// if further checks fail, no need to continue.
	// We don't consider the case of a EightBit (Windows 1252) file that would begin with these three BOM,
	// it's *really* unlikely, now that EightBit files are getting rare.

	// UTF-8 BOM?
	// Since we have a BOM, no need to check for ASCII subset
	if n >= 3 && buffer_1000[0] == 0xEF && buffer_1000[1] == 0xBB && buffer_1000[2] == 0xBF {
		s, ok := check_utf8(buffer_1000, n)

		if !ok {
			return TextAutoDecode{Text: "", Encoding: NotText}, nil
		}

		if n < MILLE {
			return TextAutoDecode{Text: s[3:], Encoding: UTF8BOM}, nil
		}

		return final_read(&is_buffer_full_read, &buffer_full, f, UTF8BOM)
	}

	// UTF-16 LE BOM? (Windows)
	if n >= 2 && buffer_1000[0] == 0xFF && buffer_1000[1] == 0xFE {
		s, ok := check_utf16(buffer_1000, n, UTF16LEBOM)
		if !ok {
			return TextAutoDecode{Text: "", Encoding: NotText}, nil
		}

		if n < MILLE {
			return TextAutoDecode{Text: s, Encoding: UTF16LEBOM}, nil
		}

		return final_read(&is_buffer_full_read, &buffer_full, f, UTF16LEBOM)
	}

	// UTF-16 BE BOM?
	if n >= 2 && buffer_1000[0] == 0xFE && buffer_1000[1] == 0xFF {
		s, ok := check_utf16(buffer_1000, n, UTF16BEBOM)
		if !ok {
			return TextAutoDecode{Text: "", Encoding: NotText}, nil
		}

		if n < MILLE {
			return TextAutoDecode{Text: s, Encoding: UTF16BEBOM}, nil
		}

		return final_read(&is_buffer_full_read, &buffer_full, f, UTF16BEBOM)
	}

	// Then check encodings without BOM

	// UTF-8 without BOM?
	// Note that if string is only ASCII text, then type is assumed ASCII instead of UTF-8
	s, ok := check_utf8(buffer_1000, n)
	if ok {
		if n < 1000 {
			var e TextFileEncoding
			if is_ascii_text(&s) {
				e = ASCII
			} else {
				e = UTF8
			}
			return TextAutoDecode{Text: s, Encoding: e}, nil
		} else {
			// Special case, first 1000 bytes are ASCII so we got there, but after 1000 bytes, we get 8-bit
			// characters so we can't return if we didn't recognize the whole file as UTF-8
			tad, err := final_read(&is_buffer_full_read, &buffer_full, f, UTF8)
			if err == nil {
				if tad.Encoding != NotText {
					return tad, err
				}
			}
		}

		// We skip checking UTF-16, since it's a match for UTF-8/ASCII on the furst 1000 chars
		return final_read(&is_buffer_full_read, &buffer_full, f, EightBit)
	}

	// UTF-16 LE? (Windows)
	// Only files with more than 10 characters (20 bytes) are tested and checked for 75% ASCII, or many small binary non text-files will match
	if n > 20 {
		s, ok := check_utf16(buffer_1000, n, UTF16LE)
		if ok {
			if n < 1000 {
				return TextAutoDecode{Text: s, Encoding: UTF16LE}, nil
			}

			return final_read(
				&is_buffer_full_read,
				&buffer_full,
				f,
				UTF16LE)
		}

		// UTF-16 BE?
		s, ok = check_utf16(buffer_1000, n, UTF16BE)
		if ok {
			if n < 1000 {
				return TextAutoDecode{Text: s, Encoding: UTF16BE}, nil
			}
			return final_read(
				&is_buffer_full_read,
				&buffer_full,
				f,
				UTF16BE)
		}
	}

	// 8-bit?
	s, ok = check_eightbit(&buffer_1000, n)
	if ok {
		if n < 1000 {
			return TextAutoDecode{Text: s, Encoding: EightBit}, nil
		} else {
			return final_read(
				&is_buffer_full_read,
				&buffer_full,
				f,
				EightBit)
		}
	}

	// None of the encodings worked without error

	return TextAutoDecode{Text: "??", Encoding: InProgress}, nil
}

// The 75% ASCII test is too restrictive, some valid UTF-8 files are rejected (ex: output of tree command)
// So we only detect control characters that should not be present in a text file
// Old text files may contain FF (Form Feed, 12) or VT (Vertical Tab, 11), but it's unlikely for common files
func contains_binary_chars(s *string, also_check_block_c1 bool) bool {
	for _, c := range *s {
		if c < 32 && (c != 9 && c != 10 && c != 13) {
			return true
		}
		// If requested, no characters of C1 is accepted (for all encodings but 8-bit)
		if also_check_block_c1 && c >= 128 && c < 160 {
			return true
		}
	}
	return false
}

// Check that string s doesn't contain a null char and contains at least 75% of ASCII 32..127, CR, LF, TAB
func is_75percent_ascii(s *string) bool {
	acount := 0
	l := len(*s)
	for _, c := range *s {
		// For 8-bit files, we only exclude non-comon elements of C0 block, and DEL (127) char
		// Anything in [128..255] is accepted
		if c == 127 || c < 32 && (c != 9 && c != 10 && c != 13) {
			return false
		}
		if c >= 32 && c < 127 || c == 9 || c == 10 || c == 13 {
			acount += 1
		}
	}
	if l < 10 {
		return true
	} else {
		return float64(acount)/float64(l) >= 0.75
	}
}

func final_read(is_buffer_full_read *bool, buffer_full *[]byte, file *os.File, encoding TextFileEncoding) (TextAutoDecode, error) {
	// If the whole file has not been read yet, then read it
	if !*is_buffer_full_read {
		// Rewind the file position to the beginning
		_, _ = file.Seek(0, io.SeekStart)
		temp_buffer_full, err := io.ReadAll(file)
		if err != nil {
			return TextAutoDecode{}, err
		}
		*buffer_full = temp_buffer_full
		*is_buffer_full_read = true
	}

	text := ""
	switch encoding {
	case UTF8, UTF8BOM:
		if utf8.Valid(*buffer_full) {
			if encoding == UTF8 {
				text = string(*buffer_full)
			} else {
				text = string(*buffer_full)[3:]
			}
		} else {
			return TextAutoDecode{Text: "", Encoding: NotText}, nil
		}

	case UTF16LE, UTF16LEBOM, UTF16BE, UTF16BEBOM:
		s, ok := utf16_decode(*buffer_full, encoding)
		if !ok {
			return TextAutoDecode{Text: "", Encoding: NotText}, nil
		}
		text = s

	case EightBit:
		s, ok := eightbit_decode(*buffer_full)
		if !ok {
			return TextAutoDecode{Text: "", Encoding: NotText}, nil
		}
		text = s

	default:
		panic("final_read: encoding not supported yet!")
	}

	check_ascii := encoding == UTF8 // UTF8_BOM is never considered ASCII

	// Without BOM, we add heuristics to be sure that what has been decoded makes sense
	check_75percent_text := encoding == EightBit || encoding == UTF16BE || encoding == UTF16LE

	// Special heuristics to be sure it's a valid text files
	if check_75percent_text && !is_75percent_ascii(&text) {
		return TextAutoDecode{Text: "", Encoding: NotText}, nil
	}
	if encoding != EightBit && contains_binary_chars(&text, true) {
		return TextAutoDecode{Text: "", Encoding: NotText}, nil
	}

	e := encoding
	if check_ascii {
		if is_ascii_text(&text) {
			e = ASCII
		} else {
			e = UTF8
		}
	}

	return TextAutoDecode{Text: text, Encoding: e}, nil
}

// check_utf8 checks if a small byte buffer of n bytes (max 1000) contains a valid UTF-8 string.
// If valid, it returns the string and true.
// If not valid, it returns an empty string and false.
// Note that if buffer len is exactly 1000, it's possible that the last UTF-8 character is truncated,
// so we reduce the buffer to be safe. Anyway, in this case, we'll reread the whole file and do a global check
func check_utf8(buffer_1000 []byte, n int) (string, bool) {
	if buffer_1000 == nil || n < 0 {
		panic("Internal error")
	}

	nsafe := n
	if n == MILLE {
		for nsafe = MILLE - 1; ; nsafe-- {
			// If buffer[nsafe] is a valid beginning for UTF-8 encoding, we can stop here
			if buffer_1000[nsafe] < 128 || (buffer_1000[nsafe]&0b11100000) == 0b11000000 || (buffer_1000[nsafe]&0b11110000) == 0b11100000 || (buffer_1000[nsafe]&0b11111000) == 0b11110000 {
				break
			}
			// If it's a continuation character, we can continue, but at most three continuation characters are valid
			if (buffer_1000[nsafe]&0b11000000) == 0b10000000 && n >= MILLE-3 {
				continue
			}
			// Sorry, that's not valid UTF-8...
			return "", false
		}
	}

	// If last character is <128, it's not truncated and can be kept
	if buffer_1000[nsafe] < 128 {
		nsafe++
	}
	buffer_safe := buffer_1000[:nsafe]

	// Use utf8.Valid to check if the byte slice is valid UTF-8
	if utf8.Valid(buffer_safe) {
		s := string(buffer_safe)
		if contains_binary_chars(&s, true) {
			return "", false
		}

		return string(buffer_safe), true
	}

	// If not valid, return an empty string and false
	return "", false
}

func check_utf16(buffer_1000 []byte, n int, encoding TextFileEncoding) (string, bool) {
	if buffer_1000 == nil || n < 0 {
		panic("Internal error")
	}

	// We have to check whether we truncated reading in the middle of a surrogate sequence when reading 1000 bytes max.
	// Lead surrogate is 0xD800-0xDBFF (and tail surrogate is 0xDC00-0xDFFF), if the byte at index 998 is 0xD8, then
	// we cut a surrogate. Note that optional byte order header (0xFF, 0xFE) is two bytes long, so all UTF-16 words
	// start at even index.
	nsafe := n

	if n == MILLE {
		off := 0
		if encoding == UTF16BE || encoding == UTF16BEBOM {
			off = 1
		}

		pa := 998
		if buffer_1000[pa+off] >= 0xD8 && buffer_1000[pa+off] <= 0xDB {
			pa -= 2
		}
		nsafe = pa + 2
	}
	buffer_safe := buffer_1000[:nsafe]

	s, ok := utf16_decode(buffer_safe, encoding)
	if !ok {
		return s, ok
	}

	// If there is no BOM, actually UTF-16 BE can be decoded as UTF-16 LE and also the reverse in most of cases.
	// To be sure there is no confusion, add an extra heuristics to check that content is 75% ASCII
	if (encoding == UTF16LE || encoding == UTF16BE) && !is_75percent_ascii(&s) {
		return "", false
	}

	if !contains_binary_chars(&s, true) {
		return s, ok
	}
	return "", false
}

func utf16_decode(buffer []byte, encoding TextFileEncoding) (string, bool) {
	// Buffer len must be even for UTF-16
	if len(buffer)&1 == 1 {
		return "", false
	}

	if len(buffer) == 0 {
		return "", encoding == UTF16LE || encoding == UTF16BE
	}

	off := 0
	if encoding == UTF16BE || encoding == UTF16BEBOM {
		off = 1
	}

	start := 0
	if encoding == UTF16BEBOM || encoding == UTF16LEBOM {
		// Check BOM
		if buffer[off] != 0xFF || buffer[1-off] != 0xFE {
			return "", false
		}
		start = 2
	}

	const (
		// 0xd800-0xdc00 encodes the high 10 bits of a pair.
		// 0xdc00-0xe000 encodes the low 10 bits of a pair.
		// the value is those 20 bits plus 0x10000.
		surr1 = 0xd800
		surr2 = 0xdc00
		surr3 = 0xe000

		surrSelf = 0x10000
	)

	var buf []rune
	for start < len(buffer) {
		r := (int(buffer[start+off]) + (int(buffer[start+1-off]) << 8))
		var ar rune
		switch {
		case r < surr1, surr3 <= r:
			// normal rune
			ar = rune(r)
		case r >= surr1 && r < surr2: // High surrogate
			if start+2 >= len(buffer) {
				return "", false
			}
			start += 2
			r2 := (int(buffer[start+off]) + (int(buffer[start+1-off]) << 8))
			ar = rune((r-surr1)<<10 | (r2 - surr2) + surrSelf)
		default:
			return "", false
		}
		buf = append(buf, ar)
		start += 2
	}
	return string(buf), true
}

func check_eightbit(buffer_1000 *[]byte, _ int) (string, bool) {
	s, ok := eightbit_decode(*buffer_1000)
	if ok && is_75percent_ascii(&s) {
		return s, ok
	}

	return "", false
}

func eightbit_decode(buffer []byte) (string, bool) {
	// Create a new decoder for Windows CP 1252
	decoder := charmap.Windows1252.NewDecoder()

	// Use ioutil.ReadAll with the decoder to convert the byte slice
	// transform.NewReader creates a new reader that decodes the input
	utf8Bytes, err := io.ReadAll(decoder.Reader(bytes.NewReader(buffer)))
	if err != nil {
		return "", false
	}

	// Convert the UTF-8 byte slice to a Go string
	return string(utf8Bytes), true
}

func is_ascii_text(s *string) bool {
	for _, b := range *s {
		if b > 126 || (b < 32 && b != '\r' && b != '\n' && b != '\t') {
			return false
		}
	}
	return true
}

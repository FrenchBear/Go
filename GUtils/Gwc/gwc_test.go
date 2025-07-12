// Tests for Gwc
//
// 2025-07-10 	PV 		First version

package main

import (
	"testing"
)

func assert_eq(t *testing.T, a, b int) {
	if a != b {
		t.Errorf("Expected %d, got %d", b, a)
	}
}


func TestCount1(t *testing.T) {
	o := Options{ShowOnlyTotal: true }
	b := DataBag{}
    processText(&b, "Once upon a time\nWas a King and a Prince\nIn a far, far away kingdom.", "(test)", &o, 68)
    assert_eq(t, b.files_count, 1)
    assert_eq(t, b.lines_count, 3)
    assert_eq(t, b.words_count, 16)
    assert_eq(t, b.chars_count, 68)
    assert_eq(t, b.bytes_count, 68)
}

func TestCount2(t *testing.T) {
	o := Options{ShowOnlyTotal: true }
	b := DataBag{}
    processText(&b, " Aé♫山𝄞🐗   🐷🐽🐖 ", "(test)", &o, 34)
    assert_eq(t, b.files_count, 1)
    assert_eq(t, b.lines_count, 1)
    assert_eq(t, b.words_count, 2)
    assert_eq(t, b.chars_count, 14)
    assert_eq(t, b.bytes_count, 34)
}

func TestFileAscii(t *testing.T) {
	o := Options{ShowOnlyTotal: true }
	b := DataBag{}
    processFile(&b, `C:\DocumentsOD\Doc tech\Encodings\prenoms-ascii.txt`, &o)
    assert_eq(t, b.files_count, 1)
    assert_eq(t, b.lines_count, 9)
    assert_eq(t, b.words_count, 143)
    assert_eq(t, b.chars_count, 1145)
    assert_eq(t, b.bytes_count, 1145)
}

func TestFileRtf8(t *testing.T) {
	o := Options{ShowOnlyTotal: true }
	b := DataBag{}
    processFile(&b, `C:\DocumentsOD\Doc tech\Encodings\prenoms-utf8.txt`, &o)
    assert_eq(t, b.files_count, 1)
    assert_eq(t, b.lines_count, 9)
    assert_eq(t, b.words_count, 143)
    assert_eq(t, b.chars_count, 1145)
    assert_eq(t, b.bytes_count, 1194)
}

func TestFileUtf16lebom(t *testing.T) {
	o := Options{ShowOnlyTotal: true }
	b := DataBag{}
    processFile(&b, `C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16lebom.txt`, &o)
    assert_eq(t, b.files_count, 1)
    assert_eq(t, b.lines_count, 9)
    assert_eq(t, b.words_count, 143)
    assert_eq(t, b.chars_count, 1145)
    assert_eq(t, b.bytes_count, 2292)
}

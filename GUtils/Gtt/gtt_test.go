// Tests for Gtt
//
// 2025-07-05 	PV 		Translation of Rust equivalent by Gemini

package main

import (
	"os"
	"testing"
)

func TestEmpty(t *testing.T) {
	tempFile, err := os.CreateTemp("", "rtt-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.Write([]byte{})
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Sync()

	b := NewDataBag()
	res := processFile(b, tempFile.Name(), "(test empty)")

	if res != "(test empty): «Empty file»" {
		t.Errorf("Expected \"(test empty): «Empty file»\", got \"%s\"", res)
	}

	if b.FilesTypes.Total != 1 {
		t.Errorf("Expected FilesTypes.Total 1, got %d", b.FilesTypes.Total)
	}
	if b.FilesTypes.Empty != 1 {
		t.Errorf("Expected FilesTypes.Empty 1, got %d", b.FilesTypes.Empty)
	}
	if b.FilesTypes.Ascii != 0 {
		t.Errorf("Expected FilesTypes.Ascii 0, got %d", b.FilesTypes.Ascii)
	}
	if b.FilesTypes.Utf8 != 0 {
		t.Errorf("Expected FilesTypes.Utf8 0, got %d", b.FilesTypes.Utf8)
	}
	if b.FilesTypes.Utf16 != 0 {
		t.Errorf("Expected FilesTypes.Utf16 0, got %d", b.FilesTypes.Utf16)
	}
	if b.FilesTypes.EightBit != 0 {
		t.Errorf("Expected FilesTypes.EightBit 0, got %d", b.FilesTypes.EightBit)
	}
	if b.FilesTypes.NonText != 0 {
		t.Errorf("Expected FilesTypes.NonText 0, got %d", b.FilesTypes.NonText)
	}

	if b.EolStyles.Total != 0 {
		t.Errorf("Expected EolStyles.Total 0, got %d", b.EolStyles.Total)
	}
	if b.EolStyles.Windows != 0 {
		t.Errorf("Expected EolStyles.Windows 0, got %d", b.EolStyles.Windows)
	}
	if b.EolStyles.Unix != 0 {
		t.Errorf("Expected EolStyles.Unix 0, got %d", b.EolStyles.Unix)
	}
	if b.EolStyles.Mac != 0 {
		t.Errorf("Expected EolStyles.Mac 0, got %d", b.EolStyles.Mac)
	}
	if b.EolStyles.Mixed != 0 {
		t.Errorf("Expected EolStyles.Mixed 0, got %d", b.EolStyles.Mixed)
	}
}

func TestAscii(t *testing.T) {
	tempFile, err := os.CreateTemp("", "rtt-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.Write([]byte{'H', 'e', 'l', 'l', 'o', '\r', '\n'})
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Sync()

	b := NewDataBag()
	res := processFile(b, tempFile.Name(), "(test ascii)")

	if res != "(test ascii): ASCII, Windows" {
		t.Errorf("Expected \"(test ascii): ASCII, Windows\", got \"%s\"", res)
	}

	if b.FilesTypes.Total != 1 {
		t.Errorf("Expected FilesTypes.Total 1, got %d", b.FilesTypes.Total)
	}
	if b.FilesTypes.Empty != 0 {
		t.Errorf("Expected FilesTypes.Empty 0, got %d", b.FilesTypes.Empty)
	}
	if b.FilesTypes.Ascii != 1 {
		t.Errorf("Expected FilesTypes.Ascii 1, got %d", b.FilesTypes.Ascii)
	}
	if b.FilesTypes.Utf8 != 0 {
		t.Errorf("Expected FilesTypes.Utf8 0, got %d", b.FilesTypes.Utf8)
	}
	if b.FilesTypes.Utf16 != 0 {
		t.Errorf("Expected FilesTypes.Utf16 0, got %d", b.FilesTypes.Utf16)
	}
	if b.FilesTypes.EightBit != 0 {
		t.Errorf("Expected FilesTypes.EightBit 0, got %d", b.FilesTypes.EightBit)
	}
	if b.FilesTypes.NonText != 0 {
		t.Errorf("Expected FilesTypes.NonText 0, got %d", b.FilesTypes.NonText)
	}

	if b.EolStyles.Total != 1 {
		t.Errorf("Expected EolStyles.Total 1, got %d", b.EolStyles.Total)
	}
	if b.EolStyles.Windows != 1 {
		t.Errorf("Expected EolStyles.Windows 1, got %d", b.EolStyles.Windows)
	}
	if b.EolStyles.Unix != 0 {
		t.Errorf("Expected EolStyles.Unix 0, got %d", b.EolStyles.Unix)
	}
	if b.EolStyles.Mac != 0 {
		t.Errorf("Expected EolStyles.Mac 0, got %d", b.EolStyles.Mac)
	}
	if b.EolStyles.Mixed != 0 {
		t.Errorf("Expected EolStyles.Mixed 0, got %d", b.EolStyles.Mixed)
	}
}

func TestNonText1(t *testing.T) {
	tempFile, err := os.CreateTemp("", "rtt-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.Write([]byte{0xCA, 0xFE, 0xDE, 0xAD, 0xBE, 0xEF})
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Sync()

	b := NewDataBag()
	res := processFile(b, tempFile.Name(), "(test non-text)")

	if res != "" {
		t.Errorf("Expected \"\", got \"%s\"", res)
	}

	if b.FilesTypes.Total != 1 {
		t.Errorf("Expected FilesTypes.Total 1, got %d", b.FilesTypes.Total)
	}
	if b.FilesTypes.Empty != 0 {
		t.Errorf("Expected FilesTypes.Empty 0, got %d", b.FilesTypes.Empty)
	}
	if b.FilesTypes.Ascii != 0 {
		t.Errorf("Expected FilesTypes.Ascii 0, got %d", b.FilesTypes.Ascii)
	}
	if b.FilesTypes.Utf8 != 0 {
		t.Errorf("Expected FilesTypes.Utf8 0, got %d", b.FilesTypes.Utf8)
	}
	if b.FilesTypes.Utf16 != 0 {
		t.Errorf("Expected FilesTypes.Utf16 0, got %d", b.FilesTypes.Utf16)
	}
	if b.FilesTypes.EightBit != 0 {
		t.Errorf("Expected FilesTypes.EightBit 0, got %d", b.FilesTypes.EightBit)
	}
	if b.FilesTypes.NonText != 1 {
		t.Errorf("Expected FilesTypes.NonText 1, got %d", b.FilesTypes.NonText)
	}

	if b.EolStyles.Total != 0 {
		t.Errorf("Expected EolStyles.Total 0, got %d", b.EolStyles.Total)
	}
	if b.EolStyles.Windows != 0 {
		t.Errorf("Expected EolStyles.Windows 0, got %d", b.EolStyles.Windows)
	}
	if b.EolStyles.Unix != 0 {
		t.Errorf("Expected EolStyles.Unix 0, got %d", b.EolStyles.Unix)
	}
	if b.EolStyles.Mac != 0 {
		t.Errorf("Expected EolStyles.Mac 0, got %d", b.EolStyles.Mac)
	}
	if b.EolStyles.Mixed != 0 {
		t.Errorf("Expected EolStyles.Mixed 0, got %d", b.EolStyles.Mixed)
	}
}

func TestNonText2(t *testing.T) {
	tempFile, err := os.CreateTemp("", "rtt-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.Write([]byte{0xCA, 0xFE, 0xDE, 0xAD, 0xBE, 0xEF})
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Sync()

	b := NewDataBag()
	res := processFile(b, tempFile.Name(), "non-text.rs")

	if res != "non-text.rs: «Non-text file detected, but extension rs is usually a text file»" {
		t.Errorf("Expected \"non-text.rs: «Non-text file detected, but extension rs is usually a text file»\", got \"%s\"", res)
	}

	if b.FilesTypes.Total != 1 {
		t.Errorf("Expected FilesTypes.Total 1, got %d", b.FilesTypes.Total)
	}
	if b.FilesTypes.Empty != 0 {
		t.Errorf("Expected FilesTypes.Empty 0, got %d", b.FilesTypes.Empty)
	}
	if b.FilesTypes.Ascii != 0 {
		t.Errorf("Expected FilesTypes.Ascii 0, got %d", b.FilesTypes.Ascii)
	}
	if b.FilesTypes.Utf8 != 0 {
		t.Errorf("Expected FilesTypes.Utf8 0, got %d", b.FilesTypes.Utf8)
	}
	if b.FilesTypes.Utf16 != 0 {
		t.Errorf("Expected FilesTypes.Utf16 0, got %d", b.FilesTypes.Utf16)
	}
	if b.FilesTypes.EightBit != 0 {
		t.Errorf("Expected FilesTypes.EightBit 0, got %d", b.FilesTypes.EightBit)
	}
	if b.FilesTypes.NonText != 1 {
		t.Errorf("Expected FilesTypes.NonText 1, got %d", b.FilesTypes.NonText)
	}

	if b.EolStyles.Total != 0 {
		t.Errorf("Expected EolStyles.Total 0, got %d", b.EolStyles.Total)
	}
	if b.EolStyles.Windows != 0 {
		t.Errorf("Expected EolStyles.Windows 0, got %d", b.EolStyles.Windows)
	}
	if b.EolStyles.Unix != 0 {
		t.Errorf("Expected EolStyles.Unix 0, got %d", b.EolStyles.Unix)
	}
	if b.EolStyles.Mac != 0 {
		t.Errorf("Expected EolStyles.Mac 0, got %d", b.EolStyles.Mac)
	}
	if b.EolStyles.Mixed != 0 {
		t.Errorf("Expected EolStyles.Mixed 0, got %d", b.EolStyles.Mixed)
	}
}

func TestUtf8(t *testing.T) {
	model := []byte{
		0x41,                   // A
		0xC3, 0xA9,             // é
		0xE2, 0x99, 0xAB,       // ♫
		0xE5, 0xB1, 0xB1,       // 山
		0xF0, 0x9D, 0x84, 0x9E, // 𝄞
		0xF0, 0x9F, 0x90, 0x97, // 🐗
	}

	tempFile, err := os.CreateTemp("", "rtt-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.Write(model)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Sync()

	b := NewDataBag()
	res := processFile(b, tempFile.Name(), "(test utf8)")

	if res != "(test utf8): UTF-8, No EOL detected" {
		t.Errorf("Expected \"(test utf8): UTF-8, No EOL detected\", got \"%s\"", res)
	}

	if b.FilesTypes.Total != 1 {
		t.Errorf("Expected FilesTypes.Total 1, got %d", b.FilesTypes.Total)
	}
	if b.FilesTypes.Empty != 0 {
		t.Errorf("Expected FilesTypes.Empty 0, got %d", b.FilesTypes.Empty)
	}
	if b.FilesTypes.Ascii != 0 {
		t.Errorf("Expected FilesTypes.Ascii 0, got %d", b.FilesTypes.Ascii)
	}
	if b.FilesTypes.Utf8 != 1 {
		t.Errorf("Expected FilesTypes.Utf8 1, got %d", b.FilesTypes.Utf8)
	}
	if b.FilesTypes.Utf16 != 0 {
		t.Errorf("Expected FilesTypes.Utf16 0, got %d", b.FilesTypes.Utf16)
	}
	if b.FilesTypes.EightBit != 0 {
		t.Errorf("Expected FilesTypes.EightBit 0, got %d", b.FilesTypes.EightBit)
	}
	if b.FilesTypes.NonText != 0 {
		t.Errorf("Expected FilesTypes.NonText 0, got %d", b.FilesTypes.NonText)
	}

	if b.EolStyles.Total != 0 {
		t.Errorf("Expected EolStyles.Total 0, got %d", b.EolStyles.Total)
	}
	if b.EolStyles.Windows != 0 {
		t.Errorf("Expected EolStyles.Windows 0, got %d", b.EolStyles.Windows)
	}
	if b.EolStyles.Unix != 0 {
		t.Errorf("Expected EolStyles.Unix 0, got %d", b.EolStyles.Unix)
	}
	if b.EolStyles.Mac != 0 {
		t.Errorf("Expected EolStyles.Mac 0, got %d", b.EolStyles.Mac)
	}
	if b.EolStyles.Mixed != 0 {
		t.Errorf("Expected EolStyles.Mixed 0, got %d", b.EolStyles.Mixed)
	}
}

func TestUtf8Bom(t *testing.T) {
	model := []byte{
		0xEF, 0xBB, 0xBF, // UTF-8 BOM
		0x41,       // A
		'\r',       // Mac EOL
	}

	tempFile, err := os.CreateTemp("", "rtt-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.Write(model)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Sync()

	b := NewDataBag()
	res := processFile(b, tempFile.Name(), "(test utf8bom)")

	if res != "(test utf8bom): UTF-8 «with BOM», Mac" {
		t.Errorf("Expected \"(test utf8bom): UTF-8 «with BOM», Mac\", got \"%s\"", res)
	}

	if b.FilesTypes.Total != 1 {
		t.Errorf("Expected FilesTypes.Total 1, got %d", b.FilesTypes.Total)
	}
	if b.FilesTypes.Empty != 0 {
		t.Errorf("Expected FilesTypes.Empty 0, got %d", b.FilesTypes.Empty)
	}
	if b.FilesTypes.Ascii != 0 {
		t.Errorf("Expected FilesTypes.Ascii 0, got %d", b.FilesTypes.Ascii)
	}
	if b.FilesTypes.Utf8 != 1 {
		t.Errorf("Expected FilesTypes.Utf8 1, got %d", b.FilesTypes.Utf8)
	}
	if b.FilesTypes.Utf16 != 0 {
		t.Errorf("Expected FilesTypes.Utf16 0, got %d", b.FilesTypes.Utf16)
	}
	if b.FilesTypes.EightBit != 0 {
		t.Errorf("Expected FilesTypes.EightBit 0, got %d", b.FilesTypes.EightBit)
	}
	if b.FilesTypes.NonText != 0 {
		t.Errorf("Expected FilesTypes.NonText 0, got %d", b.FilesTypes.NonText)
	}

	if b.EolStyles.Total != 1 {
		t.Errorf("Expected EolStyles.Total 1, got %d", b.EolStyles.Total)
	}
	if b.EolStyles.Windows != 0 {
		t.Errorf("Expected EolStyles.Windows 0, got %d", b.EolStyles.Windows)
	}
	if b.EolStyles.Unix != 0 {
		t.Errorf("Expected EolStyles.Unix 0, got %d", b.EolStyles.Unix)
	}
	if b.EolStyles.Mac != 1 {
		t.Errorf("Expected EolStyles.Mac 1, got %d", b.EolStyles.Mac)
	}
	if b.EolStyles.Mixed != 0 {
		t.Errorf("Expected EolStyles.Mixed 0, got %d", b.EolStyles.Mixed)
	}
}

func TestUtf16LeBom(t *testing.T) {
	model := []byte{
		0xFF, 0xFE, // UTF-16 LE BOM
		0x41, 0x00, // A
		0x42, 0x00, // B
		'\n', 0x00, // Unix EOL
		0x43, 0x00, // C
		0x44, 0x00, // D
		'\r', 0x00, '\n', 0x00, // Windows EOL
	}

	tempFile, err := os.CreateTemp("", "rtt-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.Write(model)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Sync()

	b := NewDataBag()
	res := processFile(b, tempFile.Name(), "(test utf16lebom)")

	if res != "(test utf16lebom): UTF-16 LE, «Mixed EOL styles»" {
		t.Errorf("Expected \"(test utf16lebom): UTF-16 LE, «Mixed EOL styles»\", got \"%s\"", res)
	}

	if b.FilesTypes.Total != 1 {
		t.Errorf("Expected FilesTypes.Total 1, got %d", b.FilesTypes.Total)
	}
	if b.FilesTypes.Empty != 0 {
		t.Errorf("Expected FilesTypes.Empty 0, got %d", b.FilesTypes.Empty)
	}
	if b.FilesTypes.Ascii != 0 {
		t.Errorf("Expected FilesTypes.Ascii 0, got %d", b.FilesTypes.Ascii)
	}
	if b.FilesTypes.Utf8 != 0 {
		t.Errorf("Expected FilesTypes.Utf8 0, got %d", b.FilesTypes.Utf8)
	}
	if b.FilesTypes.Utf16 != 1 {
		t.Errorf("Expected FilesTypes.Utf16 1, got %d", b.FilesTypes.Utf16)
	}
	if b.FilesTypes.EightBit != 0 {
		t.Errorf("Expected FilesTypes.EightBit 0, got %d", b.FilesTypes.EightBit)
	}
	if b.FilesTypes.NonText != 0 {
		t.Errorf("Expected FilesTypes.NonText 0, got %d", b.FilesTypes.NonText)
	}

	if b.EolStyles.Total != 1 {
		t.Errorf("Expected EolStyles.Total 1, got %d", b.EolStyles.Total)
	}
	if b.EolStyles.Windows != 1 {
		t.Errorf("Expected EolStyles.Windows 1, got %d", b.EolStyles.Windows)
	}
	if b.EolStyles.Unix != 1 {
		t.Errorf("Expected EolStyles.Unix 1, got %d", b.EolStyles.Unix)
	}
	if b.EolStyles.Mac != 0 {
		t.Errorf("Expected EolStyles.Mac 0, got %d", b.EolStyles.Mac)
	}
	if b.EolStyles.Mixed != 1 {
		t.Errorf("Expected EolStyles.Mixed 1, got %d", b.EolStyles.Mixed)
	}
}

func TestUtf16Le1(t *testing.T) {
	model := []byte{
		0x41, 0x00, // A
		0x42, 0x00, // B
		0x43, 0x00, // C
		0x44, 0x00, // D
		0x45, 0x00, // E
		'\n', 0x00, // Unix EOL
		0x61, 0x00, // a
		0x62, 0x00, // b
		0x63, 0x00, // c
		0x64, 0x00, // d
		0x65, 0x00, // e
		'\n', 0x00, // Unix EOL
	}

	tempFile, err := os.CreateTemp("", "rtt-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.Write(model)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Sync()

	b := NewDataBag()
	res := processFile(b, tempFile.Name(), "(test utf16le1)")

	if res != "(test utf16le1): UTF-16 LE «without BOM», Unix" {
		t.Errorf("Expected \"(test utf16le1): UTF-16 LE «without BOM», Unix\", got \"%s\"", res)
	}

	if b.FilesTypes.Total != 1 {
		t.Errorf("Expected FilesTypes.Total 1, got %d", b.FilesTypes.Total)
	}
	if b.FilesTypes.Empty != 0 {
		t.Errorf("Expected FilesTypes.Empty 0, got %d", b.FilesTypes.Empty)
	}
	if b.FilesTypes.Ascii != 0 {
		t.Errorf("Expected FilesTypes.Ascii 0, got %d", b.FilesTypes.Ascii)
	}
	if b.FilesTypes.Utf8 != 0 {
		t.Errorf("Expected FilesTypes.Utf8 0, got %d", b.FilesTypes.Utf8)
	}
	if b.FilesTypes.Utf16 != 1 {
		t.Errorf("Expected FilesTypes.Utf16 1, got %d", b.FilesTypes.Utf16)
	}
	if b.FilesTypes.EightBit != 0 {
		t.Errorf("Expected FilesTypes.EightBit 0, got %d", b.FilesTypes.EightBit)
	}
	if b.FilesTypes.NonText != 0 {
		t.Errorf("Expected FilesTypes.NonText 0, got %d", b.FilesTypes.NonText)
	}

	if b.EolStyles.Total != 1 {
		t.Errorf("Expected EolStyles.Total 1, got %d", b.EolStyles.Total)
	}
	if b.EolStyles.Windows != 0 {
		t.Errorf("Expected EolStyles.Windows 0, got %d", b.EolStyles.Windows)
	}
	if b.EolStyles.Unix != 1 {
		t.Errorf("Expected EolStyles.Unix 1, got %d", b.EolStyles.Unix)
	}
	if b.EolStyles.Mac != 0 {
		t.Errorf("Expected EolStyles.Mac 0, got %d", b.EolStyles.Mac)
	}
	if b.EolStyles.Mixed != 0 {
		t.Errorf("Expected EolStyles.Mixed 0, got %d", b.EolStyles.Mixed)
	}
}

func TestUtf16Le2(t *testing.T) {
	model := []byte{
		0x41, 0x00, // A
		'\n', 0x00, // Unix EOL
		0x61, 0x00, // a
		'\n', 0x00, // Unix EOL
	}

	tempFile, err := os.CreateTemp("", "rtt-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = tempFile.Write(model)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Sync()

	b := NewDataBag()
	res := processFile(b, tempFile.Name(), "(test utf16le1)")

	if res != "" {
		t.Errorf("Expected \"\", got \"%s\"", res)
	}

	if b.FilesTypes.Total != 1 {
		t.Errorf("Expected FilesTypes.Total 1, got %d", b.FilesTypes.Total)
	}
	if b.FilesTypes.Empty != 0 {
		t.Errorf("Expected FilesTypes.Empty 0, got %d", b.FilesTypes.Empty)
	}
	if b.FilesTypes.Ascii != 0 {
		t.Errorf("Expected FilesTypes.Ascii 0, got %d", b.FilesTypes.Ascii)
	}
	if b.FilesTypes.Utf8 != 0 {
		t.Errorf("Expected FilesTypes.Utf8 0, got %d", b.FilesTypes.Utf8)
	}
	if b.FilesTypes.Utf16 != 0 {
		t.Errorf("Expected FilesTypes.Utf16 0, got %d", b.FilesTypes.Utf16)
	}
	if b.FilesTypes.EightBit != 0 {
		t.Errorf("Expected FilesTypes.EightBit 0, got %d", b.FilesTypes.EightBit)
	}
	if b.FilesTypes.NonText != 1 {
		t.Errorf("Expected FilesTypes.NonText 1, got %d", b.FilesTypes.NonText)
	}

	if b.EolStyles.Total != 0 {
		t.Errorf("Expected EolStyles.Total 0, got %d", b.EolStyles.Total)
	}
	if b.EolStyles.Windows != 0 {
		t.Errorf("Expected EolStyles.Windows 0, got %d", b.EolStyles.Windows)
	}
	if b.EolStyles.Unix != 0 {
		t.Errorf("Expected EolStyles.Unix 0, got %d", b.EolStyles.Unix)
	}
	if b.EolStyles.Mac != 0 {
		t.Errorf("Expected EolStyles.Mac 0, got %d", b.EolStyles.Mac)
	}
	if b.EolStyles.Mixed != 0 {
		t.Errorf("Expected EolStyles.Mixed 0, got %d", b.EolStyles.Mixed)
	}
}

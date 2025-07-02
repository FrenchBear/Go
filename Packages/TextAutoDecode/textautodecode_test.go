// textautodecode_test.go
// Tests for package textautodecode
//
// 2025-06-23	PV		First version
// 2026-07-02 	PV 		External test project moved along package itself

package TextAutoDecode

import (
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {	
	test(t, `C:\DocumentsOD\Doc tech\Encodings\inexistent`, TFE_FileError)
	test(t, `C:\utils\bookApps\astructw.exe`, TFE_NotText)
	test(t, `C:\DocumentsOD\Doc tech\Encodings\prenoms-empty.txt`, TFE_Empty)
	test(t, `C:\DocumentsOD\Doc tech\Encodings\prenoms-ascii.txt`, TFE_ASCII)
	test(t, `C:\DocumentsOD\Doc tech\Encodings\prenoms-utf8bom.txt`, TFE_UTF8BOM)
	test(t, `C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16lebom.txt`, TFE_UTF16LEBOM)
	test(t, `C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16bebom.txt`, TFE_UTF16BEBOM)
	test(t, `C:\DocumentsOD\Doc tech\Encodings\prenoms-utf8.txt`, TFE_UTF8)
	test(t, `C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16le.txt`, TFE_UTF16LE)
	test(t, `C:\DocumentsOD\Doc tech\Encodings\prenoms-utf16be.txt`, TFE_UTF16BE)
	test(t, `C:\DocumentsOD\Doc tech\Encodings\prenoms-1252.txt`, TFE_EightBit)
}

func test(t *testing.T, filename string, expected TextFileEncoding) {
		tad, err := ReadTextFile(filename)
	if err != nil {
		if expected == TFE_FileError {
			return 
		} 
		t.Errorf("%-65.65s Err: %v\n", filename, err)
		return
	}

	if tad.Encoding!=expected {
		t.Errorf("Decoding %s, expected %s, got %s", filename, expected, tad.Encoding)
		return
	}
	if tad.Encoding==TFE_NotText || tad.Encoding==TFE_Empty{
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
		t.Errorf("Decoding %s, got \n%s\ninstead of\n%s\n", filename, tad.Text[:l], beginning)
	}
}
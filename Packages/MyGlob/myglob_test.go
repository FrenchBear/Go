// myglob_test.go
// Tests for MyGlob package
//
// 2025-07-01	PV 		Converted from Rust by Gemini
// 2025-07-13   PV      Tests with chinese characters
// 2025-08-11   PV      Added getRoot tests
// 2025-09-07   PV      Added MaxDepth tests

package MyGlob

import (
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup for search tests
	setupSearchTests()
	// Run tests
	code := m.Run()
	// Teardown for search tests
teardownSearchTests()
	// Exit
	os.Exit(code)
}

// -----------------------------------------------------------------------------
// Tests for regexp conversions

func TestRegexpConversions(t *testing.T) {
	t.Run("Simple constant string", func(t *testing.T) {
		globOneSegmentTest(t, "Pomme", "pomme", true)
		globOneSegmentTest(t, "Pomme", "pommerol", false)
	})

	t.Run("* pattern", func(t *testing.T) {
		globOneSegmentTest(t, "*", "rsgresp.d", true)
		globOneSegmentTest(t, "*.d", "rsgresp.d", true)
		globOneSegmentTest(t, "*.*", "rsgresp.d", true)
		globOneSegmentTest(t, "*.*", "rsgresp", false)
	})

	t.Run("** pattern", func(t *testing.T) {
		_, err := globToSegments("**.d" + string(os.PathSeparator))
		if err == nil {
			t.Errorf("Expected error for **.d, got nil")
		}
		globOneSegmentTest(t, "**", "", true)
	})

	t.Run("Alternations", func(t *testing.T) {
		globOneSegmentTest(t, "a{b,c}d", "abd", true)
		globOneSegmentTest(t, "a{b,c}d", "ad", false)
		globOneSegmentTest(t, "a{{b,c},{d,e}}f", "acf", true)
		globOneSegmentTest(t, "a{{b,c},{d,e}}f", "adf", true)
		globOneSegmentTest(t, "a{{b,c},{d,e}}f", "acdf", false)
		globOneSegmentTest(t, "a{b,c}{d,e}f", "acdf", true)
		globOneSegmentTest(t, "file.{cs,py,rs,vb}", "file.bat", false)
		globOneSegmentTest(t, "file.{cs,py,rs,vb}", "file.rs", true)
	})

	t.Run("? pattern", func(t *testing.T) {
		globOneSegmentTest(t, "file.?s", "file.rs", true)
		globOneSegmentTest(t, "file.?s", "file.cds", false)
	})

	t.Run("Character classes", func(t *testing.T) {
		globOneSegmentTest(t, "file.[cr]s", "file.rs", true)
		globOneSegmentTest(t, "file.[cr]s", "file.cs", true)
		globOneSegmentTest(t, "file.[cr]s", "file.py", false)
		globOneSegmentTest(t, "file.[a-r]s", "file.rs", true)
		globOneSegmentTest(t, "file.[-+]s", "file.-s", true)
		globOneSegmentTest(t, "file.[!abc]s", "file.rs", true)
		globOneSegmentTest(t, "file.[!abc]s", "file.cs", false)
		globOneSegmentTest(t, "file.[]]s", "file.]s", true)
		globOneSegmentTest(t, "file.[!]]s", "file.[s", true)
		globOneSegmentTest(t, `file[\d].cs`, "file1.cs", true)
		globOneSegmentTest(t, `file[\D].cs`, "filed.cs", true)
	})
}

func globOneSegmentTest(t *testing.T, globPattern, testString string, isMatch bool) {
	segments, err := globToSegments(globPattern + string(os.PathSeparator))
	if err != nil {
		t.Errorf("globToSegments failed for %s: %v", globPattern, err)
		return
	}

	if globPattern != "**" && len(segments) != 1 {
		t.Errorf("Expected 1 segment for %s, got %d", globPattern, len(segments))
		return
	}

	if len(segments) == 0 {
		return
	}

	switch s := segments[0].(type) {
	case ConstantSegment:
		if (strings.EqualFold(s.Value, testString)) != isMatch {
			t.Errorf("Constant match failed for %s with %s", globPattern, testString)
		}
	case RecurseSegment:
		if !isMatch {
			t.Errorf("Recurse match failed for %s with %s", globPattern, testString)
		}
	case FilterSegment:
		if s.Regexp.MatchString(testString) != isMatch {
			t.Errorf("Filter match failed for %s with %s (regex: %s)", globPattern, testString, s.Regexp.String())
		}
	}
}

// -----------------------------------------------------------------------------
// Tests for search functionality

func setupSearchTests() {
	_ = os.MkdirAll(`C:\Temp\search1\fruits`, 0755)
	_ = os.MkdirAll(`C:\Temp\search1\légumes`, 0755)
	_ = os.MkdirAll(`C:\Temp\search1\我爱你`, 0755)
	_ = os.MkdirAll(`C:\Temp\search1\我爱你\\Ƥḭҽɾɾҽ ѵìǫłҽղէ`, 0755)
	_ = os.WriteFile(`C:\Temp\search1\fruits et légumes.txt`, []byte("Des fruits et des légumes"), 0644)
	_ = os.WriteFile(`C:\Temp\search1\info`, []byte("Information"), 0644)
	_ = os.WriteFile(`C:\Temp\search1\fruits\pomme.txt`, []byte("Pomme"), 0644)
	_ = os.WriteFile(`C:\Temp\search1\fruits\poire.txt`, []byte("Poire"), 0644)
	_ = os.WriteFile(`C:\Temp\search1\fruits\ananas.txt`, []byte("Ananas"), 0644)
	_ = os.WriteFile(`C:\Temp\search1\fruits\tomate.txt`, []byte("Tomate"), 0644)
	_ = os.WriteFile(`C:\Temp\search1\légumes\épinard.txt`, []byte("Épinard"), 0644)
	_ = os.WriteFile(`C:\Temp\search1\légumes\tomate.txt`, []byte("Tomate"), 0644)
	_ = os.WriteFile(`C:\Temp\search1\légumes\pomme.de.terre.txt`, []byte("Pomme de terre"), 0644)
	_ = os.WriteFile(`C:\Temp\search1\我爱你\你好世界.txt`,                  []byte("Hello world"), 0644)
    _ = os.WriteFile(`C:\Temp\search1\我爱你\tomate.txt`,                    []byte("Hello Tomate"), 0644)
    _ = os.WriteFile(`C:\Temp\search1\我爱你\Ƥḭҽɾɾҽ ѵìǫłҽղէ\tomate.txt`,     []byte("Hello Tomate"), 0644)
    _ = os.WriteFile(`C:\Temp\search1\我爱你\Ƥḭҽɾɾҽ ѵìǫłҽղէ\Aé♫山𝄞🐗.txt`,  []byte("Random 1"), 0644)
    _ = os.WriteFile(`C:\Temp\search1\我爱你\Ƥḭҽɾɾҽ ѵìǫłҽղէ\œæĳøß≤≠Ⅷﬁﬆ.txt`, []byte("Random 2"), 0644)
}

func teardownSearchTests() {
	_ = os.RemoveAll(`C:\Temp\search1`)
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name          string
		glob          string
		autorecurse   bool
		ignore        []string
		maxDepth      int
		expectedFiles int
		expectedDirs  int
	}{
		// Basic testing
		{"InfoFile", `C:\Temp\search1\info`, false, nil, 0, 1, 0},
		{"AllInRoot", `C:\Temp\search1\*`, false, nil, 0, 2, 3},
		{"TxtInRoot", `C:\Temp\search1\*.*`, false, nil, 0, 1, 0},
		{"FilesInFruits", `C:\Temp\search1\fruits\*`, false, nil, 0, 4, 0},
		{"PFilesInTwoDirs", `C:\Temp\search1\{fruits,légumes}\p*`, false, nil, 0, 3, 0},
		{"RecursivePFiles", `C:\Temp\search1\**\p*`, false, nil, 0, 3, 0},
		{"RecursiveTxtFiles", `C:\Temp\search1\**\*.txt`, false, nil, 0, 13, 0},
		{"RecursiveDoubleExt", `C:\Temp\search1\**\*.*.*`, false, nil, 0, 1, 0},
		{"FilesInLegumes", `C:\Temp\search1\légumes\*`, false, nil, 0, 3, 0},
		{"ComplexFilter", `C:\Temp\search1\*s\to[a-z]a{r,s,t}e.t[xX]t`, false, nil, 0, 2, 0},

		// Multibyte runes
		{"Multibytes1", `C:\Temp\search1\**\*爱*\*a*.txt`, false, nil, 0, 1, 0},
		{"Multibytes2", `C:\Temp\search1\**\*爱*\**\*a*.txt`, false, nil, 0, 3, 0},
		{"Multibytes3", `C:\Temp\search1\我爱你\**\*🐗*`, false, nil, 0, 1, 0},

		// Testing autorecurse
		{"AutorecurseTxtOff", `C:\Temp\search1\*.txt`, false, nil, 0, 1, 0},
		{"AutorecurseTxtOn", `C:\Temp\search1\*.txt`, true, nil, 0, 13, 0},
		{"AutorecurseRootOff", `C:\Temp\search1`, false, nil, 0, 0, 1},
		{"AutorecurseRootOn", `C:\Temp\search1`, true, nil, 0, 14, 4},
		{"AutorecurseRootOnEndSlash", `C:\Temp\search1\`, true, nil, 0, 14, 4},		// Test with final \

		// Testing ignore
		{"IgnoreLegumes", `C:\Temp\search1\**\*.txt`, false, []string{"Légumes"}, 0, 10, 0},
		{"IgnoreLegumesAndOther", `C:\Temp\search1\**\*.txt`, false, []string{"Légumes","我爱你"}, 0, 5, 0},

		// Testing MaxDepth
		{"MaxDepth1", `C:\Temp\search1\**\*.txt`, true, nil, 1, 10, 0},
		{"MaxDepth2", `C:\Temp\search1\**\*.txt`, true, nil, 2, 13, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := New(tt.glob).Autorecurse(tt.autorecurse).MaxDepth(tt.maxDepth).ChannelSize(10)
			for _, ignore := range tt.ignore {
				builder.AddIgnoreDir(ignore)
			}
			gs, err := builder.Compile()
			if err != nil {
				t.Fatalf("Compile failed: %v", err)
			}

			nf, nd := 0, 0
			for m := range gs.Explore() {
				if m.Err != nil {
					t.Errorf("Explore error: %v", m.Err)
					continue
				}
				if m.IsDir {
					nd++
				} else {
					nf++
				}
			}

			if nf != tt.expectedFiles || nd != tt.expectedDirs {
				t.Errorf("got (files: %d, dirs: %d), want (files: %d, dirs: %d)", nf, nd, tt.expectedFiles, tt.expectedDirs)
			}
		})
	}
}

func TestSearchErrors(t *testing.T) {
	t.Run("InvalidGlob", func(t *testing.T) {
			_, err := New(`C:\**z\\z`).Compile()
			if err == nil {
				t.Error("Expected error for invalid glob, got nil")
			}
		})

	t.Run("InvalidRegex", func(t *testing.T) {
			_, err := New(`C:\[\d&&\p{ascii]`).Compile()
			if err == nil {
				t.Error("Expected error for invalid regex, got nil")
			}
		})

	t.Run("UnclosedBracket", func(t *testing.T) {
		_, err := New(`C:\[Hello`).Compile()
		if err == nil {
			t.Error("Expected error for invalid regex, got nil")
		}
	})

}

// -----------------------------------------------------------------------------
// Tests for getRoot

func TestGetRoot(t *testing.T) {
    tgr(t, "", ".", "*");
    tgr(t, "*", ".", "*");
    tgr(t, "C:", "C:", "");
    tgr(t, "C:\\", "C:\\", "");
    tgr(t, "file.ext", "file.ext", "");
    tgr(t, "C:file.ext", "C:file.ext", "");
    tgr(t, "C:\\file.ext", "C:\\file.ext", "");
    tgr(t, "path\\file.ext", "path\\file.ext", "");
    tgr(t, "path\\*.jpg", "path\\", "*.jpg");
    tgr(t, "path\\**\\*.jpg", "path\\", "**\\*.jpg");
    tgr(t, "C:path\\file.ext", "C:path\\file.ext", "");
    tgr(t, "C:\\path\\file.ext", "C:\\path\\file.ext", "");
    tgr(t, "\\\\server\\share", "\\\\server\\share", "");
    tgr(t, "\\\\server\\share\\", "\\\\server\\share\\", "");
    tgr(t, "\\\\server\\share\\file.txt", "\\\\server\\share\\file.txt", "");
    tgr(t, "\\\\server\\share\\path\\file.txt", "\\\\server\\share\\path\\file.txt", "");
    tgr(t, "\\\\server\\share\\*.jpg", "\\\\server\\share\\", "*.jpg");
    tgr(t, "\\\\server\\share\\path\\*.jpg", "\\\\server\\share\\path\\", "*.jpg");
    tgr(t, "\\\\server\\share\\**\\*.jpg", "\\\\server\\share\\", "**\\*.jpg");
}

func tgr(t *testing.T, pat, root, rem string) {
	r, s := getRoot(pat)	
	if r != root {
		t.Errorf("Pattern %s: Expected root %s, got %s", pat, root, r)
	}
	if s != rem {
		t.Errorf("Pattern: %s, Expected remainder %s, got %s", pat, rem, s)
	}
}
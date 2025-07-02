// MyMarkup.go
// MyMarkup package, my own implementation of Markup search
//
// 2025-07-02	PV 		First version
//
// MyMarkup use pecialized brackets for formatting text:
// ⟪Bold⟫           ~W  ~X
// ⟨Italic⟩         ~w  ~x
// ⌊Underline⌋      ~D  ~F
// ⌈Striketrough⌉   ~Q  ~S
// ⟦Color1⟧         ~c  ~v  Cyan
// ⦃Color2⦄         ~C  ~V  Yellow
// ⟮⟯               ~à  ~)  (Unused for now)
// ¬ (AltGr+7) sets left margin

package MyMarkup

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	libVersion = "1.0.0"
	showLimits = false 		// For dev
)

const (
// Styles
	STYLE_CLEAR = "\x1b[0m"
	STYLE_BOLD_ON = "\x1b[1m"
	STYLE_BOLD_OFF = "\x1b[22m" // Clears Dim and Bold
	STYLE_DIM_ON = "\x1b[2m"
	STYLE_DIM_OFF = "\x1b[22m" // Clears Dim and Bold
	STYLE_ITALIC_ON = "\x1b[3m"
	STYLE_ITALIC_OFF = "\x1b[23m"
	STYLE_UNDERLINE_ON = "\x1b[4m"
	STYLE_UNDERLINE_OFF = "\x1b[24m"
	STYLE_BLINK_ON = "\x1b[5m"
	STYLE_BLINK_OFF = "\x1b[25m"
	STYLE_REVERSE_ON = "\x1b[7m"
	STYLE_REVERSE_OFF = "\x1b[27m"
	STYLE_HIDDEN_ON = "\x1b[8m"
	STYLE_HIDDEN_OFF = "\x1b[28m"
	STYLE_STRIKETHROUGH_ON = "\x1b[9m"
	STYLE_STRIKETHROUGH_OFF = "\x1b[29m"

// Colors
	FG_BLACK = "\x1b[30m"
	FG_RED = "\x1b[31m"
	FG_GREEN = "\x1b[32m"
	FG_YELLOW = "\x1b[33m"
	FG_BLUE = "\x1b[34m"
	FG_MAGENTA = "\x1b[35m"
	FG_CYAN = "\x1b[36m"
	FG_WHITE = "\x1b[37m"
	FG_DEFAULT = "\x1b[39m"
	FG_BRIGHT_BLACK = "\x1b[90m"
	FG_BRIGHT_RED = "\x1b[91m"
	FG_BRIGHT_GREEN = "\x1b[92m"
	FG_BRIGHT_YELLOW = "\x1b[93m"
	FG_BRIGHT_BLUE = "\x1b[94m"
	FG_BRIGHT_MAGENTA = "\x1b[95m"
	FG_BRIGHT_CYAN = "\x1b[96m"
	FG_BRIGHT_WHITE = "\x1b[97m"

	BG_BLACK = "\x1b[40m"
	BG_RED = "\x1b[41m"
	BG_GREEN = "\x1b[42m"
	BG_YELLOW = "\x1b[43m"
	BG_BLUE = "\x1b[44m"
	BG_MAGENTA = "\x1b[45m"
	BG_CYAN = "\x1b[46m"
	BG_WHITE = "\x1b[47m"
	BG_DEFAULT = "\x1b[49m"
	BG_BRIGHT_BLACK = "\x1b[100m"
	BG_BRIGHT_RED = "\x1b[101m"
	BG_BRIGHT_GREEN = "\x1b[102m"
	BG_BRIGHT_YELLOW = "\x1b[103m"
	BG_BRIGHT_BLUE = "\x1b[104m"
	BG_BRIGHT_MAGENTA = "\x1b[105m"
	BG_BRIGHT_CYAN = "\x1b[106m"
	BG_BRIGHT_WHITE = "\x1b[107m"	
)

// Version returns the library version.
func Version() string {
	return libVersion
}

func RenderMarkup(txt_str string) {
        // To simplify code, ensure that string always ends with \n
       if !strings.HasSuffix(txt_str, "\n") {
			txt_str += "\n"
		}

        width := get_terminal_width()

        if showLimits {
            for i:=0 ; i<width ; i++ {
                fmt.Print("─")
            }
            fmt.Println()
        }

        // ToDo: Tabs expansion

        word := ""
        len := 0

        col := 0
        tab := 0
        for _, c := range txt_str {

            switch c {
			case '⟪':
                    word += STYLE_BOLD_ON
                    continue
			case '⟫':
                    word += STYLE_BOLD_OFF
                    continue
                
			case '⟨':
                    word += STYLE_ITALIC_ON
                    continue
                
			case '⟩':
                    word += STYLE_ITALIC_OFF
                    continue
                
			case '⌊':
                    word += STYLE_UNDERLINE_ON
                    continue
                
			case '⌋':
                    word += STYLE_UNDERLINE_OFF
                    continue
                
			case '⟦':
                    word += FG_CYAN
                    continue
                
			case '⟧':
                    word += FG_DEFAULT
                    continue
                
			case '⦃':
                    word += FG_YELLOW
                    continue
                
			case '⦄':
                    word += FG_DEFAULT
                    continue
                
                case '\r': continue

			case '\n' :
                    if word!="" && !is_only_spaces(&word) {
                        if col + len <= width {
                            fmt.Print(word)
                            if showLimits {
                                col += len
                                for ; col < width ; {
                                    col += 1
                                    fmt.Print(" ")
                                }
                                fmt.Print("|")
                            }
                            fmt.Println()
                        } else {
                            if showLimits {
                                for ; col < width ; {
                                    col += 1
                                    fmt.Print(" ")
                                }
                                fmt.Print("|")
                            }
                            fmt.Println()
                            for i:=0 ; i<tab ; i++ {
                                fmt.Print(" ")
                            }
                            for ; strings.HasPrefix(word, " ") ; {
                                word = word[1:]
                                len -= 1
                            }
                            fmt.Print(word)
                            col = tab + len
                            if showLimits {
                                for ; col < width ; {
                                    col += 1
                                    fmt.Print(" ")
                                }
                                fmt.Print("|")
                            }
                            fmt.Println()
                        }
                    } else {
                        for ; col < width ; {
                            col += 1
                            fmt.Print(" ")
                        }
                        if showLimits {
                            fmt.Print("|")
                        }
                        fmt.Println()
                    }
                    word=""
                    len = 0
                    col = 0
                    tab = 0
                
			case '¬':
				fmt.Print(word)
				col += len
				tab = col
				word=""
				len = 0
                

			case ' ':
                    if word!="" {
                        if is_only_spaces(&word) {
                            word += string(c)
                            len += 1
                            continue
                        }

                        if col + len <= width {
                            fmt.Print(word)
                            col += len
                            word = ""
                            len = 0
                        } else {
                            if showLimits {
                                for ; col < width ; {
                                    col += 1
                                    fmt.Print(" ")
                                }
                                fmt.Print("|")
                            }
                            fmt.Println()
                            for i:=0 ; i<tab; i++ {
                                fmt.Print(" ")
                            }
                            col = tab
                            for ; strings.HasPrefix(word, " ") ; {
                                word = word[1:]
                                len -= 1
                            }

                            fmt.Print(word)
                            col += len
                            word = ""
                            len = 0
                        }
                    }
                    word += " "
                    len += 1
                
                default: 
                    if tab + len >= width - 1 {
                        // We can't accumulate char, it would be longer than width
                        if col > tab {
                            // if we have already printed some chars, we need to flush and start a new line
                            if showLimits {
                                for ; col < width ; {
                                    col += 1
                                    fmt.Print(" ")
                                }
                                fmt.Print("|")
                            }
                            fmt.Println()

                            for i:=0 ; i<tab; i++ {
                                fmt.Print(" ")
                            }
                            col = tab
                        }

                        for ; strings.HasPrefix(word, " ") ;{
                            word = word[1:]
                            len -= 1
                        }

                        if col + len >= width {
                            fmt.Println("{}|", word)

                            word = ""
                            len = 0
                            for i:=0 ; i<tab; i++ {
                                fmt.Print(" ")
                            }
                            col = tab
                        }
                    }

                    word += string(c)
                    len += 1
                }
            }
        
        fmt.Println()
    }

func is_only_spaces(txt_str *string) bool {	
	for _, c := range *txt_str {
		if c != ' ' {
			return false
		}
	}
	return true
}

func get_terminal_width() int {
	// Check if the program is running in a terminal.
	// This is important because file descriptors for pipes or files will cause an error.
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return 80
	}

	// Get the terminal width and height.
	// os.Stdout.Fd() returns the file descriptor for standard output.
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		//log.Fatalf("failed to get terminal size: %v", err)
		return 80
	}

	return width
}

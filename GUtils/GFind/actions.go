// actions.go, definitions of actions
//
// 2025-07-12 	PV 		First version

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type IAction interface {
	action(path string, info os.FileInfo, noAction bool, verbose bool)
	name() string
}

// ===============================================================
// Print action

type action_print struct {
	detailed_output bool
}

func (ctx *action_print) action(path string, info os.FileInfo, noAction bool, verbose bool) {
	if !info.IsDir() {
		if ctx.detailed_output {
			fileSize := info.Size()
			strFileSize := formatIntWithThousandsSeparator(int(fileSize)) // Insert thousands separators
			modifiedTime := info.ModTime().Local()
			strModifiedTime := modifiedTime.Format("02/01/2006 15:04:05")	// Format date and time d/%m/%Y %H:%M:%S

			fmt.Printf("%19s    %15s  %s\n", strModifiedTime, strFileSize, path)
		} else {
			fmt.Println(path)
		}
	} else {
		if ctx.detailed_output {
			modifiedTime := info.ModTime().Local()
			strModifiedTime := modifiedTime.Format("02/01/2006 15:04:05")	// Format date and time d/%m/%Y %H:%M:%S
			fmt.Printf("%19s    %-15s  %s\n", strModifiedTime, "<DIR>", path)
		} else {
			fmt.Printf("%s%c\n", path, os.PathSeparator)
		}
	}
}

func (ctx *action_print) name() string {
	if ctx.detailed_output {
		return "Dir"
	} else {
		return "Print"
	}
}

// formatIntWithThousandsSeparator formats an integer with a thousand separator.
func formatIntWithThousandsSeparator(n int) string {
	s := strconv.Itoa(n)
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

// ===============================================================
// Delete action (remove files)

type action_delete struct {
	recycle bool
}

func (ctx *action_delete) action(path string, info os.FileInfo, noAction bool, verbose bool) {
	if !info.IsDir() {
		qp := quotedPath(path)
		if !ctx.recycle {
			fmt.Println("DEL", qp)
			if !noAction {
				err := os.Remove(path)
				if err==nil {
					fmt.Println("File", qp, "deleted successfully.")
				} else {
					fmt.Println("*** Error deleting file", qp, ":", err)
				}

			}
		} else {
			fmt.Println("RECYCLE", qp)
			if !noAction {
				fmt.Println("-> Recycle bin not implemented in this Go version.")
			}
		}
	}
}

func (ctx *action_delete) name() string {
	if ctx.recycle {
		return "Delete files (use recycle bin for local files, permanently for remote files)"
	} else {
		return "Delete files (permanently)"
	}
}

func quotedPath(path string) string {
	if strings.Contains(path, " ") {
		return "\"" + path + "\""
	} else {
		return path
	}
}

// ===============================================================
// Rmdir action (remove directories)

type action_rmdir struct {
	recycle bool
}

func (ctx *action_rmdir) action(path string, info os.FileInfo, noAction bool, verbose bool) {
	if info.IsDir() {
		qp := quotedPath(path)
		if !ctx.recycle {
			fmt.Println("RS /S", qp)
			if !noAction {
				err := os.RemoveAll(path)
				if err==nil {
					fmt.Println("Dir", qp, "deleted successfully.")
				} else {
					fmt.Println("*** Error deleting dir", qp, ":", err)
				}

			}
		} else {
			fmt.Println("RECYCLE (dir)", qp)
			if !noAction {
				fmt.Println("-> Recycle bin not implemented in this Go version.")
			}
		}
	}
}

func (ctx *action_rmdir) name() string {
	if ctx.recycle {
		return "Delete directories (use recycle bin for local files, permanently for remote files)"
	} else {
		return "Delete directories (permanent)"
	}
}

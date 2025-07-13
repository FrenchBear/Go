// readdirstream.go
// A faster, better, using less memory version of os.ReadDir...
//
// 2025-07-13 	PV 		First version from Gemini

package main

import (
	"io/fs"
	"os"
)

// DirEntry is a struct that holds the directory entry and any potential error.
type DirEntry struct {
	Entry fs.DirEntry
	Err   error
}

// readDirStream reads directory entries in a separate goroutine and sends them to a channel.
func readDirStream(dirName string) <-chan DirEntry {
	// Create a channel to return the directory entries.
	// The buffer size can be tuned for performance.
	entries := make(chan DirEntry, 500)

	go func() {
		// Ensure the channel is closed when the goroutine finishes.
		defer close(entries)

		dir, err := os.Open(dirName)
		if err != nil {
			entries <- DirEntry{Err: err}
			return
		}
		defer dir.Close()

		for {
			// Read a batch of directory entries. A value of -1 would read all,
			// but reading in smaller batches allows for more responsive streaming.
			// A positive value like 100 strikes a good balance.
			subEntries, err := dir.ReadDir(100)
			for _, entry := range subEntries {
				entries <- DirEntry{Entry: entry}
			}

			// io.EOF signals that we've reached the end of the directory.
			if err != nil {
				if err.Error() != "EOF" {
					entries <- DirEntry{Err: err}
				}
				return
			}
		}
	}()

	return entries
}


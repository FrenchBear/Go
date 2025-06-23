// g24_files_1.go
// Learning go, System programming, files
//
// 2025-06-23	PV		First version

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

func main() {
	fmt.Println("Go files")

	//test_s1_s2()
	test_simple()
}

type S1 struct {
	F1 int
	F2 string
}
type S2 struct {
	F1   S1
	text []byte
}

// Implementing the io.Reader() interface for S1.
// Using pointer to S1 for changes to be persistent when the method exits
// Note that this methods deads from stdin, p argument is not read from but actually written to... (dumb example I'd say)
func (s *S1) Read(p []byte) (n int, err error) {
	fmt.Print("Give me your name: ")
	fmt.Scanln(&p)		// Altough the buffer is initially 2 bytes, scanln will expand it as needed
	s.F2 = string(p)
	return len(p), nil
}

// Implement the io.Writer interface for S1.
// Actually prints the content if p buffer and ignores s content (besides s.F1, but it's actual value is almost irrelevant)
func (s *S1) Write(p []byte) (n int, err error) {
	if s.F1 < 0 {
		return -1, nil
	}
	for i := 0; i < s.F1; i++ {
		fmt.Printf("%s ", p)
	}
	fmt.Println()
	return s.F1, nil
}

func (s S2) eof() bool {
	return len(s.text) == 0
}

// an implementation of bytes.Buffer.ReadByte from the standard library
func (s *S2) readByte() byte {
	// this function assumes that eof() check was done before
	temp := s.text[0]
	s.text = s.text[1:]
	return temp
}

func (s *S2) Read(p []byte) (n int, err error) {
	if s.eof() {
		err = io.EOF
		return
	}
	l := len(p)
	if l > 0 {
		for n < l {
			// The previous functions read from the given buffer until the buffer is empty. This is where we implement io.Reader for S2.
			p[n] = s.readByte()
			n++
			if s.eof() {
				s.text = s.text[0:0]
				break
			}
		}
	}
	return
}

func test_s1_s2() {
	// We initialize an S1 variable
	s1var := S1{4, "Hello"}
	fmt.Println(s1var)

	// Reading for the s1var variable using a buffer with 2 bytes.
	buf := make([]byte, 2)
	_, err := s1var.Read(buf)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Read:", s1var.F2)

	// We call the Write() method for s1var in order to write the contents of a byte slice.
	_, _ = s1var.Write([]byte("Hello There!"))

	s2var := S2{F1: s1var, text: []byte("Hello world!!")}
	// In the previous code, we initialize an S2 variable that is named s2var.
	// Read s2var.text
	r := bufio.NewReader(&s2var)
	// We now create a reader for s2var.
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
			// We keep reading from s2var until there is an io.EOF condition.
		} else if err != nil {
			fmt.Println("*", err)
			break
		}
		fmt.Println("**", n, string(buf[:n]))
	}
	fmt.Println()
}

// ==================================================================================

type HashedFile struct {
	File string
	Width, Height int
	Size int64
	Hash, FlipHash uint64
	HashedOn time.Time
}

// HashedFileReadWriter wraps HashedFile to provide Reader/Writer capabilities.
// It uses an internal buffer to hold the serialized data.
type HashedFileReadWriter struct {
	HashedFile HashedFile
	buffer     *bytes.Buffer
}

// NewHashedFileReadWriter creates a new HashedFileReadWriter.
// The initial HashedFile can be empty.
func NewHashedFileReadWriter(hf HashedFile) *HashedFileReadWriter {
	return &HashedFileReadWriter{
		HashedFile: hf,
		buffer:     new(bytes.Buffer),
	}
}

// Read implements the io.Reader interface for HashedFileReadWriter.
// It reads from the internal buffer, which is populated by the Write method
// or by explicitly serializing the HashedFile.
func (hfrw *HashedFileReadWriter) Read(p []byte) (n int, err error) {
	if hfrw.buffer.Len() == 0 {
		// If the buffer is empty, try to serialize the HashedFile
		// This makes it so you can call Read directly after creating
		// the HashedFileReadWriter with a populated HashedFile.
		data, err := json.Marshal(hfrw.HashedFile)
		if err != nil {
			return 0, err
		}
		hfrw.buffer.Write(data)
	}
	return hfrw.buffer.Read(p)
}

// Write implements the io.Writer interface for HashedFileReadWriter.
// It writes bytes into the internal buffer. When the entire serialized
// HashedFile is received, it attempts to unmarshal it.
func (hfrw *HashedFileReadWriter) Write(p []byte) (n int, err error) {
	n, err = hfrw.buffer.Write(p)
	if err != nil {
		return n, err
	}

	// Attempt to unmarshal the buffer content into HashedFile
	// This approach assumes the entire JSON is written in one go,
	// or in a way that allows for eventual complete unmarshaling.
	// For streaming JSON, a more robust parser would be needed.
	var tempHashedFile HashedFile
	if json.Unmarshal(hfrw.buffer.Bytes(), &tempHashedFile) == nil {
		hfrw.HashedFile = tempHashedFile
		// Optionally, clear the buffer after successful unmarshal if you
		// expect new data to be written.
		// hfrw.buffer.Reset()
	}
	return n, nil
}

func test_simple() {
	// --- Example 1: Writing and then Reading ---
	fmt.Println("--- Example 1: Writing and then Reading ---")
	hf1 := HashedFile{
		File:     "image1.jpg",
		Width:    800,
		Height:   600,
		Size:     123456,
		Hash:     1234567890,
		FlipHash: 9876543210,
		HashedOn: time.Now(),
	}

	hfrw1 := NewHashedFileReadWriter(HashedFile{}) // Start with an empty HashedFile in the ReadWriter

	// Write the HashedFile data (serialized)
	data, err := json.Marshal(hf1)
	if err != nil {
		fmt.Println("Error marshalling hf1:", err)
		return
	}

	// s1 := string(data)
	// fmt.Println("$1!", s1)
	// // $1! {"File":"image1.jpg","Width":800,"Height":600,"Size":123456,"Hash":1234567890,"FlipHash":9876543210,"HashedOn":"2025-06-23T21:16:32.6908327+02:00"}

	_, err = hfrw1.Write(data)
	if err != nil {
		fmt.Println("Error writing data:", err)
		return
	}

	fmt.Printf("HashedFile after writing: %+v\n", hfrw1.HashedFile)

	// Now read from it
	readBuffer := make([]byte, 1024)
	n, err := hfrw1.Read(readBuffer)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading data:", err)
		return
	}
	fmt.Printf("Read %d bytes: %s\n", n, string(readBuffer[:n]))

	var hf1Read HashedFile
	err = json.Unmarshal(readBuffer[:n], &hf1Read)
	if err != nil {
		fmt.Println("Error unmarshalling read data:", err)
		return
	}
	fmt.Printf("HashedFile unmarshaled from read data: %+v\n", hf1Read)

	// --- Example 2: Creating with a HashedFile and then Reading ---
	fmt.Println("\n--- Example 2: Creating with a HashedFile and then Reading ---")
	hf2 := HashedFile{
		File:     "document.pdf",
		Width:    0, // Not applicable for PDF
		Height:   0,
		Size:     987654,
		Hash:     1122334455,
		FlipHash: 6677889900,
		HashedOn: time.Now().Add(-24 * time.Hour), // A day ago
	}

	hfrw2 := NewHashedFileReadWriter(hf2) // Initialize with a HashedFile

	// Read directly from it
	readBuffer2 := make([]byte, 1024)
	n2, err := hfrw2.Read(readBuffer2)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading data:", err)
		return
	}
	fmt.Printf("Read %d bytes: %s\n", n2, string(readBuffer2[:n2]))

	var hf2Read HashedFile
	err = json.Unmarshal(readBuffer2[:n2], &hf2Read)
	if err != nil {
		fmt.Println("Error unmarshalling read data:", err)
		return
	}
	fmt.Printf("HashedFile unmarshaled from read data: %+v\n", hf2Read)

	// --- Example 3: Copying between HashedFileReadWriters ---
	fmt.Println("\n--- Example 3: Copying between HashedFileReadWriters ---")
	hfrwSource := NewHashedFileReadWriter(hf1) // Use hf1 again
	hfrwDest := NewHashedFileReadWriter(HashedFile{}) // Destination starts empty

	// Populate the source's buffer first
	_, err = hfrwSource.Read(make([]byte, 100)) // Just to trigger serialization
	if err != nil && err != io.EOF {
		fmt.Println("Error preparing source:", err)
		return
	}

	// io.Copy can be used since both implement Reader/Writer
	bytesCopied, err := io.Copy(hfrwDest, hfrwSource)
	if err != nil {
		fmt.Println("Error copying:", err)
		return
	}
	fmt.Printf("Copied %d bytes\n", bytesCopied)
	fmt.Printf("HashedFile in destination after copy: %+v\n", hfrwDest.HashedFile)
}

/*
in Go, I have a HashedFile struct:

type HashedFile struct { File string Width, Height int Size int64 Hash, FlipHash uint64 HashedOn time.Time
}

Can you show me an example implementing io.Reader and io.Writer interfaces for this struct? 

---------

Sure, implementing io.Reader and io.Writer directly on the HashedFile struct itself in a meaningful way can be a bit
unconventional, as these interfaces are typically used for byte streams. A HashedFile struct primarily holds metadata
about a file, not the file's content itself.

However, I can show you an example where we treat the serialized form of the HashedFile struct as the stream that
io.Reader and io.Writer operate on. We'll use JSON for serialization in this example, but you could use any encoding
(gob, protobuf, etc.).

---------

Explanation:

    HashedFileReadWriter Struct:
        This helper struct wraps the HashedFile and embeds a *bytes.Buffer. The bytes.Buffer acts as the underlying
        data store for the Read and Write operations.

    NewHashedFileReadWriter Function:
        A constructor to easily create instances of HashedFileReadWriter.

    Read(p []byte) (n int, err error) Method:
        This method satisfies the io.Reader interface.
        When Read is called, it first checks if the internal buffer is empty.
        If the buffer is empty, it means we need to serialize the HashedFile content into the buffer before reading. 
        This allows you to create a HashedFileReadWriter with a pre-populated HashedFile and then immediately Read its serialized form.
        It then delegates the actual reading to the Read method of the bytes.Buffer.

    Write(p []byte) (n int, err error) Method:
        This method satisfies the io.Writer interface.
        It writes the incoming bytes p into the internal bytes.Buffer.
        Crucially, after writing, it attempts to json.Unmarshal the entire content of the bytes.Buffer back into a HashedFile.
            Important Consideration: This approach assumes that the Write calls will eventually accumulate a complete and valid
            JSON representation of the HashedFile. If the JSON is sent in fragmented chunks, json.Unmarshal might fail until
            all parts are received. For truly streaming JSON, you'd need a more sophisticated JSON stream parser (e.g., using json.Decoder).
            I've added a comment to optionally hfrw.buffer.Reset() after a successful unmarshal if you expect to be writing new,
            distinct HashedFile data after the current one is fully received.

How to use it:

    Writing: You can Write the byte-serialized form of a HashedFile into the HashedFileReadWriter. As bytes are written, 
	the HashedFile inside the HashedFileReadWriter will be updated (if a complete and valid JSON object is formed).

    Reading: You can Read the byte-serialized form of the HashedFile that is currently held by the HashedFileReadWriter. 
	If the HashedFileReadWriter was initialized with a HashedFile, calling Read will first serialize that HashedFile into its internal buffer.

    io.Copy: Because HashedFileReadWriter implements both io.Reader and io.Writer, you can use standard Go io utilities like io.Copy
	to transfer the serialized HashedFile data between HashedFileReadWriter instances or with other io.Reader/io.Writer implementations
	(e.g., os.File, net.Conn, etc.).

This example provides a functional way to make your HashedFile struct interact with io.Reader and io.Writer by treating
its serialized representation as the byte stream. Remember to choose the serialization format (JSON, Gob, etc.) that
best fits your needs.

*/
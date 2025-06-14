// g17_interfaces.go
// Learning go, Interfaces
//
// 2025-06-14	PV		First version

// Key Principles of Go Interfaces:
// - Implicit Implementation: No implements keyword. If a type has the methods of an interface, it implicitly satisfies
//   that interface.
// - Small Interfaces: Go encourages small, single-method interfaces (like io.Reader, fmt.Stringer). This follows the
//   Interface Segregation Principle.
// - Interfaces are Contracts: They define a contract of behavior, separating what a type can do from how it does it.
// - Consumer-Defined Interfaces: Often, interfaces are defined in the package that uses them, rather than the package
//   that implements them. This promotes loose coupling and allows multiple implementations to satisfy the same consumer
//   contract.
// - Zero Value: The zero value of an interface type is nil. An interface value is nil if both its type and value
//   components are nil. A non-nil interface value holding a nil concrete type can still be non-nil.
//
// By leveraging these principles, Go interfaces provide a powerful and idiomatic way to design robust, flexible, and
// maintainable software.

package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
)

func main() {
	BasicInterfaces()
	StandardLibraryInterfaces()
	InterfaceComposition()
	TestingForMocking()
	TypeAssertionsAndTypesSwitches()
}

// ------------------------------------------------------------------

func BasicInterfaces() {
	fmt.Println("------ Basic interface: Defining behavior")
	c := Circle{Radius: 5}
	r := Rectangle{Width: 3, Height: 4}
	Measure(c)
	Measure(r)

	// Cast from interface -> Circle or Rectangle
	TestAndConvert(c)
	TestAndConvert(r)
	TestAndConvert(nil)
}

func TestAndConvert(s Shape) {
	// Type assertion with the "comma ok" idiom
	// This is the SAFE and RECOMMENDED way.
	if re, ok := s.(Rectangle); ok {
		// If 'ok' is true, it means 's' was indeed a Rectangle.
		// 'sq' is now of type Rectangle and you can access its specific fields/methods.
		fmt.Printf("The shape is a Rectangle! Width: %.2f  Height: %.2f\n", re.Width, re.Height)
		fmt.Printf("Area of Rectangle: %.2f\n", re.Area()) // Using method specific to Rectangle (or general Shape)
	} else {
		// If 'ok' is false, it means 's' was not a Rectangle (or was nil).
		fmt.Printf("The shape is NOT a Rectangle. It's of type %T or nil.\n", s)
	}

	// You can also use a type switch for multiple types
	fmt.Println("--- Using Type Switch ---")
	switch v := s.(type) {
	case Rectangle:
		fmt.Printf("Type switch: It's a Rectangle with width %.2f and height %.2f\n", v.Width, v.Height)
	case Circle:
		fmt.Printf("Type switch: It's a Circle with radius %.2f\n", v.Radius)
	case nil:
		fmt.Println("Type switch: The shape is nil.")
	default:
		fmt.Printf("Type switch: Unknown shape type: %T\n", v)
	}

	// !!! DANGER: Type assertion without "comma ok" will panic if the assertion fails !!!
	// Do NOT do this unless you are absolutely certain of the type.
	/*
		fmt.Println("\n--- DANGEROUS TYPE ASSERTION (Will Panic if not Rectangle) ---")
		RectanglePanic := s.(Rectangle) // This will panic if s is not a Rectangle
		fmt.Printf("Directly asserted Rectangle side: %.2f\n", RectanglePanic.Side)
	*/
}

// Shape is an interface that defines methods for calculating Area and Perimeter.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle struct
type Circle struct {
	Radius float64
}

// Area implements the Area method for Circle
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter implements the Perimeter method for Circle
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// Rectangle struct
type Rectangle struct {
	Width, Height float64
}

// Area implements the Area method for Rectangle
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter implements the Perimeter method for Rectangle
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Measure takes any type that implements the Shape interface
func Measure(s Shape) {
	fmt.Printf("Shape: %T\n", s)
	fmt.Printf("Area: %.2f\n", s.Area())
	fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
	fmt.Println("---")
}

// ------------------------------------------------------------------

// Go's standard library extensively uses interfaces. io.Reader and io.Writer are prime examples, enabling generic I/O
// operations.

// Benefits:
// Abstraction: processData doesn't care where the data comes from or where it goes. It only needs an io.Reader and an
//              io.Writer.
// Reusability: The same processData function can be used with files, network connections, in-memory buffers, or any
//              other type that implements these standard interfaces.
// Testability: You can easily mock io.Reader and io.Writer in unit tests to simulate various input/output scenarios
//              without actual file or network operations.

// processData takes an io.Reader and an io.Writer
func processData(reader io.Reader, writer io.Writer) error {
	fmt.Println("\n------ Standard Library Interfaces")

	// Copy data from reader to writer
	_, err := io.Copy(writer, reader)
	if err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}
	return nil
}

func StandardLibraryInterfaces() {
	// Example 1: Reading from a string and writing to standard output
	stringReader := bytes.NewBufferString("Hello from string!\n")
	fmt.Println("--- String to Stdout ---")
	err := processData(stringReader, os.Stdout)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Example 2: Reading from a file and writing to a buffer
	fileContent := []byte("This is file content.\n")
	fileName := "example.txt"
	err = os.WriteFile(fileName, fileContent, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer os.Remove(fileName) // Clean up the file

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var buffer bytes.Buffer
	fmt.Println("--- File to Buffer ---")
	err = processData(file, &buffer)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Buffer content:\n%s", buffer.String())
	}
}

// ------------------------------------------------------------------

// MyReaderWriter implements both io.Reader and io.Writer since there is only one member with no name,
// it's just a synonym of bytes.Buffer
type MyReaderWriter struct {
	bytes.Buffer
}

// The io.ReadWriter interface is a composition of io.Reader and io.Writer.
func InterfaceComposition() {
	fmt.Println("\n------ Interface composition")

	var rw io.ReadWriter // Declare a variable of type io.ReadWriter
	myRW := &MyReaderWriter{}
	rw = myRW // Assign MyReaderWriter to io.ReadWriter because it implements both Reader and Writer

	// Write some data
	_, err := rw.Write([]byte("Hello, ReadWriter!\n"))
	if err != nil {
		fmt.Println("Error writing:", err)
	}

	// Read the data back
	p := make([]byte, 100)
	n, err := rw.Read(p)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading:", err)
	}
	fmt.Printf("Read %d bytes: %s", n, p[:n])
}

// ------------------------------------------------------------------

// Database is an interface for database operations
type Database interface {
	GetUser(id int) (string, error)
	SaveUser(id int, name string) error
}

// RealDatabase implements Database for a real database
type RealDatabase struct {
	// connection details
}

func (r *RealDatabase) GetUser(id int) (string, error) {
	// Simulate database query
	if id == 1 {
		return "Alice", nil
	}
	return "", fmt.Errorf("user %d not found in real DB", id)
}

func (r *RealDatabase) SaveUser(id int, name string) error {
	// Simulate saving to database
	fmt.Printf("Saving user %d: %s to real DB\n", id, name)
	return nil
}

// MockDatabase implements Database for testing purposes
type MockDatabase struct {
	Users map[int]string
}

func (m *MockDatabase) GetUser(id int) (string, error) {
	if name, ok := m.Users[id]; ok {
		return name, nil
	}
	return "", fmt.Errorf("user %d not found in mock DB", id)
}

func (m *MockDatabase) SaveUser(id int, name string) error {
	m.Users[id] = name
	fmt.Printf("Saving user %d: %s to mock DB\n", id, name)
	return nil
}

// UserService uses the Database interface
type UserService struct {
	db Database
}

func (us *UserService) GetUserDetails(id int) (string, error) {
	return us.db.GetUser(id)
}

func (us *UserService) CreateUser(id int, name string) error {
	return us.db.SaveUser(id, name)
}

func TestingForMocking() {
	fmt.Println("\n------ Testing for Mocking")

	// Using RealDatabase
	realDB := &RealDatabase{}
	realUserService := &UserService{db: realDB}
	user, err := realUserService.GetUserDetails(1)
	if err != nil {
		fmt.Println("Real DB Error:", err)
	} else {
		fmt.Println("Real DB User:", user)
	}
	realUserService.CreateUser(2, "Bob")
	fmt.Println("---")

	// Using MockDatabase for testing
	mockDB := &MockDatabase{
		Users: map[int]string{
			10: "Charlie",
		},
	}
	mockUserService := &UserService{db: mockDB}
	user, err = mockUserService.GetUserDetails(10)
	if err != nil {
		fmt.Println("Mock DB Error:", err)
	} else {
		fmt.Println("Mock DB User:", user)
	}
	mockUserService.CreateUser(11, "David")
	user, err = mockUserService.GetUserDetails(11)
	if err != nil {
		fmt.Println("Mock DB Error:", err)
	} else {
		fmt.Println("Mock DB User after save:", user)
	}
}

// ------------------------------------------------------------------

// Message interface
type Message interface {
	Content() string
}

// TextMessage concrete type
type TextMessage struct {
	Text string
}

func (t TextMessage) Content() string {
	return t.Text
}

// ErrorMessage concrete type
type ErrorMessage struct {
	Code    int
	Details string
}

func (e ErrorMessage) Content() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Details)
}

// ProcessMessage processes different message types using a type switch
func ProcessMessage(m Message) {
	switch v := m.(type) {
	case TextMessage:
		fmt.Printf("Received Text Message: \"%s\"\n", v.Text)
	case ErrorMessage:
		fmt.Printf("Received Error Message (Code: %d, Details: \"%s\")\n", v.Code, v.Details)
	default:
		fmt.Printf("Received Unknown Message Type: %T, Content: \"%s\"\n", v, v.Content())
	}
}

// You can check the underlying concrete type of an interface value using type assertions or type switches. This is
// useful when you need to access methods or fields specific to a concrete type that are not part of the interface.
// Example: Processing different message types

func TypeAssertionsAndTypesSwitches() {
	fmt.Println("\n------ Type Assertions and Types Switches")

	msg1 := TextMessage{Text: "Hello, world!"}
	msg2 := ErrorMessage{Code: 500, Details: "Internal server error"}
	var msg3 Message = TextMessage{Text: "Another text"}

	ProcessMessage(msg1)
	ProcessMessage(msg2)
	ProcessMessage(msg3)

	// Example of type assertion (less common for dispatch, more for specific checks)
	if textMsg, ok := msg3.(TextMessage); ok {
		fmt.Printf("Directly accessed TextMessage: %s\n", textMsg.Text)
	}
}

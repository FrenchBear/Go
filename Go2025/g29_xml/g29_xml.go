// g29_xml.go
// Learning go, System programming, files, Working with XML
//
// 2025-06-27	PV		First version

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math/rand"
)

// Field must be public (start with uppercase)
type User struct {
	Name      string  `xml:"username"`
	ID        int     `xml:"id,attr"`
	FirstName string  `xml:"name>first"`
	LastName  string  `xml:"name>last"`
	Height    float32 `xml:"height,omitempty"`
	Male      bool    `xml:"ismale"`
	Year      int     `xml:"created"`
	Comment   string  `xml:",comment"` // Comment in the output
}

func main() {
	fmt.Printf("Go XML\n\n")

	exm_encode_decode()
	xml_streams()
	xml_pretty_print()
}

func exm_encode_decode() {
	fmt.Printf("Xml Records Marshalling/Unmarshalling\n\n")

	user := User{
		Name:      "Pierre Violent",
		ID:        23,
		FirstName: "Pierre",
		LastName:  "Violent",
		Height:    186,
		Year:      1965,
		Comment:   "Ceci est un test",
	}

	// Regular Structure
	// Encoding XML data -> Convert Go Structure to XML record with fields
	t, err := xml.Marshal(&user)
	// The XML.Marshal() function requires a pointer to a structure variable—its real
	// data type is an empty interface variable—and returns a byte slice with the encoded
	// information and an error variable.

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Value %s\n", t)
	}

	// Decoding XML data given as a string
	str := `
<User id="23">
	<username>Pierre Violent</username>
	<name>
		<first>Pierre</first>
		<last>Violent</last>
	</name>
	<height>186</height>
	<ismale>false</ismale>
	<created>1965</created>
	<!--Ceci est un test-->
</User>`

	// Convert string into a byte slice
	XMLRecord := []byte(str)
	// However, as XML.Unmarshal() requires a byte slice, you need to convert that string
	// into a byte slice before passing it to XML.Unmarshal().
	// Create a structure variable to store the result (note that comment is deserialized)
	temp := User{}
	err = xml.Unmarshal(XMLRecord, &temp)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Data type: %T with value %v\n", temp, temp)
	}
}

// ----------------------------------------------
// Process multiple records: streams

type Data struct {
	Key string "xml:\"key\"" // standard strings also work
	Val int    `xml:"value"`
}

//var DataRecords []Data

// DataRecordsWrapper is a wrapper struct to provide a root element for the slice of Data.
type DataRecordsWrapper struct {
	XMLName xml.Name `xml:"DataRecords"` // Define the root element name for the slice
	Records []Data   `xml:"Data"`        // Each element in the slice will be named "Data"
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

var MIN = 0
var MAX = 26

func getString(l int64) string {
	startChar := "A"
	temp := ""
	var i int64 = 1
	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == l {
			break
		}
		i++
	}
	return temp
}

// DeSerialize decodes a serialized slice with XML records
func DeSerialize(e *xml.Decoder, slice interface{}) error {
	return e.Decode(slice)
}

// Serialize serializes a slice with XML records
func Serialize(e *xml.Encoder, slice interface{}) error {
	return e.Encode(slice)
}

func xml_streams() {
	fmt.Printf("\n\nXml streams\n\n")

	DataRecords := DataRecordsWrapper{Records: []Data{}}

	// Create sample data
	var i int
	var t Data
	for i = 0; i < 2; i++ {
		t = Data{
			Key: getString(5),
			Val: random(1, 100),
		}
		DataRecords.Records = append(DataRecords.Records, t)
	}

	// bytes.Buffer is both an io.Reader and io.Writer
	buf := new(bytes.Buffer)

	encoder := xml.NewEncoder(buf)
	err := Serialize(encoder, DataRecords)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("After Serialize:", buf)

	decoder := xml.NewDecoder(buf)
	var temp DataRecordsWrapper
	err = DeSerialize(decoder, &temp)
	fmt.Println("\nAfter DeSerialize:")
	for index, value := range temp.Records {
		fmt.Println(index, value)
	}
}

// ----------------------------------------------
// Pretty printing

func xml_pretty_print() {
	fmt.Printf("\n\nXml pretty printing\n\n")

	user := User{
		Name:      "Pierre Violent",
		ID:        23,
		FirstName: "Pierre",
		LastName:  "Violent",
		Height:    186,
		Year:      1965,
		Comment:   "Ceci est un test",
	}
	PrettyPrint_record(user)
	fmt.Println()

	// Create sample data
	DataRecords := DataRecordsWrapper{Records: []Data{}}
	var i int
	var t Data
	for i = 0; i < 2; i++ {
		t = Data{
			Key: getString(5),
			Val: random(1, 100),
		}
		DataRecords.Records = append(DataRecords.Records, t)
	}

	s, err := PrettyPrint_stream(DataRecords)
	if err != nil {
		fmt.Println("PrettyPrint_stream error")
	} else {
		fmt.Println(s)
	}
}

func PrettyPrint_record(v interface{}) (err error) {
	b, err := xml.MarshalIndent(v, "", "\t")
	if err == nil {
		fmt.Println(string(b))
	}
	return err
}

func PrettyPrint_stream(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := xml.NewEncoder(buffer)
	encoder.Indent("", "\t")
	// The XML.NewEncoder() function returns a new Encoder that writes to a writer that is
	// passed as a parameter to XML.NewEncoder(). An Encoder writes XML values to an
	// output stream. Similarly to XML.MarshalIndent(), the SetIndent() method allows
	// you to apply a customizable indent to a stream.
	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// g28_json.go
// Learning go, System programming, files, Working with Json
//
// 2025-06-26	PV		First version

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
)

// Field must be public (start with uppercase)
type UseAll struct {
	Name    string `json:"username"`
	Surname string `json:"surname"`
	Year    int    `json:"created"`
}

// Ignoring empty fields in JSON
type NoEmpty struct {
	Name    string `json:"username"`
	Surname string `json:"surname"`
	Year    int    `json:"creationyear,omitempty"`
}

// Removing private fields and ignoring empty fields
type Password struct {
	Name    string `json:"username"`
	Surname string `json:"surname,omitempty"`
	Year    int    `json:"creationyear,omitempty"`
	Pass    string `json:"-"`
}

func main() {
	fmt.Printf("Go JSON\n\n")

	encode_decode()
	streams()
	json_pretty_print()
}

func encode_decode() {
	fmt.Printf("Records Marshalling/Unmarshalling\n\n")

	useall := UseAll{Name: "Mike", Surname: "Tsoukalos", Year: 2021}

	// Regular Structure
	// Encoding JSON data -> Convert Go Structure to JSON record with fields
	t, err := json.Marshal(&useall)
	// The json.Marshal() function requires a pointer to a structure variable—its real
	// data type is an empty interface variable—and returns a byte slice with the encoded
	// information and an error variable.

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Value %s\n", t)
	}

	// Decoding JSON data given as a string
	str := `{"username": "M.", "surname": "Ts", "created":2020}`

	// Convert string into a byte slice
	jsonRecord := []byte(str)
	// However, as json.Unmarshal() requires a byte slice, you need to convert that string
	// into a byte slice before passing it to json.Unmarshal().
	// Create a structure variable to store the result
	temp := UseAll{}
	err = json.Unmarshal(jsonRecord, &temp)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Data type: %T with value %v\n", temp, temp)
	}
}

// ----------------------------------------------
// Process multiple records: streams

type Data struct {
	Key string "json:\"key\"" // standard strings also work
	Val int    `json:"value"`
}

var DataRecords []Data

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

// DeSerialize decodes a serialized slice with JSON records
func DeSerialize(e *json.Decoder, slice interface{}) error {
	return e.Decode(slice)
}

// Serialize serializes a slice with JSON records
func Serialize(e *json.Encoder, slice interface{}) error {
	return e.Encode(slice)
}

func streams() {
	fmt.Printf("\n\nJson streams\n\n")

	// Create sample data
	var i int
	var t Data
	for i = 0; i < 2; i++ {
		t = Data{
			Key: getString(5),
			Val: random(1, 100),
		}
		DataRecords = append(DataRecords, t)
	}

	// bytes.Buffer is both an io.Reader and io.Writer
	buf := new(bytes.Buffer)

	encoder := json.NewEncoder(buf)
	err := Serialize(encoder, DataRecords)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print("After Serialize:", buf)

	decoder := json.NewDecoder(buf)
	var temp []Data
	err = DeSerialize(decoder, &temp)
	fmt.Println("After DeSerialize:")
	for index, value := range temp {
		fmt.Println(index, value)
	}
}

// ----------------------------------------------
// Pretty printing

func json_pretty_print() {
	fmt.Printf("\n\nJson pretty printing\n\n")

	useall := UseAll{Name: "Mike", Surname: "Tsoukalos", Year: 2021}
	PrettyPrint_record(useall)
	fmt.Println()
	
	// Create sample data
	var i int
	var t Data
	for i = 0; i < 2; i++ {
		t = Data{
			Key: getString(5),
			Val: random(1, 100),
		}
		DataRecords = append(DataRecords, t)
	}

	s, err := PrettyPrint_stream(DataRecords)
	if err!=nil {
		fmt.Println("PrettyPrint_stream error")
	} else {
		fmt.Println(s)
	}
}

func PrettyPrint_record(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err == nil {
		fmt.Println(string(b))
	}
	return err
}

func PrettyPrint_stream(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")
	// The json.NewEncoder() function returns a new Encoder that writes to a writer that is
	// passed as a parameter to json.NewEncoder(). An Encoder writes JSON values to an
	// output stream. Similarly to json.MarshalIndent(), the SetIndent() method allows
	// you to apply a customizable indent to a stream.
	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

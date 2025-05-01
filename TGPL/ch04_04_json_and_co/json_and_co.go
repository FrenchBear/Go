// Learning Go
// §4, JSON format and other external representations
// 2019-10-23	PV

package main

import (
	"encoding/asn1"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
)

// Movie is a simple struct
type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"` // omitmpty = do not generate for default value, false here
	Actors []string
}

var movies = []Movie{
	{Title: "Casablanca", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
	// ...
}

func main() {
	// JSON
	{
		data, err := json.MarshalIndent(movies, "", "    ")
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		fmt.Printf("JSON\n%s\n\n", data)
	}

	// XML
	{
		data, err := xml.MarshalIndent(movies, "", "    ")
		if err != nil {
			log.Fatalf("XML marshaling failed: %s", err)
		}
		fmt.Printf("XML\n%s\n\n", data)
	}

	// ASN.1
	{
		data, err := asn1.Marshal(movies)
		if err != nil {
			log.Fatalf("ASN1 marshaling failed: %s", err)
		}
		fmt.Printf("ASN.1\n%s\n\n", data)
	}

}

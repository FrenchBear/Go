// g11_maps.go
// Learning go, Maps
//
// 2025-06-10	PV		First version

package main

import (
	"fmt"
)

func main() {
	// Map literal
	m := map[string]int {
		"key0": 42,
		"key1": -1,
		"key2": 123,
	}
	fmt.Println(m)

	v, ok := m["key2"]
	if ok {
		fmt.Printf("m[\"key2\"]=%v\n", v)
	} else {
		fmt.Println("m[\"key2\"] not found")
	}

	v, ok = m["key3"]
	if ok {
		fmt.Printf("m[\"key3\"]=%v\n", v)
	} else {
		fmt.Println("m[\"key3\"] not found")
	}

	// Iterating over a map. By design, order of keys is randomized
	for k,v := range m {
		fmt.Println(k, "->", v)
	}

	// Delete the whole map (not just keys/values)
	m = nil;

}

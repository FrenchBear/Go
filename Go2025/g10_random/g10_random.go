// g10_random.go
// Learning go, random numbers
//
// 2025-06-09	PV		First version

package main

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
)

func main() {
	n := random(1, 7)
	fmt.Println(n)

	s := getString(12)
	fmt.Println(s)

	a, err := generateBytes(32)
	if err == nil {
		fmt.Println(a)
	}
}

// Generate a random number in [min..max[ range
func random(min, max int) int {
	return rand.Intn(max-min) + min
}

const MIN = 0
const MAX = 94

func getString(len int64) string {
	temp := ""
	startChar := "!"
	var i int64 = 1
	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == len {
			break
		}
		i++
	}
	return temp
}

func generateBytes(n int64) ([]byte, error) {
	b := make([]byte, n)
	_, err := crand.Read(b)	// Read fills b with cryptographically secure random bytes. It never returns an error, and always fills b entirely.
	if err != nil {
		return nil, err
	}
	return b, nil
}

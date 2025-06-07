// g02_module/greetings.go
// Learning go, Example of submodule
//
// 2025-06-04	PV		First version

package greetings

import "fmt"

// Hello returns a greeting for the named person
func Hello(name string) string {
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}

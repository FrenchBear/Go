// functions.go
// Learning go, example of source code split in two files
// Second source, must have the same package name
//
// 2025-06-04	PV		First version

package main

import "fmt"

// The function is exported, so it starts with a capital letter
func Bonjour() {
	//fmt.Println("Bonjour à tous")
	salut()
}

// Not exported
func salut() {
	fmt.Println("Salut à tous!")
}

// g15_reflection.go
// Learning go, Exploring refection in Go
//
// 2025-06-12	PV		First version

package main

import (
	"fmt"
	"reflect"
)

// MyStruct is a sample struct to demonstrate reflection
type MyStruct struct {
	Name    string
	Age     int
	IsActive bool
	unexportedField string // Unexported field
}

// Greet is a sample method for MyStruct
func (m MyStruct) Greet(message string) {
	fmt.Printf("Hello, my name is %s. %s\n", m.Name, message)
}

// privateMethod is a private method, reflection can't call it directly
func (m MyStruct) privateMethod() {
	fmt.Println("This is a private method.")
}

func main() {
	fmt.Println("--- 1. Inspecting Types and Values ---")
	strVar := "Hello Go!"
	intVar := 123
	myStructVar := MyStruct{Name: "Alice", Age: 30, IsActive: true, unexportedField: "secret"}

	// Getting the Type and Value of basic types
	typeOfStr := reflect.TypeOf(strVar)
	valueOfStr := reflect.ValueOf(strVar)
	fmt.Printf("strVar: Type=%v, Kind=%v, Value=%v\n", typeOfStr, typeOfStr.Kind(), valueOfStr)

	typeOfInt := reflect.TypeOf(intVar)
	valueOfInt := reflect.ValueOf(intVar)
	fmt.Printf("intVar: Type=%v, Kind=%v, Value=%v\n", typeOfInt, typeOfInt.Kind(), valueOfInt)

	// Getting the Type and Value of a struct
	typeOfMyStruct := reflect.TypeOf(myStructVar)
	valueOfMyStruct := reflect.ValueOf(myStructVar)
	fmt.Printf("myStructVar: Type=%v, Kind=%v, Value=%v\n", typeOfMyStruct, typeOfMyStruct.Kind(), valueOfMyStruct)

	fmt.Println("\n--- 2. Accessing Fields of Structs ---")
	// To modify a struct using reflection, you need to pass a pointer to it.
	// Otherwise, reflect.ValueOf will create a copy, and changes won't be reflected in the original.
	ptrToMyStruct := &myStructVar
	valueOfPtrToMyStruct := reflect.ValueOf(ptrToMyStruct)
	// Elem() gets the value that the pointer points to
	elemOfMyStruct := valueOfPtrToMyStruct.Elem()

	if elemOfMyStruct.Kind() == reflect.Struct {
		// Accessing fields by name
		nameField := elemOfMyStruct.FieldByName("Name")
		ageField := elemOfMyStruct.FieldByName("Age")
		isActiveField := elemOfMyStruct.FieldByName("IsActive")
		unexportedField := elemOfMyStruct.FieldByName("unexportedField") // Can be accessed, but not settable if unexported

		fmt.Printf("Original: Name=%v, Age=%v, IsActive=%v\n",						// , Unexported=%v\n",
			nameField.Interface(), ageField.Interface(), isActiveField.Interface())	// , unexportedField.Interface())  -> panic: reflect.Value.Interface: cannot return value obtained from unexported field or method

		// Modifying fields (only if settable)
		if nameField.CanSet() {
			nameField.SetString("Bob")
		} else {
			fmt.Println("Name field is not settable (likely because valueOfMyStruct was not a pointer's Elem())")
		}
		if ageField.CanSet() {
			ageField.SetInt(35)
		}
		if isActiveField.CanSet() {
			isActiveField.SetBool(false)
		}
		if unexportedField.CanSet() {
			unexportedField.SetString("new secret") // This will panic because it's unexported
		} else {
			fmt.Printf("Unexported field 'unexportedField' is not settable: CanSet() = %v\n", unexportedField.CanSet())
		}

		fmt.Printf("Modified: Name=%v, Age=%v, IsActive=%v\n",
			elemOfMyStruct.FieldByName("Name").Interface(),
			elemOfMyStruct.FieldByName("Age").Interface(),
			elemOfMyStruct.FieldByName("IsActive").Interface())
		fmt.Printf("Original struct after modification: %+v\n", myStructVar)
	}

	fmt.Println("\n--- 3. Calling Methods ---")
	// For method calls, the receiver must be addressable if the method has a pointer receiver
	// or if the method modifies the receiver.
	// In our case, Greet has a value receiver, so it works directly on the struct value.
	methodValue := reflect.ValueOf(myStructVar).MethodByName("Greet")
	if methodValue.IsValid() {
		args := []reflect.Value{reflect.ValueOf("How are you?")}
		methodValue.Call(args)
	} else {
		fmt.Println("Greet method not found or not callable.")
	}

	// Attempting to call a private method (won't work directly)
	privateMethodValue := reflect.ValueOf(myStructVar).MethodByName("privateMethod")
	if privateMethodValue.IsValid() {
		fmt.Println("Attempting to call privateMethod (this usually won't work):")
		// privateMethodValue.Call([]reflect.Value{}) // This would likely panic or not execute
	} else {
		fmt.Println("privateMethod not found or not callable (expected for unexported methods).")
	}

	fmt.Println("\n--- 4. Creating New Values ---")
	// Creating a new instance of MyStruct dynamically
	typeOfMyStructInstance := reflect.TypeOf(MyStruct{})
	newMyStructPtr := reflect.New(typeOfMyStructInstance) // Returns a Value representing a pointer
	newMyStruct := newMyStructPtr.Elem()                  // Get the struct value it points to

	// Set fields of the new struct
	if newMyStruct.Kind() == reflect.Struct {
		newMyStruct.FieldByName("Name").SetString("Charlie")
		newMyStruct.FieldByName("Age").SetInt(25)
		newMyStruct.FieldByName("IsActive").SetBool(true)
		fmt.Printf("Newly created struct: %+v\n", newMyStruct.Interface())
	}

	fmt.Println("\n--- 5. Working with Pointers ---")
	var intPtr *int
	intPtr = new(int)
	*intPtr = 42

	valueOfIntPtr := reflect.ValueOf(intPtr)
	fmt.Printf("Pointer Value: Kind=%v, IsNil=%v\n", valueOfIntPtr.Kind(), valueOfIntPtr.IsNil())

	// Dereferencing a pointer
	if valueOfIntPtr.Kind() == reflect.Ptr {
		elemOfIntPtr := valueOfIntPtr.Elem()
		fmt.Printf("Dereferenced Pointer Value: Kind=%v, Value=%v\n", elemOfIntPtr.Kind(), elemOfIntPtr.Interface())

		// Modifying the value through the pointer
		if elemOfIntPtr.CanSet() {
			elemOfIntPtr.SetInt(99)
			fmt.Printf("Modified value through pointer: %d\n", *intPtr)
		}
	}

	// Creating a new pointer to an int
	newIntPtr := reflect.New(reflect.TypeOf(0)) // Create a pointer to an int
	newIntPtr.Elem().SetInt(100)
	fmt.Printf("Newly created int pointer value: %d\n", newIntPtr.Elem().Int())

	fmt.Println("\n--- 6. Inspecting Slices and Maps ---")
	mySlice := []int{10, 20, 30}
	valueOfSlice := reflect.ValueOf(mySlice)

	fmt.Printf("Slice Length: %d, Capacity: %d\n", valueOfSlice.Len(), valueOfSlice.Cap())
	for i := range valueOfSlice.Len() {
		fmt.Printf("Slice element %d: %v\n", i, valueOfSlice.Index(i).Interface())
	}

	myMap := map[string]int{"one": 1, "two": 2, "three": 3}
	valueOfMap := reflect.ValueOf(myMap)

	fmt.Printf("Map Length: %d\n", valueOfMap.Len())
	for _, key := range valueOfMap.MapKeys() {
		value := valueOfMap.MapIndex(key)
		fmt.Printf("Map element: Key=%v, Value=%v\n", key.Interface(), value.Interface())
	}

	// Adding an element to a slice using reflection (less common)
	// Note: This often involves creating a new slice and reassigning,
	// as slices are value types internally.
	// For simple appends, direct append is much easier.
	// We need a settable value to append to, so we pass a pointer to the slice.
	slicePtr := &mySlice
	valueOfSlicePtr := reflect.ValueOf(slicePtr)
	elemOfSlice := valueOfSlicePtr.Elem()

	if elemOfSlice.Kind() == reflect.Slice && elemOfSlice.CanSet() {
		newElement := reflect.ValueOf(40)
		elemOfSlice.Set(reflect.Append(elemOfSlice, newElement))
		fmt.Printf("Slice after appending 40: %v\n", mySlice)
	}

	// Setting a map element using reflection
	if valueOfMap.CanSet() { // Will be false because maps are not settable in this way (value-copy)
		fmt.Println("Map is not settable directly using reflect.ValueOf(map) for modification.")
	} else {
		// To modify a map, you need a pointer to the map, similar to structs
		mapPtr := &myMap
		valueOfMapPtr := reflect.ValueOf(mapPtr)
		elemOfMap := valueOfMapPtr.Elem()

		if elemOfMap.Kind() == reflect.Map {
			newKey := reflect.ValueOf("four")
			newValue := reflect.ValueOf(4)
			elemOfMap.SetMapIndex(newKey, newValue)
			fmt.Printf("Map after adding 'four': %v\n", myMap)
		}
	}
}
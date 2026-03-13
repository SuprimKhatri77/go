package main

import "fmt"

// IndexInt searches for an integer `x` inside a slice of integers `s`.
// It returns the index of the first occurrence if found.
// If the value does not exist in the slice, it returns -1.
func IndexInt(s []int, x int) int {
	for i, value := range s { // iterate through the slice with index and value
		if value == x { // check if the current value matches x
			return i // return the index where the match is found
		}
	}
	return -1 // return -1 if the value is not found
}

// IndexString searches for a string `x` inside a slice of strings `s`.
// It returns the index of the first occurrence if found.
// If the value does not exist in the slice, it returns -1.
func IndexString(s []string, x string) int {
	for i, v := range s { // iterate through the slice
		if v == x { // check if the current string matches x
			return i // return the index where the match is found
		}
	}
	return -1 // return -1 if the string is not found
}

// Index is a generic version of the search function.
// T is a type parameter constrained by `comparable`, meaning
// any type used here must support the `==` and `!=` operators.
//
// This allows the function to work with multiple types
// (int, string, float, custom comparable structs, etc)
// without writing separate functions like IndexInt or IndexString.
func Index[T comparable](s []T, x T) int {
	for i, v := range s { // iterate through the slice of type T
		if v == x { // compare values using == (allowed because of comparable constraint)
			return i // return index if match found
		}
	}
	return -1 // return -1 if value is not found
}

// GenericsExample demonstrates how the specific and generic
// versions of the index search function can be used.
func GenericsExample() {

	// Example slice of integers
	si := []int{1, 2, 3, 4, 5}

	// Using the non-generic function specifically made for integers
	fmt.Println(IndexInt(si, 1))

	// Example slice of strings
	ss := []string{"apple", "mango", "gauva"}

	// Using the non-generic function specifically made for strings
	fmt.Println(IndexString(ss, "mango"))

	// Using the generic function.
	// Go compiler automatically infers that T = int
	fmt.Println(Index(si, 1))

	// Using the same generic function with strings.
	// Here T = string
	fmt.Println(Index(ss, "grapes"))
}

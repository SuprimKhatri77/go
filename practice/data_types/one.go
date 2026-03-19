package main

import "fmt"

func PracOne() {

	// Explicit variable declarations with types
	// We define the type manually (string, int, float64, bool)
	// Go will NOT allow assigning a different type later
	var name string = "suprim"
	var age int = 20
	var height float64 = 1.8
	var isStudent bool = true

	// If we tried:
	// age = "twenty"
	// ❌ Compile-time error: cannot use string as int

	// Implicit type declaration using :=
	// Go infers the type from the assigned value → here it's string
	city := "Kathmandu"

	// If we later try:
	// city = 100
	// ❌ Error: cannot use int as string
	// Because type is FIXED after first assignment

	// Basic printing using fmt.Println (adds newline automatically)
	fmt.Println("Name: ", name)
	fmt.Println("Age: ", age)
	fmt.Println("Height: ", height)
	fmt.Println("isStudent: ", isStudent)

	// Print the type of variable using %T
	// No newline here → next output will be on same line
	fmt.Printf("%T", city)

	// This will continue on the same line if no newline above
	fmt.Println("Hello Go!")

	// fmt.Print does NOT add newline automatically
	// Everything stays on same line unless you add \n
	fmt.Print("Age: ", age, "Height: ", height)

	// fmt.Sprintf formats and RETURNS a string (does not print)
	// Useful when you want to store formatted text
	message := fmt.Sprintf("My name is %v and I am %v years old", name, age)

	// Now we print the formatted string
	fmt.Println(message)

	// Multiple variable declaration using :=
	a, b := 10, 3

	// Arithmetic operations
	fmt.Println(a + b) // Addition → 13
	fmt.Println(a - b) // Subtraction → 7

	// If both were integers:
	// a / b → integer division (result = 3, not 3.33)
	// For float division, one must be float:
	// float64(a) / float64(b)

	// Formatted printing with type information
	fmt.Printf("Name: %s (type: %T)\nAge: %d (type: %T)", name, name, age, age)

	// %s → string
	// %d → integer
	// %T → type of variable

	// No newline at end → in terminals like zsh, you may see a '%' prompt
}

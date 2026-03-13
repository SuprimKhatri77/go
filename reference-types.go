package main

import "fmt"

func ReferenceTypes() {
	// 1. Declare two integers.
	// Go finds two empty slots in RAM and puts these numbers there.
	i, j := 42, 2701

	// 2. The '&' (Address-of) operator.
	// We aren't storing '42'; we are storing the LOCATION of 'i'.
	// Read this as: "p is a pointer to i".
	p := &i

	// 3. The '*' (Dereference) operator.
	// Read this as: "Go to the address stored in p and tell me the value."
	fmt.Println(*p) // Prints: 42

	// 4. Changing value through the pointer.
	// This tells the CPU: "Go to the address p is pointing at and put 21 there."
	// This directly overwrites the '42' that was in 'i'.
	*p = 21

	// 5. Proof: 'i' has changed because p was pointing directly at its memory.
	fmt.Println(i) // Prints: 21

	// 6. Redirecting the pointer.
	// Now 'p' stops looking at 'i' and starts looking at 'j'.
	p = &j

	// 7. Math through a pointer.
	// "Take the value at p (2701), divide by 37, and store it back at p."
	*p = *p / 37

	// 8. Proof: 'j' is now changed to 73.
	fmt.Println(j) // Prints: 73
}

package main

import "fmt"

/*
A 'struct' is a collection of fields (different data types grouped together).
Think of it like an Object in JS or a Class without methods in Java.
*/
type Vertex struct {
	X int // Exported (Public) because it starts with Capital
	Y int
}

type ShortHandVertex struct {
	X, Y int // we can group fields of the same type to keep it clean
}

// Global variable declarations
var (
	v1 = Vertex{1, 2} // Ordered: Go assumes 1 is X and 2 is Y

	v2 = Vertex{X: 1, Y: 2} // Named: Best practice for readability

	v3 = Vertex{} // Zero-Value: Both X and Y become 0 automatically

	// Pointer to Struct: This doesn't store the data, it stores the ADDRESS.
	// The type of 'ptr' is *Vertex.
	ptr = &Vertex{1, 2}

	v5 = Vertex{X: 1} // Partial: X=1, Y defaults to 0
)

func StructExample() {
	v := Vertex{1, 2}

	// Accessing fields via dot notation
	// 1e2 is scientific notation (1 * 10^2) = 100
	v.X = 1e2

	/*
	   POINTER INDIRECTION:
	   In languages like C, you'd have to write (*p).X = 1e5.
	   In Go, we can just write p.X.
	   Go automatically "dereferences" the pointer for you because it's nicer to read!
	*/
	p := &v
	p.X = 1e5 // Directly modifies 'v.X' in memory

	fmt.Println(v)

	// Note: When we print 'ptr', Go adds an '&' prefix to tell you it's a pointer.
	fmt.Println(v1, v2, v3, ptr, v5)
}

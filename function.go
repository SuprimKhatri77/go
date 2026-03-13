package main

import (
	"fmt"
	"math"
)

// compute is a "Higher-Order Function."
// It takes another function (fn) as an argument.
// For this to work, 'fn' must match the signature: func(float64, float64) float64
func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}

func AnonymousFunction() {
	// 1. ANONYMOUS FUNCTION ASSIGNMENT
	// We define a function without a name and assign it to the variable 'hypot'.
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}

	// We can call it directly via the variable name
	fmt.Println("Direct call (5, 12):", hypot(5, 12))

	// 2. PASSING FUNCTIONS AS ARGUMENTS
	// We pass our 'hypot' variable into 'compute'.
	fmt.Println("Compute with hypot:", compute(hypot))

	// 3. PASSING BUILT-IN FUNCTIONS
	// Since math.Pow has the signature func(float64, float64) float64,
	// it satisfies compute's requirement perfectly.
	fmt.Println("Compute with math.Pow:", compute(math.Pow))
}

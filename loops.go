package main

import (
	"fmt"
	"math"
)

func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}

/*
pow calculates x^n but caps the result at a maximum 'lim'.
It uses a "Short Statement" if-block, which is unique to Go.
*/
func pow(x, n, lim float64) float64 {
	// 1. INITIALIZATION: v := math.Pow(x, n) runs first.
	// 2. SCOPE: 'v' is only "alive" inside this if/else block.
	// 3. CONDITION: v < lim is checked immediately after initialization.
	if v := math.Pow(x, n); v < lim {
		return v // Case: 3^2 = 9. Since 9 < 10, return 9.
	} else {
		// %g is "General" floating-point formatting.
		// It's "Smart": it hides trailing zeros and uses scientific notation for huge numbers.
		fmt.Printf("%g >= %g\n", v, lim)
	}

	// If the condition is false, 'v' is still accessible in an 'else' block,
	// but here we just drop down to return the limit.
	return lim // Case: 4^2 = 16. Since 16 < 15 is false, return 15.
}
func Loops() {
	sum := 1

	// for i := 0; i < 10; i++ {
	// 	sum += i
	// }
	// fmt.Println("sum: ", sum)

	// infinite loop
	for {

		// If we dont have a break statement, the loop will run indefinitely
		if sum > 100 {
			break
		}
		sum += sum
	}
	fmt.Println("sum: ", sum)

	fmt.Println(sqrt(2), sqrt(-4))

	fmt.Println(
		"Result 1 (Under limit):", pow(3, 2, 10), // Returns 9
		"\nResult 2 (Over limit):", pow(4, 2, 15), // Returns 15
	)
}

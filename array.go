package main

import (
	"fmt"
)

func ArrayExample() {
	// Arrays have a FIXED size. [2]string is a different type than [3]string.
	var arr [2]string
	arr[0] = "Hello"
	arr[1] = "World"

	fmt.Println(arr[0], arr[1])
	fmt.Println(arr)

	// Array literal: size is part of the type.
	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)

	/* SLICES: A window into the array.
	   Syntax: [low : high]
	   - low: inclusive (start here)
	   - high: exclusive (stop BEFORE here)
	*/
	var s []int = primes[1:4] // Grabs indices 1, 2, and 3.
	fmt.Println(s)
}

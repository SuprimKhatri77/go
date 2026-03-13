package main

import "fmt"

// adder returns a "closure".
// Even though adder() finishes executing, the anonymous function it returns
// "closes over" the variable 'sum', keeping it alive in memory.
func adder() func(int) int {
	sum := 0 // This variable is captured by the returned function
	return func(x int) int {
		sum += x // It remembers 'sum' every time it is called
		return sum
	}
}

func ClosureExample() {
	// Each call to adder() creates a NEW, isolated scope.
	// 'pos' has its own 'sum', and 'neg' has its own 'sum'.
	pos, neg := adder(), adder()

	for i := 0; i < 10; i++ {
		// pos(i) adds positive numbers to its internal sum.
		// neg(-2*i) adds negative numbers to its separate internal sum.
		fmt.Printf("Pos Sum: %d | Neg Sum: %d\n", pos(i), neg(-2*i))
	}
}

package main

import "fmt"

func sum(n1, n2 int) int {
	return n1 + n2
}
func DeferExample() {
	fmt.Println("Counting")

	/*
	   DEFER LOGIC:
	   1. LIFO (Last-In, First-Out): Statements are pushed onto a stack.
	      The LAST one deferred is the FIRST one executed.
	   2. Arguments are evaluated IMMEDIATELY, but the
	      function call happens at the very end of the surrounding function.
	   3. PANIC SAFETY: Defer runs even if the function panics (crashes),
	      making it essential for 'Cleanup' tasks.

	   Stack visualization for this loop:
	   [ Push 0 ] -> [ Push 1 ] ... -> [ Push 9 ]

	   Execution (Popping):
	   9, 8, 7, 6, 5, 4, 3, 2, 1, 0
	*/
	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}

	defer fmt.Println("sum is: ", sum(5, 2))

	fmt.Println("done")
}

package main

import (
	"fmt"
	"strings"
)

func SliceExample() {
	names := [4]string{
		"Alice",
		"Bob",
		"Charlie",
		"David",
	}
	fmt.Println(names)

	// Slices 'a' and 'b' overlap at index 1 of the 'names' array.
	a := names[0:2] // [Alice, Bob]
	b := names[1:3] // [Bob, Charlie]
	fmt.Println(a, b)

	/*
	   REFERENCE TYPE BEHAVIOR:
	   Changing 'b[0]' is the same as changing 'names[1]'.
	   Since 'a[1]' also points to 'names[1]', 'a' reflects the change too.
	   Slices do not store data; they point to it.
	*/
	b[0] = "X"
	fmt.Println(a, b)  // a: [Alice, X], b: [X, Charlie]
	fmt.Println(names) // names: [Alice, X, Charlie, David]

	// SLICE LITERALS:
	// This looks like an array literal, but without the [number].
	// Go secretly creates an array, then builds a slice that points to it.
	q := []int{2, 3, 5, 7, 11, 13}
	fmt.Println(q)

	r := []bool{true, false, true, true, false, true}
	fmt.Println(r)

	// Slice of Anonymous Structs:
	// Common in API responses or test cases.
	s := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, false},
		{5, true},
	}
	fmt.Println(s)

	/*
	   RE-SLICING:
	   we can slice a slice to narrow the window.
	*/
	sl := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	sl = sl[1:4] // [2, 3, 4] | Original 1 is now "behind" the pointer.
	sl = sl[:3]  // [2, 3, 4] | No change here as length was already 3.
	sl = sl[2:]  // [4]       | Moves the start pointer forward 2 steps.
	fmt.Println(sl)
}

func SliceExampleTwo() {
	s := []int{2, 3, 5, 7, 11, 13}
	printSlice(s) // len=6, cap=6

	/*
	   THE 'CURTAIN' ANALOGY:
	   s[:0] sets length to 0. It's like closing the curtain.
	   The pointer is still at the start of the array (House #0).
	*/
	s = s[:0]
	printSlice(s) // len=0, cap=6 (Capacity remains because we didn't move the pointer)

	/* THE 'PANIC' ZONE:
	   If we uncomment the line below, it works because :4 stretches the RIGHT side.
	   The right side (High-bound) can look at the CAPACITY.
	*/
	// s = s[:4]

	/*
	   WHY s[2:] FAILS AFTER s[:0]:
	   The left side (Low-bound) can only move based on the CURRENT LENGTH.
	   Since len=0, there is no 'index 2' to move the pointer to.
	   Go stops we from 'walking through a wall' into memory we can't see.
	*/
	// s = s[2:] // This would panic!

	// To move the pointer forward when len is 0, we must use both bounds:
	s = s[2:4]    // Tells Go: "Move pointer to index 2 AND set length to 2."
	printSlice(s) // len=2, cap=4 [5 7]
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func SliceExampleThree() {
	// Zero value of a slice is nil.
	// PS: Zero value of a REFERENCE TYPE is nil.

	/* A 'nil' slice has no underlying array.
	   The "Pointer" inside the slice header is pointing to 0 (nothing).
	*/
	var s []int

	// Even though it's nil, Go is smart: len() and cap() won't crash.
	// They safely return 0.
	fmt.Println(s, len(s), cap(s)) // [] 0 0

	if s == nil {
		/*
		   This is the standard way to check if a slice has been initialized.
		*/
		fmt.Println("Slice is nil!")
	}
}

// declaring a slice using make
func SliceExampleFour() {
	/*
	   make([]T, len, cap) allocates an underlying array of type T.
	   - len: Initial number of elements visible in the slice.
	   - cap: Size of the allocated underlying array.
	   All elements are initialized to the type's zero-value (0 for int).

	   If capacity is not provided, it defaults to the same as the length.
	*/

	// a points to a new array of 5 ints. len=5, cap=5.
	a := make([]int, 5)
	printSlice(a) // [0 0 0 0 0] len=5 cap=5

	// b points to a new array of 5 ints, but the window is set to 0.
	// The memory is allocated, but no elements are accessible yet via index.
	b := make([]int, 0, 5)
	printSlice(b) // [] len=0 cap=5

	// c is a re-slice of b.
	// High-bound slicing ([:2]) is allowed up to the capacity (5).
	// The pointer stays at the start of the underlying array.
	c := b[:2]
	printSlice(c) // [0 0] len=2 cap=5

	/*
	   d := c[2:5]
	   1. Low-bound (2): Shifts the pointer forward 2 indices.
	   2. High-bound (5): Sets the end of the window to index 5 of the original capacity.

	   Result: Since the pointer moved forward by 2, the new capacity becomes 3 (5 - 2).
	   The slice now looks at indices 2, 3, and 4 of the original backing array.
	*/
	d := c[2:5]
	printSlice(d) // [0 0 0] len=3 cap=3
}

func TwoDimensionSlice() {
	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}

	board[0][0] = "A"
	board[2][2] = "B"
	board[1][2] = "C"
	board[1][0] = "D"
	board[0][2] = "E"

	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}
}

func AppendExample() {
	// 1. DYNAMIC GROWTH (RE-ALLOCATION MODE)
	// Starting with a 'nil' slice (len: 0, cap: 0).
	// Every time the capacity is exceeded, Go performs a 'growslice' operation:
	//   a) Allocates a new, larger block of memory.
	//   b) Copies existing elements to the new block.
	//   c) Switches the slice pointer to the new memory location.
	var nums []int

	fmt.Println("--- Dynamic Growth ---")
	for i := 0; i < 20; i++ {
		nums = append(nums, i)
		// Notice: Capacity doesn't just increment by 1.
		// Go rounds up to the nearest 'size class' to minimize OS memory requests.
		fmt.Printf("Len: %d, Cap: %d\n", len(nums), cap(nums))
	}

	// 2. PRE-ALLOCATED (APPEND-READY MODE)
	// BEST PRACTICE: Use make([]T, 0, capacity) when you know the final size.
	// Length is 0: The slice is empty and ready for append().
	// Capacity is 5: The memory is reserved upfront. No re-allocations will occur.
	fmt.Println("\n--- Pre-allocated (Len 0, Cap 5) ---")
	s := make([]int, 0, 5)

	for i := 0; i < 5; i++ {
		s = append(s, i) // Starts filling from index 0
		fmt.Printf("Len: %d, Cap: %d, Data: %v\n", len(s), cap(s), s)
	}

	// 3. INITIALIZED (INDEX-READY MODE)
	// Use make([]T, length) if you want to access by index (e.g., s[i] = val).
	// WARNING: Using append() here will add elements AFTER the five zeros.
	fmt.Println("\n--- Initialized (Len 5, Cap 5) ---")
	prefilled := make([]int, 5) // Creates [0, 0, 0, 0, 0]

	// This append will force a capacity grow to ~10 because index 0-4 are "full"
	prefilled = append(prefilled, 100)
	fmt.Printf("After append: Len: %d, Cap: %d, Data: %v\n", len(prefilled), cap(prefilled), prefilled)
}

// 'power' is a slice of integers.
// Go slices are "descriptors" of an underlying array.
var power = []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512}

func RangeExample() {
	// USE CASE 1: Both Index and Value
	// 'i' is the current index, 'v' is a COPY of the element at that index.
	fmt.Println("--- Both Index and Value ---")
	for i, v := range power {
		fmt.Printf("2**%d = %d\n", i, v)
	}

	// USE CASE 2: Index Only
	// If you only need the position, you can omit the second variable entirely.
	fmt.Println("\n--- Index Only ---")
	for i := range power {
		fmt.Printf("Processing index: %d\n", i)
	}

	// USE CASE 3: Value Only (Ignoring Index)
	// Use the blank identifier '_' to discard the index.
	// This is the most common way to "read" through a collection.
	fmt.Println("\n--- Value Only ---")
	sum := 0
	for _, v := range power {
		sum += v
	}
	fmt.Printf("Total Sum: %d\n", sum)
}

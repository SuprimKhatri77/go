// ============================================================
// TOPIC 8: Arrays
// ============================================================
// Arrays in Go are FIXED SIZE — size is part of the type itself.
// Unlike JS where arrays are dynamic, Go arrays cannot grow.
// [3]int and [4]int are completely different types — not compatible.
//
// In real Go code you'll use slices (Topic 9) almost always.
// But arrays are the foundation slices are built on — understand
// arrays and slices will click immediately.
// ============================================================

package main

import "fmt"

// --------------------------------------------------------
// Arrays in function signatures — size is part of the type
// --------------------------------------------------------
// [4]string and [5]string are different types entirely.
// A function expecting [4]string will NOT accept [5]string.
// This is why real Go code uses slices — works with any size.
// --------------------------------------------------------
func contains(routes [4]string, target string) bool {
	for _, v := range routes {
		if v == target {
			return true
		}
	}
	return false
}

func main() {

	// --------------------------------------------------------
	// 1. DECLARING ARRAYS — three ways
	// --------------------------------------------------------
	// Syntax: [size]type{values}
	// Size is mandatory and fixed at declaration — cannot change.
	// --------------------------------------------------------

	// Explicit size with values
	arr1 := [5]int{200, 201, 400, 404, 500}
	fmt.Println(arr1)

	// Let Go count the size — [...] infers size from values
	arr2 := [...]int{200, 201, 400, 404, 500} // size = 5, inferred
	fmt.Println(arr2)

	// Declare with zero values first, assign individually
	var arr3 [5]int
	arr3[0] = 200
	arr3[1] = 201
	fmt.Println(arr3) // [200 201 0 0 0] — unfilled slots = zero value

	// Zero values depend on element type:
	// int/float → 0, string → "", bool → false
	var intArr [3]int
	var strArr [3]string
	var boolArr [3]bool
	fmt.Println(intArr)  // [0 0 0]
	fmt.Println(strArr)  // [  ]
	fmt.Println(boolArr) // [false false false]

	// --------------------------------------------------------
	// 2. ARRAYS ARE FIXED — out of bounds behavior
	// --------------------------------------------------------
	// Two kinds of out-of-bounds errors — important distinction:
	//
	// COMPILE TIME error — compiler knows the index is too large:
	//   arr1[5] = 503  // ❌ index 5 out of bounds [5]
	//   Caught before program runs — free, safe
	//
	// RUNTIME PANIC — dynamic index, compiler can't check:
	//   i := 5
	//   arr1[i] = 503  // ❌ compiles fine, panics when it runs
	//   Program crashes in production — dangerous
	//
	// Always guard dynamic indices: if i < len(arr) { arr[i] = val }
	// --------------------------------------------------------

	// [3]int and [4]int are DIFFERENT TYPES — cannot compare or assign
	// a := [3]int{1, 2, 3}
	// b := [4]int{1, 2, 3, 4}
	// fmt.Println(a == b) // ❌ compile error — mismatched types

	// --------------------------------------------------------
	// 3. ARRAYS ARE VALUE TYPES — not references
	// --------------------------------------------------------
	// In JS: const copy = original → both point to same array
	//        mutating copy mutates original
	//
	// In Go: copy := original → Go creates a completely new copy
	//        mutating copy does NOT affect original
	//
	// Arrays live on the STACK (fixed size, short-lived)
	// Slices/maps live on the HEAP (dynamic, reference types)
	// --------------------------------------------------------
	original := [3]string{"admin", "user", "guest"}
	copied := original // ✅ named `copied` not `copy` — copy is a built-in function
	copied[0] = "superadmin"

	fmt.Println("original:", original) // [admin user guest] — unchanged ✅
	fmt.Println("copied:", copied)     // [superadmin user guest]

	// ⚠️  Avoid naming variables after built-in functions:
	// copy, len, append, make, new, cap, delete, panic, recover
	// These are NOT reserved keywords (won't error) but shadowing
	// them means you can't use the built-in in that scope.
	// Reserved keywords (func, for, if, return) cannot be used as names at all.

	// --------------------------------------------------------
	// 4. ITERATING ARRAYS
	// --------------------------------------------------------
	allowedRoutes := [4]string{"/health", "/users", "/posts", "/auth"}

	// Classic for loop with index
	for i := 0; i < len(allowedRoutes); i++ {
		fmt.Println(allowedRoutes[i])
	}

	// for range — cleaner, preferred
	for i, v := range allowedRoutes {
		fmt.Printf("%d: %s\n", i, v)
	}

	fmt.Println(contains(allowedRoutes, "/users"))   // true
	fmt.Println(contains(allowedRoutes, "/unknown")) // false

	// --------------------------------------------------------
	// 5. 2D ARRAYS — permission matrix
	// --------------------------------------------------------
	// Syntax: [rows][cols]type
	// Access: arr[row][col]
	// Note: trailing comma on last element is REQUIRED in Go
	//       when closing brace } is on its own line
	// --------------------------------------------------------
	roles := [3]string{"admin", "moderator", "guest"}
	permissions := [3]string{"read", "write", "delete"}

	mat := [3][3]bool{
		{true, true, true},   // admin
		{true, true, false},  // moderator
		{true, false, false}, // guest — trailing comma required here
	}

	for i := 0; i < len(mat); i++ {
		fmt.Printf("%s: ", roles[i])
		for j := 0; j < len(mat[i]); j++ {
			// permissions[j] not permissions[i] — j changes per column
			// permissions[i] would give read/write/delete per row, not per column
			fmt.Printf("%s=%t ", permissions[j], mat[i][j])
		}
		fmt.Println()
	}

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1. Arrays are fixed size — size is part of the type
	// 2. [3]int and [4]int are different types — not interchangeable
	// 3. Zero value = zero value of element type (0, "", false)
	// 4. Out of bounds with literal index → compile error (safe)
	// 5. Out of bounds with dynamic index → runtime panic (dangerous)
	//    Guard with: if i < len(arr) { }
	// 6. Arrays are VALUE types — assignment copies the whole array
	// 7. Arrays live on the stack — slices/maps live on the heap
	// 8. In real Go code, use slices — arrays are rarely used directly
	// 9. Trailing comma required when closing } is on its own line
	// 10. Don't shadow built-in functions (copy, len, append, make...)
	// --------------------------------------------------------
}

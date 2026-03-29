// ============================================================
// TOPIC 9: Slices
// ============================================================
// Slices are Go's dynamic arrays — like JS arrays but with rules.
// Unlike Go arrays (fixed size), slices can grow and shrink.
//
// A slice has THREE things under the hood:
//   1. Pointer  → points to an underlying array in memory
//   2. Length   → how many elements are currently in the slice
//   3. Capacity → how many elements the underlying array can hold
//
// This is why slices behave differently from arrays — they're
// a VIEW into an array, not the array itself.
// ============================================================

package main

import "fmt"

// --------------------------------------------------------
// Slice functions — unlike arrays, size is NOT part of the type
// A function taking []string works for ANY size slice
// This is the main reason you use slices over arrays in real code
// --------------------------------------------------------
func contains(items []string, target string) bool {
	for _, v := range items {
		if v == target {
			return true
		}
	}
	return false
}

func filter(nums []int, min int) []int {
	res := []int{}
	for _, v := range nums {
		if v > min {
			res = append(res, v)
		}
	}
	return res // ✅ always return res — even if empty
	// ❌ Don't return original slice as fallback — caller gets unexpected data
}

func use(middlewares []string, name string) []string {
	return append(middlewares, name)
	// ⚠️  Caller MUST reassign: middlewares = use(middlewares, "logger")
	// append may create a new underlying array if capacity exceeded
	// if you don't reassign, you're pointing at the old slice
}

func removeMiddleware(middlewares []string, name string) []string {
	result := []string{}
	for _, v := range middlewares {
		if v == name {
			continue
		}
		result = append(result, v)
	}
	return result
}

func main() {

	// --------------------------------------------------------
	// 1. DECLARING SLICES — three ways
	// --------------------------------------------------------
	// Way 1: Slice literal — like array but NO size in brackets
	// --------------------------------------------------------
	sliceLiteral := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	fmt.Println(sliceLiteral)
	fmt.Printf("len: %d, cap: %d\n", len(sliceLiteral), cap(sliceLiteral))

	// Way 2: make([]type, length, capacity)
	// ✅ Use make when you have a rough idea of the size upfront
	// This avoids repeated memory allocations as slice grows
	// ⚠️  make([]string, 5) sets BOTH len and cap to 5
	//     meaning 5 empty "" slots exist before you append
	//     appending then gives: ["" "" "" "" "" "GET" "POST"...]
	// ✅  make([]string, 0, 5) = length 0, capacity 5 — correct way
	sliceWithMake := make([]string, 0, 5)
	sliceWithMake = append(sliceWithMake, "GET", "POST", "PUT", "DELETE", "PATCH")
	fmt.Println(sliceWithMake)
	fmt.Printf("len: %d, cap: %d\n", len(sliceWithMake), cap(sliceWithMake))

	// Way 3: Declare empty, append later
	emptySlice := []string{}
	emptySlice = append(emptySlice, "GET", "POST", "PUT", "DELETE", "PATCH")
	fmt.Println(emptySlice)
	fmt.Printf("len: %d, cap: %d\n", len(emptySlice), cap(emptySlice))

	// --------------------------------------------------------
	// ZERO VALUE of a slice — nil, not []
	// --------------------------------------------------------
	// var s []string → nil slice (no underlying array)
	// s := []string{} → empty slice (has underlying array, len 0)
	//
	// Both print as [] but they are different:
	// nil slice  → s == nil is TRUE
	// empty slice → s == nil is FALSE
	//
	// Both can be appended to — behave the same for most operations
	// --------------------------------------------------------
	var nilSlice []string
	emptyS := []string{}
	fmt.Println(nilSlice == nil) // true
	fmt.Println(emptyS == nil)   // false
	fmt.Println(nilSlice)        // []
	fmt.Println(emptyS)          // []  — prints same but different internally

	// --------------------------------------------------------
	// 2. APPEND and CAPACITY GROWTH
	// --------------------------------------------------------
	// When you append beyond capacity, Go:
	//   1. Allocates a NEW bigger underlying array
	//   2. Copies all existing elements into it
	//   3. Discards the old array
	//   4. Returns a new slice pointing to the new array
	//
	// Growth rule: capacity roughly DOUBLES for small slices
	//   cap 2 → exceeded → new cap 4
	//   cap 4 → exceeded → new cap 8
	//   (Go uses a more complex formula for large slices)
	//
	// Performance implication in backends:
	//   Appending 10,000 items one by one = many allocations + copies
	//   Fix: use make([]T, 0, knownSize) if you know size upfront
	//   This pre-allocates and avoids repeated copying
	// --------------------------------------------------------
	s := make([]int, 0, 2)
	fmt.Printf("len: %d, cap: %d\n", len(s), cap(s)) // 0, 2
	s = append(s, 1)
	fmt.Printf("len: %d, cap: %d\n", len(s), cap(s)) // 1, 2
	s = append(s, 2)
	fmt.Printf("len: %d, cap: %d\n", len(s), cap(s)) // 2, 2
	s = append(s, 3)                                 // exceeds cap — new array allocated, cap doubles
	fmt.Printf("len: %d, cap: %d\n", len(s), cap(s)) // 3, 4
	s = append(s, 4)
	fmt.Printf("len: %d, cap: %d\n", len(s), cap(s)) // 4, 4
	s = append(s, 5)                                 // exceeds cap again — doubles to 8
	fmt.Printf("len: %d, cap: %d\n", len(s), cap(s)) // 5, 8

	// --------------------------------------------------------
	// 3. SLICES ARE REFERENCE TYPES ⚠️
	// --------------------------------------------------------
	// A slice variable holds a POINTER to the underlying array.
	// Assigning a slice to another variable copies the pointer,
	// not the data — both variables point to the same array.
	//
	// JS behaviour: const copy = original → same reference, mutates original
	// Go behaviour: ref := original      → same reference, mutates original
	// (opposite of arrays which are value types and get fully copied)
	// --------------------------------------------------------
	original := []string{"admin", "user", "guest"}
	ref := original // ref points to the SAME underlying array
	ref[0] = "superadmin"
	fmt.Println("original:", original) // [superadmin user guest] — mutated ⚠️
	fmt.Println("ref:", ref)

	// ✅ Fix: use copy() to create an independent slice
	// copy(dst, src) copies min(len(dst), len(src)) elements
	// dst MUST be pre-allocated with make — copy doesn't allocate
	// If dst has len 0, ZERO elements get copied (common mistake)
	copied := make([]string, len(original)) // pre-allocate with same length
	copy(copied, original)                  // copy(destination, source)
	copied[0] = "superadmin"
	fmt.Println("original:", original) // [admin user guest] — unchanged ✅
	fmt.Println("copied:", copied)

	// --------------------------------------------------------
	// 4. SLICING A SLICE — the head pointer concept
	// --------------------------------------------------------
	// Syntax: slice[low:high] — low inclusive, high exclusive
	// slice[0:3] → elements at index 0, 1, 2
	// slice[3:]  → from index 3 to end
	// slice[:3]  → from start to index 2
	//
	// CAPACITY after slicing — this surprises people:
	// The "head pointer" shifts to where the slice starts.
	// Capacity = original cap - start index
	// Because Go counts how many elements are visible FROM that head
	// all the way to the end of the underlying array.
	//
	// requests = [req1, req2, req3, req4, req5, req6]  cap=6
	// [0:3] → head at 0, sees 6 slots ahead → cap=6, len=3
	// [3:]  → head at 3, sees 3 slots ahead → cap=3, len=3
	// [1:4] → head at 1, sees 5 slots ahead → cap=5, len=3
	// --------------------------------------------------------
	requests := []string{"req1", "req2", "req3", "req4", "req5", "req6"}

	batch1 := requests[0:3]
	fmt.Println(batch1)
	fmt.Printf("len: %d, cap: %d\n", len(batch1), cap(batch1)) // 3, 6

	batch2 := requests[3:]
	fmt.Println(batch2)
	fmt.Printf("len: %d, cap: %d\n", len(batch2), cap(batch2)) // 3, 3

	middle := requests[1:4]
	fmt.Println(middle)
	fmt.Printf("len: %d, cap: %d\n", len(middle), cap(middle)) // 3, 5

	// --------------------------------------------------------
	// 5. PASSING SLICES TO FUNCTIONS
	// --------------------------------------------------------
	allowedRoutes := []string{"/health", "/users", "/posts", "/auth"}
	fmt.Println(contains(allowedRoutes, "/users"))   // true
	fmt.Println(contains(allowedRoutes, "/unknown")) // false

	nums := []int{45, 12, 89, 3, 67, 23, 91}
	fmt.Println(filter(nums, 50)) // [89 67 91]
	fmt.Println(filter(nums, 1))  // [45 12 89 3 67 23 91] — all above 1
	fmt.Println(filter(nums, 99)) // [] — empty slice, not original

	// --------------------------------------------------------
	// 6. REAL WORLD — middleware chain (like Gin internally)
	// --------------------------------------------------------
	// ⚠️  Must reassign when using functions that return slices
	// middlewares = use(middlewares, "logger") ✅
	// use(middlewares, "logger") alone does nothing to middlewares ❌
	// --------------------------------------------------------
	middlewares := []string{}
	middlewares = use(middlewares, "logger")
	middlewares = use(middlewares, "auth")
	middlewares = use(middlewares, "rateLimit")
	middlewares = use(middlewares, "cors")
	fmt.Println(middlewares)

	middlewares = removeMiddleware(middlewares, "auth")
	fmt.Println(middlewares) // [logger rateLimit cors]

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1.  Slice = pointer + length + capacity (3 things under the hood)
	// 2.  var s []string → nil slice (nil == true)
	//     s := []string{} → empty slice (nil == false)
	//     Both print [] but are internally different
	// 3.  make([]T, 0, cap) → length 0, capacity cap ✅
	//     make([]T, n)      → length n AND capacity n (pre-filled with zero values)
	// 4.  Slices are REFERENCE types — assignment copies the pointer not the data
	// 5.  Use copy(dst, src) to make independent copy — dst must be pre-allocated
	// 6.  copy copies min(len(dst), len(src)) elements
	// 7.  append returns a new slice — ALWAYS reassign: s = append(s, val)
	// 8.  When capacity exceeded: new array allocated, elements copied, old discarded
	// 9.  Capacity doubles for small slices — use make with known size for performance
	// 10. slice[low:high] — low inclusive, high exclusive
	// 11. Capacity after slicing = original cap - start index (head pointer shifts)
	// 12. Concurrent access to slices is NOT safe — use mutex (Topic 25)
	// --------------------------------------------------------
}

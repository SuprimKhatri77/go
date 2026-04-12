// ============================================================
// TOPIC 15: Generics
// ============================================================
// Generics let you write functions and structs that work for
// multiple types without repeating code or losing type safety.
//
// Before generics — two bad options:
//   1. Write the same function for each type (redundant)
//   2. Use interface{} and lose type safety (runtime panics)
//
// With generics — one function, type safe, compiler checked.
//
// Go 1.18+ only.
// ============================================================

package main

import (
	"encoding/json"
	"fmt"
)

// ============================================================
// 1. BASIC GENERIC FUNCTION
// ============================================================
// [T int | float64] means T can be either int or float64.
// Caller decides which — Go enforces it at compile time.
// ============================================================
func sum[T int | float64](nums []T) T {
	var total T // var total T = zero value of T — works for any type
	// can't write T{} — int, string, bool don't use {} syntax
	// var name T is the idiomatic way to get generic zero value
	for _, v := range nums {
		total += v
	}
	return total
}

// ============================================================
// 2. TYPE CONSTRAINTS — reusable interface for allowed types
// ============================================================
// Instead of listing types inline, define a constraint interface.
// This interface can only be used as a generic constraint —
// you can't use Number as a regular function parameter type.
// ============================================================
type Number interface {
	int | int32 | int64 | float32 | float64
}

func sumConstrained[T Number](nums []T) T {
	var total T
	for _, v := range nums {
		total += v
	}
	return total
}

func min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// ============================================================
// 3. GENERIC STRUCTS
// ============================================================
// [T any] on a struct means the struct works with any type.
// T is a placeholder — caller decides what it becomes.
//
// type Stack[T any] struct { items []T }
//
//	Stack[int]    → items is []int
//	Stack[string] → items is []string
//	Stack[Todo]   → items is []Todo
//
// Methods on generic structs use (s *Stack[T]) — T is the same
// T the struct was created with.
//
// Stack = LIFO (Last In First Out)
//
//	Push adds to the END
//	Pop  removes from the END
//	Peek reads  from the END without removing
//
// ============================================================
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(val T) {
	s.items = append(s.items, val)
}

func (s *Stack[T]) Pop() (T, error) {
	var zero T // zero value of T — idiomatic pattern for generic zero value
	if len(s.items) == 0 {
		return zero, fmt.Errorf("stack is empty")
		// use len() == 0 not == nil
		// an emptied slice is [] not nil — nil check would miss it
	}
	top := s.items[len(s.items)-1]     // last element
	s.items = s.items[:len(s.items)-1] // remove last
	return top, nil
}

func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if len(s.items) == 0 {
		return zero, fmt.Errorf("stack is empty")
	}
	return s.items[len(s.items)-1], nil // read last, don't remove
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

// ============================================================
// 4. GENERIC UTILITY FUNCTIONS
// ============================================================
//
// FUNCTIONS AS PARAMETERS — fn func(T) bool
// When you see fn as a parameter, you provide the implementation
// at the call site. Two ways:
//
//	// way 1 — inline anonymous function (most common)
//	filter(todos, func(t Todo) bool {
//	    return t.Completed // your logic here
//	})
//
//	// way 2 — named function defined elsewhere
//	func isCompleted(t Todo) bool { return t.Completed }
//	filter(todos, isCompleted) // pass by name
//
// The signature in the parameter just defines the shape.
// You fill in the body when calling.
//
// comparable CONSTRAINT:
//
//	Means T supports == and != operators.
//	Nothing to do with return type — only about what you can DO with T.
//	Structs are comparable only if ALL their fields are comparable.
//	Slices, maps, functions are NOT comparable.
//
// ============================================================
func contains[T comparable](items []T, target T) bool {
	for _, v := range items {
		if v == target { // == only works because T is constrained to comparable
			return true
		}
	}
	return false
}

func filter[T any](items []T, fn func(T) bool) []T {
	result := []T{} // initialize empty slice — never return nil
	// returning nil means caller must nil-check before ranging
	// returning []T{} means caller can always range safely
	for _, v := range items {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Two type parameters — T is input type, R is output/result type
// T and R can be completely different types
// eg: []Todo → []string (T=Todo, R=string)
func mapSlice[T any, R any](items []T, fn func(T) R) []R {
	result := []R{}
	for _, v := range items {
		result = append(result, fn(v))
	}
	return result
}

// ============================================================
// 5. GENERIC STRUCTS — APIResponse
// ============================================================
// Data T means the Data field type is decided by the caller.
// Type safe — no interface{}, no type assertions, no runtime panics.
//
// Compare:
//
//	interface{} version:
//	  res.Data           → type is interface{}
//	  res.Data.(*Todo)   → type assertion needed, panics if wrong
//
//	generic version:
//	  res.Data           → type is Todo, compiler knows
//	  res.Data.Title     → works directly, no assertion needed
//
// ============================================================
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

// respond should only BUILD the response — not marshal inside
// marshaling is a side effect — let the caller decide what to do
func respond[T any](success bool, message string, data T) APIResponse[T] {
	return APIResponse[T]{
		Success: success,
		Message: message,
		Data:    data,
	}
}

type Todo struct {
	Title     string
	Completed bool
}

type User struct {
	ID   int
	Name string
}

func main() {

	// basic generic function
	fmt.Println(sum([]int{1, 2, 3, 4, 5}))     // 15
	fmt.Println(sum([]float64{1.1, 2.2, 3.3})) // 6.6

	// type constraint
	fmt.Println(sumConstrained([]int32{22, 34, 45})) // 101
	fmt.Println(sumConstrained([]int64{100, 200}))   // 300
	fmt.Println(min(5, 7))                           // 5
	fmt.Println(max(5, 7))                           // 7

	// generic stack — LIFO
	stack := &Stack[int]{}
	stack.Push(1)
	stack.Push(2)
	stack.Push(3) // top of stack

	top, _ := stack.Pop() // removes 3
	fmt.Println(top)      // 3

	peek, _ := stack.Peek()   // reads 2, doesn't remove
	fmt.Println(peek)         // 2
	fmt.Println(stack.Size()) // 2 — 3 was popped, 1 and 2 remain

	// same struct, different type
	middlewares := &Stack[string]{}
	middlewares.Push("logger")
	middlewares.Push("auth")
	middlewares.Push("rateLimit") // top
	top2, _ := middlewares.Pop()
	fmt.Println(top2) // rateLimit

	// utility functions
	todos := []Todo{
		{Title: "first", Completed: false},
		{Title: "second", Completed: true},
		{Title: "third", Completed: true},
		{Title: "fourth", Completed: false},
	}

	// contains — comparable constraint allows ==
	fmt.Println(contains([]string{"admin", "user"}, "admin")) // true
	fmt.Println(contains([]int{1, 2, 3}, 5))                  // false
	fmt.Println(contains(todos, todos[0]))                    // true — Todo is comparable (string + bool fields only)

	// filter — fn provided inline at call site
	completed := filter(todos, func(t Todo) bool {
		return t.Completed // your logic here
	})
	fmt.Println(completed) // [{second true} {third true}]

	// mapSlice — T=Todo, R=string
	titles := mapSlice(todos, func(t Todo) string {
		return t.Title // extract one field, returns []string
	})
	fmt.Println(titles) // [first second third fourth]

	// generic APIResponse — type safe
	todoRes := respond(true, "todo created", Todo{Title: "learn go", Completed: false})
	j, _ := json.Marshal(todoRes)
	fmt.Println(string(j))
	// todoRes.Data is Todo — no type assertion needed
	fmt.Println(todoRes.Data.Title) // "learn go" — direct access, compiler knows the type

	userRes := respond(true, "users fetched", []User{{ID: 1, Name: "alice"}})
	j, _ = json.Marshal(userRes)
	fmt.Println(string(j))

	strRes := respond(true, "message", "hello")
	j, _ = json.Marshal(strRes)
	fmt.Println(string(j))

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1.  [T any] — T can be any type, caller decides at usage
	// 2.  [T Number] — T constrained to types in Number interface
	// 3.  Type constraints are interfaces with | for multiple types
	//     only usable as constraints, not as regular param types
	// 4.  var zero T — idiomatic zero value for generic types
	//     can't use T{} — doesn't work for int, string, bool
	// 5.  Generic structs: type Foo[T any] struct { val T }
	//     Foo[int]{} → val is int, Foo[string]{} → val is string
	// 6.  Methods on generic structs: func (s *Stack[T]) Push(val T)
	//     T is the same T the struct was created with
	// 7.  comparable — supports == and !=, nothing to do with return type
	//     structs comparable only if all fields are comparable
	//     slices, maps, functions are NOT comparable
	// 8.  fn func(T) bool as param — provide implementation at call site
	//     inline: filter(items, func(t T) bool { return ... })
	//     named:  filter(items, myFunc) — myFunc must match signature
	// 9.  Two type params: [T any, R any] — input and output can differ
	//     mapSlice([]Todo, fn) → []string (T=Todo, R=string)
	// 10. Generic response vs interface{}: generics give compile-time
	//     type safety — res.Data.Field works directly, no assertion
	// 11. Stack is LIFO — push/pop/peek from the END (last element)
	//     use len(s.items)-1, not index 0 (that's a queue)
	//     use len() == 0 not == nil for empty check
	// --------------------------------------------------------
}

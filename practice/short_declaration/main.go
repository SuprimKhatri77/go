// ============================================================
// TOPIC 2: Short Variable Declaration (:=)
// ============================================================
// := is Go's shorthand for declaring AND initializing a variable
// in one step. No `var` keyword, no explicit type needed.
// Go infers the type from the value you assign.
//
// Think of it as the everyday variable declaration inside functions.
// You'll use this 90% of the time over `var`.
//
// JS/TS equivalent: const x = 10 (but mutable, like let)
// ============================================================

package main

import "fmt"

// ❌ ILLEGAL: := cannot be used at package level
// host := "localhost"  // compile error — := is a statement, not a declaration
//
// WHY: Package level only allows declarations (var, const, func, type)
// := is a statement — statements only live inside function bodies
// Use `var` at package level instead:
var host = "localhost" // ✅ var with inferred type at package level

func main() {

	// --------------------------------------------------------
	// 1. BASIC SHORT DECLARATION
	// --------------------------------------------------------
	// Syntax: <name> := <value>
	// Type is inferred automatically from the value.
	// Must be inside a function — always.
	// Must include a value — cannot just declare without initializing.
	// --------------------------------------------------------
	port := 9000     // inferred as int
	isSecure := true // inferred as bool
	version := "v2"  // inferred as string
	gpa := 9.0       // inferred as float64 (default for decimals)
	score := 9       // inferred as int (NOT int32 — Go's int is platform-sized, 64-bit on modern systems)

	fmt.Println(host, port, isSecure, version)
	fmt.Printf("score type: %T\n", score) // int
	fmt.Printf("gpa type: %T\n", gpa)     // float64

	// --------------------------------------------------------
	// 2. TYPE INFERENCE TRAP
	// --------------------------------------------------------
	// := infers from the literal value.
	// 9   → int
	// 9.0 → float64
	// These are DIFFERENT types — you cannot add them directly.
	// Must explicitly convert one side.
	// Safer direction: convert int → float64 (preserves decimals)
	// --------------------------------------------------------
	total := float64(score) + gpa // ✅ explicit conversion
	fmt.Println("total:", total)

	// --------------------------------------------------------
	// 3. THE REDECLARATION RULE — the most important := quirk
	// --------------------------------------------------------
	// You CANNOT redeclare a variable with := in the same scope
	// UNLESS at least one variable on the left side is brand new.
	//
	// When an existing variable appears on the left:
	//   - It is REASSIGNED (not redeclared)
	//   - Same variable, same memory address, new value
	//   - Type must stay the same
	//
	// This pattern is used constantly with `err` in real Go code:
	// --------------------------------------------------------
	userID, err := 101, error(nil) // userID: NEW, err: NEW
	username, err := "claude", nil // username: NEW, err: REASSIGNED ✅
	fmt.Println(userID, username, err)

	// ❌ This would NOT compile — nothing new on the left:
	// userID, err := 999, nil  // both already declared in this scope

	// ✅ Use = (assignment) when no new variables involved:
	// userID = 999  // just reassign

	// --------------------------------------------------------
	// 4. SCOPE TRAP — := respects block scope strictly
	// --------------------------------------------------------
	// Variables declared with := inside a block (if, for, {})
	// are NOT accessible outside that block.
	// This is the same as let/const in JS — no hoisting, strict scope.
	// --------------------------------------------------------
	status := "inactive"

	if status == "inactive" {
		updatedStatus := "active" // scoped to this if block only
		fmt.Println("inside block:", updatedStatus)
	}
	// fmt.Println(updatedStatus) // ❌ compile error — not in scope

	// ✅ Fix: declare outside the block first
	var updatedStatus string // zero value: ""
	if status == "inactive" {
		updatedStatus = "active" // assignment, not declaration
	}
	fmt.Println("outside block:", updatedStatus)

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1. := declares AND initializes — you cannot use it without a value
	// 2. := only works inside functions (it's a statement, not a declaration)
	// 3. At least one variable on the left must be new for := to be valid
	// 4. Existing variables on the left are reassigned, not redeclared
	// 5. Type is inferred once and cannot change — statically typed
	// 6. `var` is for package level or zero-value declarations
	// 7. Use `=` (not :=) when all variables on the left already exist
	// --------------------------------------------------------

	// WHEN TO USE WHICH:
	// := inside functions, when you have a value ready           (most common)
	// var inside functions, when you need a zero value first     (scope fix pattern)
	// var at package level, for config/globals                   (no := allowed here)
	// const for values that never change                         (compile-time only)
}

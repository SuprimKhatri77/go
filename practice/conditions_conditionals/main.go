// ============================================================
// TOPIC 5: Conditions and Conditionals
// ============================================================
// Similar concept to JS but with key differences:
// 1. No parentheses around conditions — gofmt removes them
// 2. Braces {} are mandatory — no single line if statements
// 3. Init statement in if — doesn't exist in JS
// 4. Early return pattern — preferred over nested else blocks
// ============================================================

package main

import "fmt"

// --------------------------------------------------------
// EARLY RETURN PATTERN — preferred in Go
// --------------------------------------------------------
// Instead of if/else chains, return early and let the rest
// of the function be the "happy path". Reduces nesting.
// You'll see this everywhere in real Go codebases.
// --------------------------------------------------------
func getRole(userID int) string {
	if userID == 1 {
		return "admin"
	}
	return "user" // only runs if the if check above failed — no else needed
}

func getUser(id int) (string, error) {
	if id == 1 {
		return "Alice", nil
	}
	return "", fmt.Errorf("user not found") // ✅ error messages are lowercase in Go
	// ❌ fmt.Errorf("User not found") — capital looks odd when errors are chained
	// eg: "failed to fetch: User not found" vs "failed to fetch: user not found"
}

func main() {

	// --------------------------------------------------------
	// 1. BASIC if/else — what's different from JS
	// --------------------------------------------------------
	// ❌ No parentheses around condition (gofmt removes them)
	// ❌ No single-line if without braces
	// ✅ Braces {} are mandatory always
	// ✅ Condition order matters — check most restrictive first
	// --------------------------------------------------------
	token := "abc123"
	isExpired := false
	role := "admin"

	// ✅ Token check first — no point checking expiry/role if no token
	if token == "" {
		fmt.Println("access denied")
	} else if isExpired {
		fmt.Println("token expired")
	} else if role == "admin" {
		fmt.Println("welcome admin")
	} else {
		fmt.Println("welcome user")
	}

	// ❌ This won't compile — no braces, no single-line if
	// if token == "" fmt.Println("denied")

	// ❌ This won't compile — parentheses alone don't create a scope
	// if (token == "") fmt.Println("denied")

	// ✅ Parentheses around condition are allowed but gofmt removes them
	// if (token == "") { } → gofmt formats to → if token == "" { }

	// --------------------------------------------------------
	// 2. INIT STATEMENT — Go's unique pattern, not in JS
	// --------------------------------------------------------
	// Syntax: if <init>; <condition> { }
	// - <init> runs first, declares a variable
	// - that variable is scoped to the entire if/else if/else block
	// - cannot be accessed outside the block
	//
	// This is the standard pattern for function calls that return
	// a value + error — keeps variables tightly scoped
	// --------------------------------------------------------

	// Without init statement — role leaks into the rest of main
	// role2 := getRole(1)
	// if role2 == "admin" { ... }

	// ✅ With init statement — role scoped to this block only
	if r := getRole(1); r == "admin" {
		fmt.Println("admin access")
	} else {
		fmt.Println("role:", r) // r is accessible in else too
	}
	// fmt.Println(r) // ❌ compile error — r not in scope here

	// --------------------------------------------------------
	// 3. VARIABLE SHADOWING with init statement ⚠️
	// --------------------------------------------------------
	// If an outer variable has the same name as the init variable,
	// the inner one SHADOWS the outer one within the block.
	// They are two completely separate variables.
	// Some teams use a `shadow` linter to catch this bug.
	// --------------------------------------------------------
	outerRole := "user"
	if outerRole := getRole(1); outerRole == "admin" { // different variable, same name
		fmt.Println("inner:", outerRole) // "admin" — inner variable
	}
	fmt.Println("outer:", outerRole) // "user" — outer unchanged

	// --------------------------------------------------------
	// 4. LOGICAL OPERATORS — same as JS
	// --------------------------------------------------------
	// &&  AND
	// ||  OR
	// !   NOT
	// No difference from JS — but no extra () wrapping the whole condition
	// Use parentheses to group sub-conditions for clarity
	// --------------------------------------------------------
	requestsThisMinute := 81
	isPremiumUser := false
	isInternalIP := false

	if (requestsThisMinute > 100) || (requestsThisMinute > 80 && !isPremiumUser && !isInternalIP) {
		fmt.Println("rate limited")
	} else {
		fmt.Println("allowed")
	}

	// --------------------------------------------------------
	// 5. INIT STATEMENT vs normal if — scope comparison
	// --------------------------------------------------------
	// Normal if — u and err live in the rest of main's scope
	// Risk: stray err gets reassigned later, you check the wrong error
	// --------------------------------------------------------
	u, err := getUser(7)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u)
	}
	// u and err still exist here — risk of accidental reuse

	// ✅ Init statement — u and err scoped to this block only
	// Preferred when you don't need the values outside the condition
	if u, err := getUser(1); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u) // "Alice"
	}
	// u and err from above are gone here — clean scope

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1. No parentheses around if conditions — gofmt removes them
	// 2. Braces {} are mandatory — no single-line ifs
	// 3. Init statement: if val := fn(); val == x { } — val scoped to block
	// 4. Init statement variable is accessible in else/else if too
	// 5. Early return > nested else — flatter code is more readable
	// 6. Error messages are lowercase — they get chained/wrapped
	// 7. Variable shadowing in init statement — same name ≠ same variable
	// 8. Stray err in wide scope is a common bug — use init statement to avoid it
	// --------------------------------------------------------
}

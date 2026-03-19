// ============================================================
// TOPIC 1: Data Types & Variables
// ============================================================
// Go is statically typed — types are checked at COMPILE time.
// Unlike JS, there is no runtime type coercion. No undefined.
// No null (for basic types). Every variable has a zero value.
// ============================================================

package main

import "fmt"

// -- CONSTANTS --
// Must be resolvable at compile time.
// Only allowed types: string, bool, int, float, rune.
// Slices, maps, structs CANNOT be constants (they're runtime values).
// Think: `const` in JS, but stricter.
const MAX_LOGIN_ATTEMPTS int = 5
const API_VERSION string = "v1"
const PI float64 = 3.14159

func main() {

	// --------------------------------------------------------
	// 1. EXPLICIT VARIABLE DECLARATION (using `var`)
	// --------------------------------------------------------
	// Syntax: var <name> <type> = <value>
	// Use this when you want to be explicit about the type,
	// or when declaring without an initial value (zero value).
	// Similar to: let port: number = 8080 in TypeScript
	// --------------------------------------------------------
	var serverHost string = "localhost"
	var port int = 8080
	var maxRequestBodySize int64 = 1048576 // Use int64 for sizes/offsets — int32 overflows at ~2GB
	var isHTTPS bool = false
	var requestTimeout float64 = 30.5 // Always prefer float64 over float32 (more precision, no cost)

	fmt.Println(serverHost)
	fmt.Println(port)
	fmt.Println(maxRequestBodySize)
	fmt.Println(isHTTPS)
	fmt.Println(requestTimeout)

	// --------------------------------------------------------
	// 2. ZERO VALUES — Go's answer to undefined/null
	// --------------------------------------------------------
	// Every type has a default "zero value" when declared without assignment.
	// string  → ""
	// int     → 0
	// bool    → false
	// float64 → 0
	// No surprises like JS's undefined or null. Always predictable.
	// --------------------------------------------------------
	var apiKey string        // ""
	var requestCount int     // 0
	var isAuthenticated bool // false
	var responseTime float64 // 0

	fmt.Println(apiKey)
	fmt.Println(requestCount)
	fmt.Println(isAuthenticated)
	fmt.Println(responseTime)

	// --------------------------------------------------------
	// 3. TYPE CONVERSION — Go never implicitly converts types
	// --------------------------------------------------------
	// Unlike JS, Go will NOT auto-convert int + float64.
	// You must explicitly convert: float64(x) or int(x).
	// CAUTION: int(9.5) → 9. Decimals are silently truncated!
	// Safer direction: convert int → float64 to preserve precision.
	// --------------------------------------------------------
	var userID int = 42
	var sessionScore float64 = 9.5

	// WRONG approach (loses decimal): int(sessionScore) → 9
	// var sum = userID + int(sessionScore)

	// CORRECT approach (preserves decimal):
	var sum float64 = float64(userID) + sessionScore
	fmt.Println("sum:", sum) // 51.5

	// --------------------------------------------------------
	// 4. VAR BLOCK — grouping related variables
	// --------------------------------------------------------
	// Use var() blocks for related config/setup values.
	// Cleaner than repeating `var` on every line.
	// Similar to declaring multiple consts/lets at the top of a JS module.
	// --------------------------------------------------------
	var (
		firstName, lastName string = "John", "Doe" // multiple vars, same type, one line
		age                 int    = 28
		isAdmin             bool   = true
	)

	fmt.Println("FirstName:", firstName)
	fmt.Println("LastName:", lastName)
	fmt.Println("Age:", age)
	fmt.Println("IsAdmin:", isAdmin)

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1. Every declared variable MUST be used — Go won't compile otherwise.
	//    (Like ESLint errors, but enforced by the compiler. No warnings, just failure.)
	// 2. Same for imports — unused imports = compile error.
	// 3. Variables cannot be redeclared in the same scope (like `let` in JS).
	// 4. No hoisting — variables are scoped to the block they're declared in.
	// 5. Type cannot change after declaration (statically typed).
	// --------------------------------------------------------
}

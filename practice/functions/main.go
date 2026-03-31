// ============================================================
// TOPIC 11: Functions
// ============================================================
// Functions in Go are first-class citizens — they can be:
//   - Assigned to variables
//   - Passed as arguments to other functions
//   - Returned from functions
//
// This is the foundation of middleware patterns, closures,
// dependency injection — everything you do in Gin is built on this.
// ============================================================

package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ============================================================
// 1. BASIC FUNCTION — single return
// ============================================================
// Syntax: func name(params) returnType { }
// ============================================================
func greetHandler(name string) string {
	return "Hello, " + name
}

// ============================================================
// 2. MULTIPLE RETURN VALUES
// ============================================================
// Go's answer to try/catch — return result AND error together
// Caller must handle both values
// Convention: error is always the LAST return value
//
// Syntax: func name(params) (type1, type2) { }
//
//	parentheses required for multiple returns
//
// ============================================================
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("b cannot be zero") // error messages lowercase
	}
	return a / b, nil
}

func parsePort(s string) (int, error) {
	num, err := strconv.Atoi(s) // strconv.Atoi = string to int, returns (int, error)
	if err != nil {
		return 0, fmt.Errorf("invalid port number: %s", s)
	}
	if num < 1 || num > 65535 {
		return 0, fmt.Errorf("port must be between 1 and 65535, got %d", num)
	}
	return num, nil
}

// ============================================================
// 3. NAMED RETURN VALUES
// ============================================================
// Name your return values — they become variables in the function.
// Use naked `return` to return them all at once.
//
// Syntax: func name(params) (returnName1 type, returnName2 type) { }
//
//	parentheses required, names are optional but useful
//
// WHEN TO USE:
//
//	✅ Short functions — documents what each return value is
//	✅ Defer + named returns (you'll see this in Topic 27)
//	❌ Long functions — naked return makes it hard to track values
//
// ============================================================
func getServerInfo() (host string, port int, err error) {
	host = "localhost"
	port = 8080
	err = nil
	return // naked return — returns host, port, err
}

func calcPagination(totalItems, pageSize, currentPage int) (offset, limit, totalPages int) {
	// ✅ Always guard against zero division
	// pageSize = 0 would panic without this check
	if pageSize == 0 {
		return 0, 0, 0
	}
	offset = pageSize * (currentPage - 1)
	limit = pageSize
	totalPages = totalItems / pageSize
	return // naked return
}

// ============================================================
// 4. VARIADIC FUNCTIONS
// ============================================================
// Accept any number of arguments of the same type.
// JS equivalent: function sum(...nums) { }
//
// Syntax: func name(param ...type) { }
// Inside the function, param is a SLICE — you can range over it.
//
// PASSING A SLICE to a variadic function:
//
//	nums := []int{1, 2, 3}
//	sum(nums...)   ✅ spread with ... — like JS spread operator
//	sum(nums)      ❌ type mismatch — []int is not int
//
// Note: arrays CANNOT be spread — only slices
//
// USE CASES:
//   - fmt.Println(a, b, c, d) — takes ...interface{}
//   - sum(1, 2, 3, 4, 5)
//   - logMessage("INFO", "server", "started", "on", "port", "8080")
//   - Anywhere the number of inputs is unknown at compile time
//
// ============================================================
func sum(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
	// sum()        → 0  (nothing to range over, total stays 0)
	// sum(1,2,3)   → 6
	// sum(nums...) → works if nums is []int
}

func logMessage(level string, parts ...string) string {
	// strings.Join avoids leading space from manual concatenation
	return fmt.Sprintf("[%s] %s", level, strings.Join(parts, " "))
}

// ============================================================
// 5. FUNCTIONS AS VALUES — the key mental model
// ============================================================
// A function type is just like any other type.
// func(string) string means: "a variable holding a function
// that takes a string and returns a string"
//
// Reading function signatures — break it down:
//   func(string) string
//     └─ takes string, returns string
//
//   func(func(string) string) func(string) string
//     └─ takes a (func that takes string returns string)
//     └─ returns a (func that takes string returns string)
//     └─ this is a MIDDLEWARE WRAPPER signature
//
// HOW TO IDENTIFY THE PATTERN:
//   1. Read left to right
//   2. Each func(...) ... is one function type
//   3. If you see func inside func — it's a higher order function
//   4. The outermost func is the signature of the function itself
// ============================================================

// logger is a MIDDLEWARE — wraps a handler and adds behavior
// Pattern: takes a handler, returns a NEW handler with extra behavior
// This is exactly how Gin middleware works internally
func logger(next func(string) string) func(string) string {
	// returns a NEW function that:
	// 1. does something (logging)
	// 2. calls the original handler (next)
	// 3. returns the result
	return func(s string) string {
		fmt.Println("[LOG] calling handler")
		result := next(s) // call the original handler
		fmt.Println("[LOG] handler returned:", result)
		return result
	}
}

// applyMiddleware connects a handler and middleware together
// Takes:  handler    → func(string) string
//
//	middleware → func(func(string)string) func(string)string
//
// Returns: the wrapped handler → func(string) string
//
// How it works:
//
//	applyMiddleware(greetHandler, logger)
//	→ logger(greetHandler)
//	→ returns the wrapper function
//	→ when called: logs, calls greetHandler, logs result
func applyMiddleware(
	handler func(string) string,
	middleware func(func(string) string) func(string) string,
) func(string) string {
	return middleware(handler)
}

// ============================================================
// 6. FUNCTIONS RETURNING FUNCTIONS — factory pattern
// ============================================================
// A function that returns a function — used for:
//   - Rate limiters
//   - Dependency injection (Gin handlers with DB)
//   - Closures (Topic 26)
//
// The returned function REMEMBERS variables from the outer function.
// This is called a CLOSURE — covered fully in Topic 26.
// For now: think of it as the returned function having a backpack
// that carries the outer function's variables with it.
// ============================================================
func makeRateLimiter(maxRequests int) func() bool {
	requestCount := 0 // this variable lives in the "backpack"
	return func() bool {
		requestCount++
		if requestCount > maxRequests { // ✅ use the param, not hardcoded 3
			return false
		}
		return true
	}
}

// ============================================================
// 7. PIPELINE — chaining functions
// ============================================================
// Each function's output becomes the next function's input.
// This is exactly how Gin's middleware chain works:
//
//	request → logger → auth → rateLimit → handler → response
//
// pipeline returns a FUNCTION that when called:
//  1. Takes the initial input
//  2. Passes it through each middleware in order
//  3. Returns the final output
//
// ============================================================
func withRequestID(req string) string {
	return req + " [id:abc123]"
}

func withTimestamp(req string) string {
	return req + " [ts:1234567890]"
}

func withAuth(req string) string {
	return "AUTHENTICATED: " + req
}

func pipeline(middlewares ...func(string) string) func(string) string {
	return func(req string) string {
		result := req
		for _, m := range middlewares {
			result = m(result) // each middleware gets previous output
		}
		return result
	}
}

func main() {

	// Multiple return values
	res, err := divide(10, 3)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("result: %.2f\n", res)
	}

	port, err := parsePort("8080")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("port:", port)
	}

	// Named returns
	host, p, err := getServerInfo()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("host: %s, port: %d\n", host, p)

	offset, limit, totalPages := calcPagination(100, 10, 3)
	fmt.Printf("offset: %d, limit: %d, totalPages: %d\n", offset, limit, totalPages)

	// Variadic
	fmt.Println(sum())        // 0
	fmt.Println(sum(1, 2, 3)) // 6
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(sum(nums...)) // 15 — spread slice into variadic

	fmt.Println(logMessage("INFO", "server", "started", "on port", "8080"))

	// Functions as values — middleware pattern
	wrapped := applyMiddleware(greetHandler, logger)
	fmt.Println(wrapped("Alice"))
	// [LOG] calling handler
	// [LOG] handler returned: Hello, Alice
	// Hello, Alice

	// Factory pattern — rate limiter
	limiter := makeRateLimiter(3)
	fmt.Println(limiter()) // true
	fmt.Println(limiter()) // true
	fmt.Println(limiter()) // true
	fmt.Println(limiter()) // false — exceeded

	// Pipeline — chained middleware
	process := pipeline(withRequestID, withTimestamp, withAuth)
	fmt.Println(process("GET /users"))
	// AUTHENTICATED: GET /users [id:abc123] [ts:1234567890]

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1.  Multiple returns: func f() (int, error) — parens required
	// 2.  Error is always the LAST return value by convention
	// 3.  Named returns: func f() (a, b int) — naked return returns both
	//     Use for short functions only — confusing in long ones
	// 4.  Variadic: func f(nums ...int) — nums is a slice inside
	// 5.  Pass slice to variadic: f(slice...) — spread with ...
	//     Arrays cannot be spread — only slices
	// 6.  func(string) string is a TYPE — a variable holding a function
	// 7.  Middleware signature: func(func(string)string) func(string)string
	//     Takes a handler, returns a wrapped handler
	// 8.  Factory pattern: func returns func — inner func remembers outer vars
	//     This is a closure — covered fully in Topic 26
	// 9.  Pipeline: chain functions so each output feeds next input
	//     This is how Gin middleware chain works internally
	// 10. sum() with no args → returns zero value (0) — safe, no panic
	// --------------------------------------------------------
}

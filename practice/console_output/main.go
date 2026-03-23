// ============================================================
// TOPIC 3: Console Output (fmt package)
// ============================================================
// Go's fmt package is the equivalent of console.log in JS,
// but more explicit and more powerful.
//
// Three families:
//   Print/Println/Printf  → write to stdout
//   Sprintf               → returns a formatted string
//   Fprintf               → writes to any writer (stderr, files, HTTP)
//
// Import: "fmt" and "os" (for Fprintf with os.Stderr/os.Stdout)
// ============================================================

package main

import (
	"fmt"
	"os"
)

// --------------------------------------------------------
// Sprintf used inside a function — returns string, doesn't print
// --------------------------------------------------------
// ✅ Correct: just return the Sprintf result directly
// ❌ Avoid: naming the return variable `err` — misleading to readers
//
//	`err` in Go always signals an error type by convention
//
// --------------------------------------------------------
func formatLog(level string, message string) string {
	return fmt.Sprintf("[%s] %s", level, message) // S = String, returns formatted string
}

type User struct {
	Name  string
	Email string
	Age   int
}

func main() {

	// --------------------------------------------------------
	// 1. Print vs Println vs Printf
	// --------------------------------------------------------
	// Print   — no newline, no spaces between args, no format verbs
	// Println — adds newline + spaces between multiple args automatically
	// Printf  — no newline, supports format verbs (%s %d %f %t %v etc)
	//
	// JS equivalent:
	//   Println → console.log()
	//   Printf  → console.log(`formatted ${var}`)
	//   Print   → rarely needed
	// --------------------------------------------------------
	fmt.Print("hello")            // no newline
	fmt.Print("world")            // prints right after: helloworld
	fmt.Println()                 // just a newline
	fmt.Println("hello", "world") // hello world  (auto space between args)
	fmt.Print("hello", "world")   // helloworld   (no auto space)
	fmt.Printf("hello world\n")   // hello world  (manual \n needed)

	// ⚠️  Literal % in Printf must be escaped as %%
	fmt.Printf("success rate: 50%%\n") // prints: success rate: 50%

	// --------------------------------------------------------
	// 2. Format Verbs — used with Printf and Sprintf
	// --------------------------------------------------------
	// %s  → string
	// %d  → integer (decimal)
	// %f  → float   (%.2f = 2 decimal places)
	// %t  → boolean
	// %v  → any value (default format)
	// %+v → struct with field names
	// %#v → Go syntax representation (type + fields)
	// %T  → type of the variable
	// %%  → literal percent sign
	// --------------------------------------------------------
	method := "POST"
	path := "/api/users"
	status := 201
	duration := 12.45
	authenticated := true

	fmt.Printf("[REQUEST] %s %s | status: %d | duration: %.2fms | authenticated: %t\n",
		method, path, status, duration, authenticated)

	// ⚠️  Wrong verb = runtime formatting error, NOT a compile error
	// fmt.Printf("%d", "hello") → prints %!d(string=hello), code keeps running
	// Silent bug in logs — hard to catch. Use the right verb always.

	// --------------------------------------------------------
	// 3. Sprintf — returns a formatted string
	// --------------------------------------------------------
	// Use when you need to BUILD a string, not print it immediately.
	// Common uses: error messages, log entries, dynamic strings
	// S = String → Sprintf returns a string
	// --------------------------------------------------------
	logLine := formatLog("ERROR", "something went wrong")
	fmt.Println(logLine) // [ERROR] something went wrong

	fmt.Println(formatLog("INFO", "server started on port 8080"))

	// Sprintf inline — store result in a variable
	msg := fmt.Sprintf("user %s logged in from %s", "alice", "192.168.1.1")
	fmt.Println(msg)
	fmt.Printf("%T\n", msg) // string

	// --------------------------------------------------------
	// 4. Fprintf — write to any writer
	// --------------------------------------------------------
	// Like console.log vs console.error in JS — same terminal,
	// but different streams. Docker, systemd, log tools treat them differently.
	//
	// os.Stdout → normal output  (like console.log)
	// os.Stderr → error output   (like console.error)
	//
	// You'll use this in real backends for structured logging,
	// and Fprintf is also how Go writes HTTP responses (more on this later)
	// --------------------------------------------------------
	fmt.Fprintf(os.Stdout, "server started on port %d\n", 8080)
	fmt.Fprintf(os.Stderr, "error: %s\n", "database connection failed")

	// --------------------------------------------------------
	// 5. Debug verbs — %v, %+v, %#v
	// --------------------------------------------------------
	// Your best friends when debugging structs and unknown types
	//
	// %v   → values only              {Alice alice@example.com 30}
	// %+v  → field names + values     {Name:Alice Email:alice@example.com Age:30}
	// %#v  → Go syntax (type+fields)  main.User{Name:"Alice", Email:"alice@example.com", Age:30}
	// %T   → just the type            main.User
	//
	// JS equivalent:
	//   %v   → console.log(obj)           basic output
	//   %+v  → JSON.stringify(obj)        with keys
	//   %#v  → seeing the full type def   for runtime type debugging
	// --------------------------------------------------------
	u := User{"Alice", "alice@example.com", 30}
	fmt.Printf("%v\n", u)
	fmt.Printf("%+v\n", u)
	fmt.Printf("%#v\n", u)
	fmt.Printf("%T\n", u)

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1. Println adds newline + spaces between args automatically
	// 2. Printf and Print do NOT add a newline — you need \n
	// 3. Print does NOT support format verbs — use Printf for that
	// 4. Wrong format verb (%d for a string) is a RUNTIME error, not compile
	// 5. Sprintf returns a string — use it to build messages before printing
	// 6. Fprintf writes to a writer — os.Stdout, os.Stderr, files, HTTP
	// 7. Escape literal % as %% in Printf/Sprintf
	// --------------------------------------------------------

}

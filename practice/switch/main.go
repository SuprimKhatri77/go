// ============================================================
// TOPIC 6: Switch
// ============================================================
// Go's switch is more powerful than JS's switch.
// Key differences:
// 1. No break needed — Go breaks automatically after each case
// 2. Multiple values per case in one line
// 3. Expressionless (naked) switch — each case is a full condition
// 4. Init statement support — just like if statements
// 5. fallthrough is explicit and rare — opposite default from JS
//
// JS default: fallthrough (need break to stop)
// Go default: break     (need fallthrough to continue)
// ============================================================

package main

import "fmt"

// --------------------------------------------------------
// Expressionless switch used inside a function
// Cleaner than if/else chains for multiple conditions
// --------------------------------------------------------
func getUserRole(userID int) string {
	switch {
	case userID == 1:
		return "admin"
	case userID == 2:
		return "moderator"
	default:
		return "guest"
	}
}

func main() {

	// --------------------------------------------------------
	// 1. BASIC SWITCH — no break needed
	// --------------------------------------------------------
	// Go automatically stops at the end of a matched case.
	// default runs when NO other case matches — like the final
	// else in an if/else chain. Always runs if nothing matches.
	// --------------------------------------------------------
	statusCode := 404

	switch statusCode {
	case 200:
		fmt.Println("OK")
	case 201:
		fmt.Println("Created")
	case 404:
		fmt.Println("Not Found") // ✅ only this runs — no fallthrough
	case 500:
		fmt.Println("Internal Server Error")
	default:
		fmt.Println("Unknown status") // runs only if no case matched
	}

	// --------------------------------------------------------
	// 2. MULTIPLE VALUES PER CASE
	// --------------------------------------------------------
	// JS: stack cases on top of each other (relies on fallthrough)
	// Go: comma-separated values in a single case line — cleaner
	// --------------------------------------------------------
	switch statusCode {
	case 200, 201, 204:
		fmt.Println("success")
	case 400, 401, 403, 404:
		fmt.Println("client error")
	case 500, 502, 503:
		fmt.Println("server error")
	default:
		fmt.Println("unknown")
	}

	// --------------------------------------------------------
	// 3. EXPRESSIONLESS (NAKED) SWITCH
	// --------------------------------------------------------
	// No variable after `switch` keyword.
	// Each case is a full boolean condition.
	// Go enters the switch block 100% of the time,
	// then evaluates cases top to bottom — stops at first match.
	// Does NOT evaluate all cases — stops at first true condition.
	//
	// Use this instead of long if/else chains. Cleaner, more readable.
	//
	// JS equivalent (rough):
	// switch(true) { case x > 100: ... case x > 80: ... }
	// --------------------------------------------------------
	requestsPerMinute := 75

	switch {
	case requestsPerMinute > 100:
		fmt.Println("critical: rate limit exceeded")
	case requestsPerMinute > 80:
		fmt.Println("warning: approaching limit")
	case requestsPerMinute > 50:
		fmt.Println("moderate traffic") // ✅ this runs for 75
	default:
		fmt.Println("normal traffic")
	}

	// --------------------------------------------------------
	// 4. INIT STATEMENT in switch
	// --------------------------------------------------------
	// Same as if — declare a variable before the condition.
	// Variable is scoped to the entire switch block only.
	//
	// switch val := fn(); val { ... }
	//   vs
	// switch fn() { ... }
	//
	// Use init statement when you need to reference the result
	// inside the case bodies. Anonymous form works but you can't
	// use the value inside the cases.
	// --------------------------------------------------------
	switch role := getUserRole(2); role {
	case "admin":
		fmt.Println("Welcome Admin —", role) // can use role here
	case "moderator":
		fmt.Println("Welcome Moderator —", role)
	default:
		fmt.Println("Welcome Guest —", role)
	}
	// fmt.Println(role) // ❌ compile error — role scoped to switch block

	// --------------------------------------------------------
	// 5. FALLTHROUGH — explicit and rare
	// --------------------------------------------------------
	// Forces execution of the NEXT case's body unconditionally.
	// Does NOT check the next case's condition — blind execution.
	// Falls through ONE level only (unless that case also has fallthrough)
	//
	// ❌ fallthrough in default case → compile error
	//    "cannot fallthrough final case in switch"
	//
	// Real world use: almost never. Treat it as a code smell.
	// A function call is almost always cleaner.
	// --------------------------------------------------------
	version := 1

	switch version {
	case 1:
		fmt.Println("v1 handler")
		fallthrough // blindly runs case 2 body — no condition check
	case 2:
		fmt.Println("v2 handler")
		// no fallthrough — stops here
	case 3:
		fmt.Println("v3 handler")
	default:
		fmt.Println("unknown version")
		// ❌ fallthrough here → compile error
	}

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1. No break needed — Go auto-breaks after each case
	// 2. JS default is fallthrough, Go default is break — opposite
	// 3. Multiple values per case: case 200, 201, 204:
	// 4. Expressionless switch — no variable, each case is a condition
	//    Stops at first match — does NOT evaluate all cases
	// 5. Init statement: switch val := fn(); val { }
	//    val scoped to switch block only
	// 6. switch fn() — anonymous, can't reference result in cases
	// 7. fallthrough — blind, unconditional, one level only
	// 8. fallthrough in default → compile error
	// 9. Use expressionless switch instead of long if/else chains
	// --------------------------------------------------------
}

// ============================================================
// TOPIC 7: Looping (for and while)
// ============================================================
// Go has ONE loop keyword: `for`
// No while, no do-while — for does everything.
// Syntax is same as JS but NO parentheses around the condition.
//
// Three forms:
//   for init; condition; post { }   → classic for loop
//   for condition { }               → while loop
//   for { }                         → infinite loop (while true)
// ============================================================

package main

import "fmt"

func main() {

	// --------------------------------------------------------
	// 1. CLASSIC FOR LOOP
	// --------------------------------------------------------
	// Same as JS: for (let i = 0; i < n; i++)
	// Difference: no parentheses around the three parts
	// --------------------------------------------------------
	for i := 0; i < 5; i++ {
		fmt.Printf("Fetching page %d...\n", i+1)
	}

	// --------------------------------------------------------
	// 2. FOR AS WHILE LOOP
	// --------------------------------------------------------
	// Only a condition — no init, no post.
	// Equivalent to JS's while(condition) { }
	//
	// ⚠️  Don't mutate your config values inside the loop.
	// Keep maxRetries intact — loop on a separate counter.
	// --------------------------------------------------------
	maxRetries := 5
	attempt := 0

	for attempt < maxRetries { // ✅ condition only = while loop
		attempt++
		fmt.Printf("Attempt %d failed, retrying...\n", attempt)
	}
	fmt.Println("Max retries reached.")
	// maxRetries still = 5 here — not mutated ✅

	// --------------------------------------------------------
	// 3. INFINITE LOOP with break
	// --------------------------------------------------------
	// for { } with no condition = infinite loop = while(true) in JS
	//
	// break → exits the loop, continues after it
	// return → exits the entire function (not just the loop)
	//
	// ⚠️  Use >= not == for loop exit conditions
	//     == assumes you hit exactly that number
	//     >= is defensive — handles skipped values too
	//     eg: if requestCount somehow jumps from 4 to 6,
	//         == 5 never triggers → infinite loop → server crash
	// --------------------------------------------------------
	requestCount := 0

	for {
		if requestCount >= 5 { // ✅ >= safer than ==
			break
		}
		fmt.Printf("Processing request %d\n", requestCount+1)
		requestCount++
	}
	fmt.Println("server shutting down")

	// ⚠️  Order matters inside loops:
	// requestCount+1 → expression, temporary value, original unchanged
	// requestCount++ → statement, permanently mutates the variable
	// These are NOT the same thing

	// --------------------------------------------------------
	// 4. CONTINUE and NESTED LOOPS
	// --------------------------------------------------------
	// continue → skips the rest of the current iteration
	// break    → exits the innermost loop only
	// return   → exits the entire function (confirmed — nothing after runs)
	// --------------------------------------------------------
	bannedID := 3

	for i := 1; i <= 7; i++ {
		if i == bannedID {
			continue // skip banned user, go to next iteration
		}
		fmt.Printf("processing user %d\n", i)
	}

	// Nested loops — break only exits innermost
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fmt.Printf("page %d item %d\n", i, j)
		}
	}

	// --------------------------------------------------------
	// 5. LABELED STATEMENTS — breaking outer loops
	// --------------------------------------------------------
	// To break out of an outer loop from inside an inner loop,
	// use a labeled statement. Label can be any name.
	// `break labelName` exits the loop that has that label.
	//
	// Officially called: labeled statement in Go spec.
	// --------------------------------------------------------
outer:
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			if j == 2 {
				break outer // exits the outer loop entirely
			}
			fmt.Printf("page %d item %d\n", i, j)
		}
	}
	fmt.Println("after labeled break — this runs") // ✅ runs, return would have killed this

	// --------------------------------------------------------
	// 6. FOR RANGE — iterating over collections
	// --------------------------------------------------------
	// range is a built-in Go feature (not a function, not a type)
	// Works on: slices, arrays, maps, strings, channels
	// Returns: index + value for slices/arrays
	//
	// Forms:
	//   for i, v := range slice { }   → both index and value
	//   for i := range slice { }      → index only ✅ valid
	//   for _, v := range slice { }   → value only (blank identifier)
	//
	// ⚠️  Cannot declare i, v and then not use i → compile error
	//     Must use _ to explicitly discard: for _, v := range slice
	// --------------------------------------------------------
	routes := []string{"/health", "/users", "/posts", "/auth"}

	for i, v := range routes {
		fmt.Printf("%d: %s\n", i, v)
	}

	// Index only
	for i := range routes {
		fmt.Printf("route index: %d\n", i)
	}

	// Value only — discard index with _
	for _, v := range routes {
		fmt.Printf("route: %s\n", v)
	}

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1. One loop keyword: `for` — no while, no do-while
	// 2. No parentheses around loop condition
	// 3. for { } = infinite loop — needs break or return to exit
	// 4. break → exits innermost loop only
	// 5. return → exits the entire function, always, no exceptions
	// 6. break labelName → exits the specifically labeled loop
	// 7. continue → skips current iteration, continues loop
	// 8. Use >= not == for loop exit conditions — defensive programming
	// 9. expr+1 is temporary, var++ mutates — not the same thing
	// 10. range returns index + value — must use or discard with _
	// 11. for i := range slice → index only, valid without _
	// 12. for _, v := range slice → value only, _ discards index
	// --------------------------------------------------------
}

// ============================================================
// TOPIC 4: Arithmetic Operators
// ============================================================
// Go is strict about types in arithmetic — unlike JS where all
// numbers are just `number`, Go has int, float64, int32, int64 etc.
// This means you cannot mix types in expressions without explicit
// conversion, and integer vs float division behaves differently.
// ============================================================

package main

import "fmt"

func main() {

	// --------------------------------------------------------
	// 1. BASIC OPERATORS
	// --------------------------------------------------------
	// +  addition
	// -  subtraction
	// *  multiplication
	// /  division      ⚠️ behavior depends on type (see below)
	// %  modulo/remainder (integers only — not valid on floats)
	// --------------------------------------------------------
	a := 7
	b := 2

	fmt.Println(a + b) // 9
	fmt.Println(a - b) // 5
	fmt.Println(a * b) // 14
	fmt.Println(a % b) // 1  — remainder after division

	// --------------------------------------------------------
	// 2. INTEGER DIVISION TRAP ⚠️
	// --------------------------------------------------------
	// In JS: 7 / 2 = 3.5 (all numbers are float under the hood)
	// In Go: 7 / 2 = 3   (integer division TRUNCATES the decimal)
	//
	// This is a silent logic bug in production if you're calculating
	// prices, percentages, averages — anything that needs decimals.
	//
	// Fix: convert to float64 before dividing
	// --------------------------------------------------------
	fmt.Println(a / b)                   // 3   ← truncated, not rounded
	fmt.Println(float64(a) / float64(b)) // 3.5 ← correct

	// --------------------------------------------------------
	// 3. DIVISION BY ZERO — behavior differs by type
	// --------------------------------------------------------
	// int    / 0 → PANIC at runtime (program crashes)
	// float64/ 0 → returns +Inf, no panic
	//
	// Always guard against zero denominators in backend code:
	// if divisor == 0 { return error }
	// --------------------------------------------------------
	// fmt.Println(a / 0)        // ❌ panic: runtime error: integer divide by zero
	f := 7.0
	fmt.Println(f / 0) // +Inf — no panic, but probably not what you want

	// --------------------------------------------------------
	// 4. INCREMENT AND DECREMENT
	// --------------------------------------------------------
	// ++ and -- exist in Go but are STATEMENTS not EXPRESSIONS.
	// This means they must be on their own line.
	// You CANNOT use them inside another expression.
	//
	// JS allows:  y := x++  (expression, returns old value)
	// Go refuses: y := x++  (compile error — ++ is a statement)
	//
	// Also: no ++x or --x (prefix form) — only x++ and x-- (postfix)
	// --------------------------------------------------------
	requestCount := 0
	requestCount++ // ✅ own line — valid
	requestCount++ // 2
	requestCount-- // 1
	fmt.Println(requestCount)

	// total := requestCount++ // ❌ compile error — can't use in expression

	// --------------------------------------------------------
	// 5. CEILING DIVISION — common backend pattern
	// --------------------------------------------------------
	// Go has no built-in ceil for integers. Standard trick:
	// (a + b - 1) / b  gives ceiling division without floats.
	//
	// Why it works: adding (b-1) ensures that any remainder
	// pushes the result up to the next integer after truncation.
	//
	// Use case: pagination, chunking, rate limit windows
	// --------------------------------------------------------
	totalItems := 57
	pageSize := 10

	totalPages := (totalItems + pageSize - 1) / pageSize // 6 ✅
	lastPageItems := totalItems % pageSize               // 7
	isPartial := totalItems%pageSize != 0                // true

	fmt.Println("Total Pages:", totalPages)
	fmt.Println("Items on Last Page:", lastPageItems)
	fmt.Println("Is Last Page Partial:", isPartial)

	// Edge case: totalItems divisible by pageSize exactly
	// (60 + 10 - 1) / 10 = 69 / 10 = 6 ✅ still correct

	// --------------------------------------------------------
	// 6. RATE LIMITER MATH — real backend pattern
	// --------------------------------------------------------
	requestsPerMinute := 100
	windowSeconds := 60
	elapsedSeconds := 23
	requestsMade := 38

	remaining := requestsPerMinute - requestsMade
	// ⚠️ remaining can go negative if limit exceeded
	// In a real API you'd clamp this to 0:
	// if remaining < 0 { remaining = 0 }

	// Elapsed % of window = (elapsed / total window) * 100
	elapsedPercentage := float64(elapsedSeconds) / float64(windowSeconds) * 100

	requestsPerSecond := float64(requestsPerMinute) / float64(windowSeconds)
	limitExceeded := requestsMade > requestsPerMinute

	fmt.Printf("Remaining requests: %d\n", remaining)
	fmt.Printf("Elapsed Percentage: %.2f%%\n", elapsedPercentage) // %% = literal %
	fmt.Printf("Requests/sec allowed: %.2f\n", requestsPerSecond)
	fmt.Printf("Limit Exceeded: %t\n", limitExceeded)

	// --------------------------------------------------------
	// 7. OPERATOR PRECEDENCE — BODMAS applies
	// --------------------------------------------------------
	// Go precedence (high to low):
	//   * / %   (evaluated first, left to right)
	//   + -     (evaluated second, left to right)
	//
	// Use parentheses to make intent explicit — don't rely on
	// precedence rules in complex expressions. Future you will
	// thank present you.
	// --------------------------------------------------------
	result := 2 + 3*4 - 10/2
	fmt.Println(result) // 2 + 12 - 5 = 9

	tokenExpiry := 3600
	bufferTime := 300

	// ❌ Relies on precedence — harder to read at a glance
	halfExpiry := tokenExpiry/2 + bufferTime*2

	// ✅ Explicit parentheses — intent is clear
	halfExpiryExplicit := (tokenExpiry / 2) + (bufferTime * 2)

	fmt.Println(halfExpiry)
	fmt.Println(halfExpiryExplicit)

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1. Integer division TRUNCATES — convert to float64 for decimals
	// 2. Cannot mix types — float64(x) before operating with int
	// 3. % (modulo) works on integers only, not floats
	// 4. ++ and -- are statements — own line only, no prefix form
	// 5. int / 0 panics, float64 / 0 returns +Inf
	// 6. No ternary operator — use if/else
	// 7. Use parentheses in complex expressions — be explicit
	// 8. Ceiling division trick: (a + b - 1) / b
	// 9. Precedence: * / % are at the SAME level (left to right), then + -
	//    eg: 10 - 4 / 2 * 3 → 4/2=2 first, then 2*3=6, then 10-6=4
	//    When in doubt — use parentheses
	// --------------------------------------------------------
}

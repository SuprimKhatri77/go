// ============================================================
// TOPIC 10: Maps
// ============================================================
// Maps are Go's key-value store.
// JS equivalent: object {} or Map — but typed and stricter.
//
// map[KeyType]ValueType
// Keys must be comparable types (string, int, bool etc.)
// Keys cannot be: slices, maps, or functions
//
// Maps ARE typed — map[string]int only accepts string keys + int values.
// For fixed named fields with autocomplete → use Structs (Topic 18)
// ============================================================

package main

import (
	"fmt"
	"sort"
)

// --------------------------------------------------------
// FUNCTIONS
// --------------------------------------------------------

func getStatusMessage(codes map[int]string, code int) string {
	// comma ok pattern — safe key lookup
	if message, ok := codes[code]; ok {
		return message
	}
	return "unknown status"
}

func canWrite(perms map[string]map[string]string, user string, resource string) bool {
	// always check if outer key exists first
	if _, ok := perms[user]; !ok {
		return false
	}
	// then check inner key and value
	return perms[user][resource] == "write"
}

func recordHit(hits map[string]int, route string) {
	hits[route]++
	// ✅ No return needed — maps are REFERENCE types
	// Passing a map passes the pointer to the same underlying map
	// hit[route]++ modifies the original directly
	//
	// This is different from slices:
	// Slices need return because append may create a NEW underlying array
	// Maps always modify in place — pointer never becomes stale
	//
	// Also works without pre-initializing the key:
	// hits["/newroute"]++ → Go uses zero value (0) then increments to 1
	// No need to do hits["/newroute"] = 0 first
}

func topRoute(hits map[string]int) string {
	top := ""
	max := 0
	for route, count := range hits {
		if count > max {
			max = count
			top = route
		}
	}
	return top
	// ⚠️  If two routes have the same count, whichever comes last
	// in the random iteration wins — map order is not guaranteed
	// For deterministic results, sort keys first (see iteration section)
}

func main() {

	// --------------------------------------------------------
	// 1. DECLARING MAPS — three ways
	// --------------------------------------------------------

	// Way 1: Map literal — declare and initialize in one step
	// Use when you know the data upfront
	mapLiteral := map[int]string{
		200: "OK",
		201: "Created",
		404: "Not Found",
		500: "Internal Server Error", // trailing comma required
	}
	fmt.Println(mapLiteral)

	// Way 2: make(map[KeyType]ValueType)
	// Use when you'll populate it dynamically
	mapWithMake := make(map[int]string)
	mapWithMake[200] = "OK"
	mapWithMake[201] = "Created"
	mapWithMake[404] = "Not Found"
	mapWithMake[500] = "Internal Server Error"
	fmt.Println(mapWithMake)

	// Way 3: Empty literal, assign later
	// Same as make for most purposes
	emptyMap := map[int]string{}
	emptyMap[200] = "OK"
	fmt.Println(emptyMap)

	// --------------------------------------------------------
	// ZERO VALUE of a map — nil (same pattern as slices)
	// --------------------------------------------------------
	// var m map[int]string → nil map
	// Reading  from nil map → returns zero value, no panic ✅
	// Writing  to nil map  → PANIC at runtime ❌
	//
	// Fix: always initialize before writing
	//   var m = map[int]string{}   ✅
	//   m := map[int]string{}      ✅
	//   m := make(map[int]string)  ✅
	//
	// With := you MUST initialize — even if empty: map[int]string{}
	// --------------------------------------------------------
	var nilMap map[int]string
	fmt.Println(nilMap == nil) // true
	fmt.Println(nilMap)        // map[] — prints same as empty map
	fmt.Println(nilMap[200])   // "" — zero value, no panic
	// nilMap[200] = "OK"      // ❌ panic: assignment to entry in nil map

	// --------------------------------------------------------
	// 2. THE COMMA OK PATTERN — safe key lookup
	// --------------------------------------------------------
	// In JS: obj[key] → undefined if missing (silent bug)
	// In Go: map[key] → zero value if missing (silent bug)
	//
	// Fix: comma ok pattern
	// value, ok := map[key]
	// ok is true if key exists, false if not
	// --------------------------------------------------------
	statusMessages := map[int]string{200: "OK", 404: "Not Found"}

	// ❌ Unsafe — returns "" silently if key missing
	msg := statusMessages[500]
	fmt.Println(msg) // ""

	// ✅ Safe — check existence first
	if message, ok := statusMessages[500]; ok {
		fmt.Println(message)
	} else {
		fmt.Println("key not found")
	}

	fmt.Println(getStatusMessage(mapLiteral, 700)) // "unknown status"

	// --------------------------------------------------------
	// 3. DELETE and ITERATION
	// --------------------------------------------------------
	// delete(map, key) — removes a key, no error if key doesn't exist
	//
	// Iteration order is NOT guaranteed in Go maps
	// Every range over a map may give different order
	// For deterministic/sorted order: collect keys → sort → iterate
	// --------------------------------------------------------
	sessions := map[string]string{
		"token_abc": "user_1",
		"token_xyz": "user_2",
		"token_123": "user_3",
		"token_def": "user_4",
	}
	delete(sessions, "token_xyz")

	// ✅ Cleaner iteration — use both k and v directly
	for k, v := range sessions {
		fmt.Printf("%s: %s\n", k, v)
	}

	// Sorted iteration — for deterministic output
	keys := []string{}
	for k := range sessions {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %s\n", k, sessions[k])
	}

	// --------------------------------------------------------
	// 4. MAPS ARE REFERENCE TYPES ⚠️
	// --------------------------------------------------------
	// Same trap as slices — assignment copies the pointer not the data
	// Unlike slices there is NO built-in copy() for maps
	//
	// Manual copy — loop and assign each key:
	// for k, v := range original { copied[k] = v }
	//
	// JS equivalents for copying:
	//   Object.assign({}, original)
	//   structuredClone(original)
	//   JSON.parse(JSON.stringify(original))
	// --------------------------------------------------------
	original := map[string]int{"requests": 100, "errors": 5}
	ref := original
	ref["requests"] = 999
	fmt.Println("original:", original["requests"]) // 999 — mutated ⚠️
	fmt.Println("ref:", ref["requests"])           // 999

	// ✅ Manual deep copy
	copied := make(map[string]int)
	for k, v := range original {
		copied[k] = v
	}
	copied["requests"] = 1
	fmt.Println("original after copy:", original["requests"]) // 999 — unchanged ✅

	// --------------------------------------------------------
	// 5. NESTED MAPS
	// --------------------------------------------------------
	// map[string]map[string]string
	// Always check outer key exists before accessing inner map
	// Accessing missing outer key returns nil inner map
	// Then accessing nil inner map returns zero value — silent bug
	// --------------------------------------------------------
	perms := map[string]map[string]string{
		"alice": {"posts": "write", "comments": "read"},
		"bob":   {"posts": "read", "comments": "read"},
	}
	fmt.Println(canWrite(perms, "alice", "posts")) // true
	fmt.Println(canWrite(perms, "bob", "posts"))   // false
	fmt.Println(canWrite(perms, "ghost", "posts")) // false — user doesn't exist

	// --------------------------------------------------------
	// 6. REAL WORLD — route hit counter
	// --------------------------------------------------------
	routeHits := map[string]int{}
	recordHit(routeHits, "/users")
	recordHit(routeHits, "/posts")
	recordHit(routeHits, "/users")
	recordHit(routeHits, "/users")
	recordHit(routeHits, "/auth")
	recordHit(routeHits, "/posts")
	fmt.Println(routeHits)
	fmt.Println("top route:", topRoute(routeHits)) // /users

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1.  Syntax: map[KeyType]ValueType
	// 2.  Three ways to declare:
	//       map literal     → map[k]v{k:v, ...}   when data known upfront
	//       make            → make(map[k]v)        when populating dynamically
	//       empty literal   → map[k]v{}            same as make for most uses
	// 3.  Zero value is nil — reading is safe, writing panics
	//     Always initialize before writing
	// 4.  Missing key returns zero value silently — use comma ok pattern
	//     value, ok := m[key]
	// 5.  delete(m, key) — safe even if key doesn't exist
	// 6.  Iteration order is NOT guaranteed — sort keys for deterministic order
	// 7.  Maps are REFERENCE types — assignment copies pointer not data
	// 8.  No built-in copy() for maps — manually loop and assign
	// 9.  Passing map to function modifies original — no return needed
	//     (unlike slices where append may create new underlying array)
	// 10. Nested maps — always check outer key exists before inner access
	// 11. Map keys must be comparable — no slices, maps, or funcs as keys
	// --------------------------------------------------------
}

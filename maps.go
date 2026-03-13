package main

import "fmt"

// LatLong defines a custom coordinate structure.
type LatLong struct {
	Lat, Long float64
}

// 1. DECLARING WITHOUT INITIALIZING
// This creates a 'nil' map. It has no underlying storage.
// READS from a nil map are safe (return zero value), but WRITES will panic.
var m map[string]LatLong

// 2. MAP LITERAL SYNTAX
// This is the cleanest way to initialize a map with known starting data.
// Note: we can omit the type 'LatLong' inside the brackets for cleaner code.
var literalMapDeclaration = map[string]LatLong{
	"Bell Labs": {40.68433, -74.39967},
	"Google":    {37.42202, -122.08408},
}

func MapExample() {
	// 3. INITIALIZING WITH MAKE
	// 'make' allocates the underlying hash table and returns a ready-to-use map.
	m = make(map[string]LatLong)
	m["Bell Labs"] = LatLong{40.6843, -74.39967}

	fmt.Println("Bell Labs Coords:", m["Bell Labs"])
	fmt.Println("All Locations:", literalMapDeclaration)

	// --- CRUD OPERATIONS (Create, Read, Update, Delete) ---
	val := make(map[string]int)

	// CREATE / UPDATE
	val["Answer"] = 7
	val["Answer"] = 42 // Overwrites existing value
	fmt.Println("Updated Value:", val["Answer"])

	// DELETE
	// If the key doesn't exist, delete does nothing (no error).
	delete(val, "Answer")

	// READ & VALIDATION (The "Comma ok" Idiom)
	// In Go, map lookups return two values:
	// 1. The value (or the zero-value of the type if not found)
	// 2. A boolean 'ok' (true if the key exists)
	v, ok := val["Answer"]

	// %d for int, %t for boolean
	fmt.Printf("Value: %d, IsPresent: %t\n", v, ok)

	// 1. ORDER IS RANDOMIZED
	// Go randomizes map iteration intentionally so developers don't rely
	// on a specific order. Each run might produce different results.
	fmt.Println("--- Iterating over Map ---")
	for key, value := range literalMapDeclaration {
		// key is the string, value is the LatLong struct
		fmt.Printf("Key: %-10s | Lat: %.2f, Long: %.2f\n", key, value.Lat, value.Long)
	}

	// 2. STRUCTS VS MAPS
	// You cannot 'range' over a struct because its fields are fixed at compile time.
	// Maps are dynamic collections; structs are structured "objects".
}

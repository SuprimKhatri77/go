package main

import (
	"fmt"
	"math"
)

// Vrtx (Vertex) is a simple struct.
type Vrtx struct {
	X, Y float64
}

// MyFloat is a 'defined type' based on float64.
// This allows us to attach methods to a primitive-like value.
type MyFloat float64

// 1. VALUE RECEIVER: (v Vrtx)
// This method operates on a COPY of the struct.
// Use this for reading data without modifying the original.
func (v Vrtx) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// 2. POINTER RECEIVER: (v *Vrtx)
// This method operates on the ORIGINAL struct in memory.
// Use this to modify (mutate) the caller or to avoid copying large structs.
func (v *Vrtx) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

// 3. NON-STRUCT METHOD
// You can define methods on any type you define in the same package.
func (f MyFloat) Fbs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func MethodExample() {
	v := Vrtx{3, 4}

	// Go is smart: even though Scale() needs a pointer (*Vrtx),
	// and 'v' is a value (Vrtx), Go interprets v.Scale(10) as (&v).Scale(10).
	v.Scale(10)

	fmt.Println("Absolute Value:", v.Abs())

	f := MyFloat(-math.Sqrt(2))
	fmt.Println("Float Absolute Value:", f.Fbs())
}

/*
---------------------------------------------------------

	CASE 1: THE METHOD (Receiver Syntax)
	---------------------------------------------------------
	USE CASE:
	- When the action is "owned" by the object (Encapsulation).
	- When you need to satisfy an Interface (e.g., Stringer, Reader).
	- When you want "Chainability" (e.g., v.Scale(2).Abs()).
*/
func (v Vrtx) AbsMethod() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

/*
---------------------------------------------------------

	CASE 2: THE FUNCTION (Standard Syntax)
	---------------------------------------------------------
	USE CASE:
	- Pure utility logic that doesn't "belong" to one thing.
	- When you are transforming one type into another (e.g., JSON to Struct).
	- When the logic is "Stateless" (doesn't care about internal data).
*/
func AbsFunction(v Vrtx) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func ComparisonExample() {
	v := Vrtx{3, 4}

	// Method call: The 'v' is the subject. Clean and readable.
	fmt.Println(v.AbsMethod())

	// Function call: 'v' is just an input parameter.
	fmt.Println(AbsFunction(v))
}

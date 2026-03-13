package main //ideally same as file name

import (
	"fmt" //this is provided by the go standard library
	"math"
	"math/cmplx"
	"math/rand"
)

// go build main.go creates an executable file
// ./main runs the executable file

// go run main.go runs the file directly by directly compiling and storing the binary temporarily somewhere in the system. It also caches the binary for future use. If the file is run again, it will use the cached binary if the file is not changed while go  build main.go creates an executable file each time.

// functions
/* func functionName(parameters) returnType {
 		function body
	}
*/
func add(a int, b int) int {
	return a + b
}

/*
we can return multiple values from a function by specifying the return types in the function signature using brackets

also if the parameters are of the same type, we can specify the type only once.
eg: func swap(x,y string)(string,string) instead of func swap(x string, y string)(string, string)
*/
func swap(x, y string) (string, string) {
	return y, x
}

/*
split demonstrates named return values and "naked" returns.

1. Named Returns: The return variables 'x' and 'y' are defined in the function
   signature. They are automatically initialized to their zero values (0 for int)
   at the start of the function.

2. Assignment: We assign values directly to 'x' and 'y'. While you can define
   other variables in the function, only the values stored in 'x' and 'y'
   will be sent back by the naked return.

3. Naked Return: The 'return' statement without arguments returns the current
   values of the named result parameters.

*/
// func split(sum int) (x, y int) {
// 	x = sum + 7
// 	y = sum - x
// 	return
// }

// this is much better than the previous one
func split(sum int) (int, int) {
	x := sum * 4 / 9
	y := sum - x
	return x, y
}

/*
if the variables are of same datatype we can declare them in a single line and define the datatype at the end
eg: var c,python,java bool
here the variables are of type bool and are initialized to their zero values (false)
*/

var c, python, java bool

/*
1. Factored Declaration: We can group multiple variable declarations into a
   single 'var' block for better readability.
2. Variable Names: 'ToBe', 'MaxInt', and 'z' are custom identifiers. They can
   be any name (e.g., 'isReady', 'myNumber') as long as they aren't Go keywords.
3. uint64 & Bit Shifting: 'MaxInt' is an unsigned 64-bit integer. The expression
   '1<<64 - 1' creates a number with 64 bits all set to 1 (the maximum value).
4. complex128: A type representing a complex number with a real and an
   imaginary part (i), using 64 bits for each part (128 bits total).
*/

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

/*
we cannot use the := operator outside of a function to declare variables , we must make use of var keyword
*/

// this is the entry point of the file
// this where the program execution will start
func main() {
	// declaring variables inside a function
	// using var keyword
	// var i int
	i := 0

	fmt.Println("hello world")
	fmt.Println("Random number using math/rnd: ", rand.Intn(8))
	fmt.Println("Value of Pi: ", math.Pi)

	fmt.Println("Sum of 7 and 8 is: ", add(7, 8))

	a, b := swap("hello", "world")
	fmt.Println(a, b)

	fmt.Println(split(17))

	fmt.Println(i, c, python, java)

	/*
		%T is used to print the type of the variable
		%v is used to print the value of the variable
	*/
	fmt.Printf("Type : %T Value: %v\n", ToBe, ToBe)
	fmt.Printf("Type : %T Value: %v\n", MaxInt, MaxInt)
	fmt.Printf("Type : %T Value: %v\n", z, z)

	/*
		Zero Value Initialization:
		In Go, variables declared without an initial value are automatically
		given a "zero value" for safety.

		1. int: 0
		2. float64: 0
		3. bool: false
		4. string: "" (empty string)

		Note on Printing:
		- %v (Value): Prints '0 0 false '. The empty string is invisible.
		- %s (String): Prints '0 0 false '. Also invisible; looks like a trailing space.
		- %q (Quoted): Prints '0 0 false ""'. This is the best way to visually
		confirm a string is empty because it shows the literal quotes.
	*/
	var integer int
	var f float64
	var boolean bool
	var s string

	// Using %q for the string is smart—it's the only way to "see" an empty string!
	fmt.Printf("%v %v %v %s\n", integer, f, boolean, s)

	/*
	   Type Conversion in Action:
	   Go is strictly typed and does not support implicit conversion (widening).
	   We must manually convert variables when the types don't match exactly.
	*/

	// 1. Initializing two integers.
	var x, y int = 3, 4

	// 2. math.Sqrt requires a float64.
	// Even though (x*x + y*y) is 25, we must wrap it in float64()
	// or the compiler will scream.
	var floatingPoint float64 = math.Sqrt(float64(x*x + y*y))

	// 3. Converting back to an unsigned integer (uint).
	// Note: Converting from float to int will truncate (chop off) decimals.
	var z uint = uint(floatingPoint)

	/*
	   WHY IT PRINTED "5" AND NOT "5.0":
	   By default, fmt.Println (and the %v verb) is "human-friendly."
	   If a float64 is exactly a whole number, Go trims the trailing zeros
	   to keep the output clean. It is still a float64 in memory!
	*/
	fmt.Println("Default Output:", floatingPoint) // Prints: 5
	/*
	   HOW TO PRINT THE DECIMAL:
	   Use fmt.Printf with a float verb like %f.
	   - %f: Prints with default precision (6 decimal places).
	   - %.1f: Prints with exactly 1 decimal place.
	*/
	fmt.Printf("Forced Float: %f\n", floatingPoint)    // Prints: 5.000000
	fmt.Printf("Fixed Decimal: %.1f\n", floatingPoint) // Prints: 5.0
	fmt.Println(x, y, floatingPoint, z)                //

	/*
	   Type Inference (The "Walrus" Operator):
	   Go determines the type of 'v' at compile-time based on the assigned value.
	   Since 42 is a whole number, 'v' becomes an 'int'.

	   CRITICAL: Go is Statically Typed. Once 'v' is an 'int', it can NEVER
	   hold a float or string. There is NO "Implicit Type Conversion" in Go.
	*/
	v := 42
	fmt.Printf("Type: %T Value: %v\n", v, v)

	/*
	   Constants:
	   Constants are immutable (read-only). Their values must be known at compile-time.

	   Note: In Go, constants can be 'untyped', allowing them to be used
	   more flexibly than variables until they are assigned to a specific type.
	*/
	const Pi float64 = 3.14159
	fmt.Printf("Pi: %f\n", Pi)

	Loops()

	SwitchStatement()
	SwitchDay()
	SwitchWithoutCondition()

	DeferExample()

	ReferenceTypes()
	StructExample()
	ArrayExample()

	SliceExample()
	SliceExampleTwo()
	SliceExampleThree()
	SliceExampleFour()
	TwoDimensionSlice()
	AppendExample()
	RangeExample()

	MapExample()

	AnonymousFunction()

	ClosureExample()

	MethodExample()

	InterfaceExample()

	GenericsExample()

	ConcurrencyExample()

	MutexExample()

}

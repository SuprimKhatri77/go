package main

import (
	"fmt"
	"time"
)

// empty interfae
type any interface{}

type I interface {
	M() // The contract: Must have a method M()
}

type T struct {
	S string
}

// T implements I implicitly.
// There is no "implements" keyword in Go.
// Because T has the method M(), it qualifies automatically.
func (t T) M() {
	fmt.Println(t.S)
}

func InterfaceExample() {
	// We create a variable of type interface 'I'.
	// We assign it a value of type 'T'.
	// Go checks: "Does T have an M() method?" Yes.
	var i I = T{"suprim77"}

	// When we call i.M(), the interface looks at its internal "Concrete Type"
	// to find the actual M() method belonging to T and executes it.
	i.M()

	var j any

	describe(j)

	j = 42
	describe(j)

	j = "slash77"
	describe(j)

	TypeAssertion()
	do(21)
	do("Hello")
	do(true)

	a := Person{"Mark Rober", 20}
	b := Person{"Cristiano Ronaldinho", 77}

	fmt.Println(a, b)

	if err := run(); err != nil {
		fmt.Println(err)
	}

}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

// type assertion

func TypeAssertion() {
	var i interface{} = "hello"
	s := i.(string)
	fmt.Println(s)

	s, ok := i.(string)
	fmt.Println(s, ok)

	f, ok := i.(float64)
	fmt.Println(f, ok)

	// f = i.(float64) //panic
	// fmt.Println(f)

}

func do(i interface{}) {
	// cannot do: i = i+2 or anything before the type assertion is done or until we cast the type explicitly
	switch v := i.(type) {
	case int:
		fmt.Printf("Twice %v is %v\n", v, v*2)
	case string:
		fmt.Printf("%v is %v\n", v, len(v))
	default:
		fmt.Printf("%v of type %T\n", v, i)
	}
}

type Person struct {
	name string
	age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%v (%v years)", p.name, p.age)
}

type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s", e.When, e.What)
}

func run() error {
	return &MyError{
		time.Now(),
		"It didn't work",
	}
}

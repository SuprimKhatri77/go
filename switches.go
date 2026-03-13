package main

import (
	"fmt"
	"runtime"
	"time"
)

func SwitchStatement() {
	// 1. INITIALIZATION: os := runtime.GOOS (runs once)
	// 2. CONDITION: os (the value we are checking)
	// 3. SCOPE: 'os' only exists inside these { } brackets

	// --- THE SHORTHAND PATTERN ---
	// Logic: [Initialization]; [Condition]

	// Use this when:
	// 1. You only need the variable for this specific check.
	// 2. You want to keep your "global" function memory clean.
	// 3. You are calling a function that returns a value + an error.
	switch os := runtime.GOOS; os {

	case "darwin": // This is for macOS
		fmt.Println("OS X.")

	case "linux":
		fmt.Println("Linux.")

	default: // Catch-all for Windows or others
		fmt.Printf("%s.", os)
	}

	// fmt.Println(os) <- If we uncommented this, it would throw an ERROR.
	// Why? Because 'os' died the moment the switch block ended.
}

func SwitchDay() {
	fmt.Println("When is Saturday?")

	// 1. Get today's day (Sunday=0, Monday=1 ... Saturday=6)
	// Since today is Saturday, today = 6.
	today := time.Now().Weekday()

	// 2. Proof of Type: 'today' is a Weekday type, but we cast to int to see the 6.
	fmt.Println("Today's integer is:", int(today))

	// 3. THE SWITCH: We are looking for the value 6 (time.Saturday).
	switch time.Saturday {

	case today:
		// Logic: Does 6 == 6? YES.
		// This is what prints on your terminal today!
		fmt.Println("Today is Saturday.")

	case today + 1:
		// Logic: Does 6 == (6 + 1)?
		// 6 + 1 = 7. Since 6 != 7, this case is skipped.
		// Note: As you saw, Go prints %!Weekday(7) because 7 isn't a real day.
		fmt.Println("Tomorrow is Saturday.")

	case today + 5:
		// Logic: Does 6 == (6 + 5)?
		// 6 + 11 = 17. Since 6 != 17, this is skipped.
		fmt.Println("Saturday is 5 days from today.")

	default:
		// If today was Tuesday(2), Wednesday(3), or Thursday(4),
		// none of the math above would equal 6, so we'd end up here.
		fmt.Println("Too far away.")
	}
}

func SwitchWithoutCondition() {
	t := time.Now()

	fmt.Println("current hour: ", t.Hour())

	/*
		This is a swtich statement without a condition aka Naked Switch Statement.
		It can be used instead of using if elseif elseif else... looks clean.
	*/
	switch {
	case t.Hour() < 12:
		fmt.Println("Good Morning")
	case t.Hour() < 18:
		fmt.Println("Good Afternoon")
	default:
		fmt.Println("Good Evening")
	}
}

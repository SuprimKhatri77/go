package main

import (
	"fmt"
	"time"
)

/*
========================
GO CONCURRENCY
========================

1. Concurrency
--------------
Concurrency means multiple tasks making progress during the same time period.

The tasks may not run at the exact same moment, but the program can switch
between them so they appear to run together.

Example:
- Handling multiple HTTP requests
- Calling multiple APIs
- Processing many files

Concurrency helps programs use time efficiently when tasks are waiting
(for network, disk, etc).

Note:
Concurrency ≠ Parallelism

Concurrency  -> managing multiple tasks at once
Parallelism  -> running tasks literally at the same time on multiple CPU cores


2. Goroutines
-------------
A goroutine is a lightweight thread managed by the Go runtime.

It allows a function to run concurrently with other functions.

Syntax:

    go functionName()

Example:

    go say("hello")

This starts the function in a separate goroutine and the program continues
executing the next lines without waiting for it.

Important characteristics:
- Very lightweight (~2KB stack)
- Go can run thousands or even millions of goroutines
- Managed by Go scheduler (not OS threads directly)


3. Channels
-----------
Channels are used for communication between goroutines.

They allow goroutines to safely send and receive data.

Philosophy of Go:

    "Do not communicate by sharing memory;
     instead, share memory by communicating."

Create a channel:

    c := make(chan int)

Send data into channel:

    c <- value

Receive data from channel:

    value := <-c

Channels are blocking by default:

    Send blocks until a receiver is ready
    Receive blocks until a value is sent

This makes channels useful for synchronization between goroutines.


4. Example Workflow
-------------------

    goroutine A ----\
                      > channel ----> main goroutine
    goroutine B ----/

Multiple goroutines send results into a channel,
and another goroutine (often main) collects them.


5. Common Real-World Use Cases
------------------------------

- Web servers handling many requests simultaneously
- Fetching multiple APIs concurrently
- Worker pools processing jobs
- Background tasks (emails, logging, analytics)
- File or image processing


6. Simple Mental Model
----------------------

goroutine  -> worker
channel    -> pipe/message queue
main       -> coordinator

Workers do tasks concurrently and communicate results
through channels.
*/

// say prints a string 5 times with a small delay between prints.
// This function will be used to demonstrate goroutines (concurrent execution).
func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond) // simulate some work / delay
		fmt.Println(s)
	}
}

// Sum calculates the total of all integers in the slice `s`
// and sends the result through the channel `c`.
//
// chan int means the channel can only send/receive integers.
//
// Channels are used for communication between goroutines.
func Sum(s []int, c chan int) {
	sum := 0

	// iterate through the slice and accumulate the sum
	for _, v := range s {
		sum += v
	}

	// send the computed sum to the channel
	// this will block send until another goroutine receives it
	c <- sum
}

// fibonacci generates the first n Fibonacci numbers and sends them into the channel c
func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x // send the current Fibonacci number into the channel
		x, y = y, x+y
	}

	// IMPORTANT: close the channel when done sending values
	// - If we don't close the channel, the receiver using `for i := range c`
	//   will wait forever for more values and cause a deadlock:
	//     fatal error: all goroutines are asleep - deadlock!
	// - Closing signals to the receiver: "no more values will come"
	close(c)
}

// fibonacciWithSwitch generates Fibonacci numbers and sends them to channel `c`
// It keeps running in an infinite loop until a signal is received on the `quit` channel
func fibonacciWithSelect(number, stop chan int) {
	x, y := 0, 1

	for {
		select {
		case number <- x:
			// Send the current Fibonacci number to the channel
			// BLOCKS until some goroutine reads from `c`
			x, y = y, x+y
		case <-stop:
			// Stop signal received → exit the function
			fmt.Println("quit")
			return
		}
	}

}

func ConcurrencyExample() {

	// `go` starts a new goroutine (lightweight thread)
	// say("world") runs concurrently with the rest of the program
	go say("world")

	// this runs in the main goroutine
	// output order between "world" and "hello" will be unpredictable
	say("hello")

	// sample slice of integers
	s := []int{7, 2, 8, -9, 4, 0}

	// create a channel that carries integers
	c := make(chan int)

	// run two goroutines that each sum half of the slice
	go Sum(s[:len(s)/2], c) // first half
	go Sum(s[len(s)/2:], c) // second half

	// receive results from the channel
	// `<-c` waits (blocks) until a value is sent to the channel
	// since two goroutines send values, we receive twice
	x, y := <-c, <-c

	// print both partial sums and their total
	fmt.Println(x, y, x+y)

	// create a buffered channel of capacity 10
	channel := make(chan int, 10)

	// start the Fibonacci generator as a goroutine
	// this runs concurrently with the main goroutine
	go fibonacci(cap(channel), channel)

	// receive values from the channel until it is closed
	for i := range channel {
		// NOTE: `i` here is NOT an index!
		// - When ranging over a channel, Go only returns the value sent to the channel
		// - There is no concept of an "index" for channels (they are queues, not arrays)
		// - You could rename it to `value` for clarity:
		//     for value := range channel { ... }
		fmt.Println(i)
	}

	numbers := make(chan int) // unbuffered channel for Fibonacci numbers
	stop := make(chan int)    // channel to signal generator to stop

	// Anonymous goroutine: reads 10 Fibonacci numbers from `chann` and prints them
	// Then it sends a stop signal to the generator via `quit`
	go func() {
		for i := 0; i < 10; i++ {
			// Receive a value from the channel
			// NOTE: <-chann will BLOCK until fibonacciWithSwitch sends a value
			fmt.Println(<-numbers)
		}
		// After receiving 10 numbers, signal the generator to stop
		stop <- 0
	}()

	// Call the Fibonacci generator in the main goroutine
	// Runs concurrently with the anonymous receiver goroutine
	// The for{} loop inside fibonacciWithSwitch will run indefinitely
	// until the `quit` channel receives a value
	fibonacciWithSelect(numbers, stop)

	// Create a tick channel that fires every 100ms
	tick := time.Tick(100 * time.Millisecond)

	// Create a boom channel that fires once after 500ms
	boom := time.After(500 * time.Millisecond)

	// Infinite loop demonstrating select with tick, boom, and default
	for {
		select {
		case <-tick:
			// Tick case executes every 100ms
			fmt.Println("tick.")
		case <-boom:
			// Boom case executes once after 500ms and exits the loop
			fmt.Println("BOOM!")
			return
		default:
			// Default case executes if no other channels are ready
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond) // simulate work
		}
	}

}

/*
========================
GO CONCURRENCY CHEAT SHEET
========================

1. Concurrency
--------------
Concurrency means structuring a program so multiple tasks can make progress
during the same time period.

It improves efficiency when tasks wait for I/O (network, disk, APIs).

Example use cases:
- Web servers handling multiple requests
- Calling multiple APIs simultaneously
- Processing files/images
- Background jobs (emails, logging)

Concurrency ≠ Parallelism
Concurrency  -> managing many tasks
Parallelism  -> tasks running literally at the same time (multiple CPU cores)


2. Goroutines
-------------
A goroutine is a lightweight thread managed by Go runtime.

Start a goroutine using the `go` keyword.

Example:

    go doWork()

The function runs concurrently while the program continues executing.

Properties:
- Very cheap (~2KB stack)
- Millions can exist
- Managed by Go scheduler


3. Channels
-----------
Channels allow goroutines to communicate safely.

Create a channel:

    c := make(chan int)

Send data:

    c <- value

Receive data:

    value := <-c

Channels block by default:

Send blocks until a receiver exists.
Receive blocks until a value is sent.

This naturally synchronizes goroutines.


4. Buffered Channels
--------------------

Create channel with capacity:

    c := make(chan int, 3)

Now up to 3 values can be sent before blocking.


5. WaitGroup (like Promise.all)
-------------------------------

Used when you want to wait for multiple goroutines to finish.

Example:

    var wg sync.WaitGroup

    wg.Add(2)

    go func() {
        defer wg.Done()
        task1()
    }()

    go func() {
        defer wg.Done()
        task2()
    }()

    wg.Wait() // waits for both goroutines


6. Select
---------
Select waits on multiple channel operations.

Example:

    select {
    case v := <-c1:
        fmt.Println(v)
    case v := <-c2:
        fmt.Println(v)
    }

Used for:
- timeouts
- multiple channel listeners
- cancellation


7. Worker Pool Pattern
----------------------

Used when processing many jobs concurrently.

Example idea:

jobs -> channel -> workers -> results

Workers pull jobs from a channel and process them.


8. Context (cancellation / timeouts)
------------------------------------

Used to cancel long-running operations.

Example:

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

Often used in:
- HTTP handlers
- database queries
- API calls


9. Common Pitfalls
------------------

Race Conditions
    multiple goroutines modify same data

Deadlocks
    goroutines waiting forever on channels

Goroutine Leaks
    goroutines that never exit

Use tools:

    go run -race main.go


10. Mental Model
----------------

goroutine  -> worker
channel    -> pipe/message queue
main       -> coordinator

Workers perform tasks concurrently and communicate
results through channels.
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

// -----------------------------------------------------------------------------
// GOROUTINES AND CONCURRENCY
// -----------------------------------------------------------------------------
// a goroutine is a lightweight thread managed by the Go runtime, not the OS
// spawning one costs ~2-8KB of stack memory vs ~1MB for an OS thread
// you can have thousands of goroutines running concurrently without issues
//
// the Go runtime has its own scheduler (M:N scheduler)
// it maps many goroutines (M) onto fewer OS threads (N)
// you don't manage threads directly — Go handles it
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// 1. LAUNCHING A GOROUTINE
// -----------------------------------------------------------------------------
// just add "go" before a function call
// the function runs concurrently — main doesn't wait for it

func basicGoroutine() {
	go fmt.Println("i run concurrently")
	fmt.Println("main keeps going immediately")
	time.Sleep(100 * time.Millisecond) // give goroutine time to finish
}

// anonymous function goroutine — most common pattern
func anonGoroutine() {
	go func() {
		fmt.Println("anonymous goroutine")
	}()
	time.Sleep(100 * time.Millisecond)
}

// -----------------------------------------------------------------------------
// 2. THE FIRE AND FORGET PROBLEM
// -----------------------------------------------------------------------------
// goroutines are fire and forget — main doesn't know when they finish
// if main exits, ALL goroutines are killed immediately, no cleanup
//
// bad way to wait: time.Sleep
// this is fragile — you're guessing how long goroutines take
//
// real way to wait: sync.WaitGroup (see section 4)

func fireAndForget() {
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("this might never print if main exits first")
	}()

	time.Sleep(1 * time.Second) // only waits 1s but goroutine needs 2s — MISSED
	fmt.Println("main done")
}

// -----------------------------------------------------------------------------
// 3. time.Sleep UNITS — common mistake
// -----------------------------------------------------------------------------
// time.Sleep takes a time.Duration (which is int64 nanoseconds under the hood)
// passing a bare number like time.Sleep(300) = 300 NANOSECONDS = basically nothing

func sleepUnits() {
	time.Sleep(300)                                // 300 nanoseconds — almost instant
	time.Sleep(300 * time.Millisecond)             // 300ms — noticeable
	time.Sleep(time.Second)                        // 1 second
	time.Sleep(2 * time.Second)                    // 2 seconds
	time.Sleep(500 * time.Millisecond)             // 500ms
	time.Sleep(time.Second + 500*time.Millisecond) // 1.5 seconds
}

// -----------------------------------------------------------------------------
// 4. WAITGROUP — proper way to wait for goroutines
// -----------------------------------------------------------------------------
// sync.WaitGroup is a counter
// Add(n) — increment by n before launching goroutines
// Done() — decrement by 1 (call inside each goroutine, usually with defer)
// Wait() — blocks until counter hits 0
//
// this is covered fully in topic 19 (WaitGroups and Mutex)
// shown here for context since Sleep-based waiting is not production code

func withWaitGroup() {
	var wg sync.WaitGroup

	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("goroutine %d done\n", id)
		}(i)
	}

	wg.Wait() // blocks until all 5 goroutines call Done()
	fmt.Println("all goroutines finished")
}

// -----------------------------------------------------------------------------
// 5. RACE CONDITION
// -----------------------------------------------------------------------------
// when multiple goroutines read/write the same variable simultaneously
// the result is non-deterministic — you get different values each run
// Go has a race detector: go run -race main.go
//
// fix: sync.Mutex (topic 19) or atomic operations or channels (topic 18)

func raceCondition() {
	counter := 0

	for range 100 {
		go func() {
			counter++ // NOT safe — read-increment-write is 3 ops, not atomic
			// another goroutine can interrupt between any of those 3 ops
		}()
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Println("counter:", counter) // will likely not be 100
	// some increments get lost silently — no panic, no error, just wrong data
}

// counter++ is actually three operations:
//   temp = counter      (read)
//   temp = temp + 1     (increment)
//   counter = temp      (write)
//
// if two goroutines both read counter=5 at the same time,
// both write 6, and one increment is lost forever

// -----------------------------------------------------------------------------
// 6. CLOSURE BUG IN GOROUTINES
// -----------------------------------------------------------------------------
// classic Go gotcha when launching goroutines in a loop
// goroutines capture the variable reference, not the value at launch time
// by the time the goroutine runs, the loop may have already moved on

func closureBug() {
	// BUGGY VERSION
	// all goroutines share the same 'i' variable
	// loop finishes (i=5) before goroutines run, so all print 5
	for i := range 5 {
		go func() {
			fmt.Println(i) // captures reference to i, not the value
		}()
	}
	time.Sleep(100 * time.Millisecond)

	// FIXED VERSION
	// pass i as an argument — each goroutine gets its own copy
	for i := range 5 {
		go func(i int) {
			fmt.Println(i) // i here is a local copy, not the loop variable
		}(i)
	}
	time.Sleep(100 * time.Millisecond)
}

// in Go 1.22+ the loop variable bug is fixed for range loops
// but the fix (passing as arg) is still the clearest and safest pattern

// -----------------------------------------------------------------------------
// 7. CONCURRENT VS SEQUENTIAL — measuring the difference
// -----------------------------------------------------------------------------

func fetchUser() string {
	time.Sleep(3 * time.Second)
	return "user data"
}

func fetchOrders() string {
	time.Sleep(2 * time.Second)
	return "orders data"
}

func fetchProducts() string {
	time.Sleep(4 * time.Second)
	return "products data"
}

// sequential: total time = 3 + 2 + 4 = 9 seconds
func sequentialFetch() {
	start := time.Now()
	fetchUser()
	fetchOrders()
	fetchProducts()
	fmt.Println("sequential took:", time.Since(start)) // ~9s
}

// concurrent: total time = slowest call = 4 seconds
func concurrentFetch() {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		fetchUser()
	}()
	go func() {
		defer wg.Done()
		fetchOrders()
	}()
	go func() {
		defer wg.Done()
		fetchProducts()
	}()

	wg.Wait()
	fmt.Println("concurrent took:", time.Since(start)) // ~4s
}

// this is the core value prop of goroutines for backend work:
// n independent operations that each take t seconds
// sequential: n * t
// concurrent: max(t1, t2, ..., tn)

// -----------------------------------------------------------------------------
// 8. WORKER POOL PATTERN (basic)
// -----------------------------------------------------------------------------
// instead of launching one goroutine per job (could be thousands)
// launch a fixed number of workers and distribute jobs across them
//
// the full version uses channels (topic 18) for dynamic job distribution
// this shows the concept with manual slice splitting

func workerPool() {
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var wg sync.WaitGroup

	chunks := [][]int{
		jobs[0:3], // worker 1 gets jobs 1-3
		jobs[3:6], // worker 2 gets jobs 4-6
		jobs[6:9], // worker 3 gets jobs 7-9
	}

	for workerID, chunk := range chunks {
		wg.Add(1)
		go func(id int, jobs []int) {
			defer wg.Done()
			for _, job := range jobs {
				fmt.Printf("worker %d processing job %d\n", id+1, job)
				time.Sleep(200 * time.Millisecond) // simulate work
			}
		}(workerID, chunk)
	}

	wg.Wait()
}

// why not just launch 1000 goroutines for 1000 jobs?
// each goroutine = ~2-8KB stack minimum
// 1000 goroutines = potentially megabytes of memory
// Go scheduler has to manage all 1000 even if most are waiting
// 3 workers = fraction of memory, same total throughput

// -----------------------------------------------------------------------------
// 9. GOROUTINE LIFECYCLE
// -----------------------------------------------------------------------------
// goroutines are tied to the main goroutine (the process)
// when main returns, ALL goroutines are killed immediately
// there is no grace period, no cleanup, no warning

func goroutineLifecycle() {
	go func() {
		for {
			fmt.Println("running...")
			time.Sleep(500 * time.Millisecond)
			// this loop runs forever UNTIL main exits
			// at that point this goroutine is killed mid-execution if needed
		}
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("main exiting — goroutine above gets killed here")
}

// for long-running goroutines you should use context for clean shutdown:
//   ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//   defer cancel()
// the goroutine checks ctx.Done() to know when to stop gracefully

// -----------------------------------------------------------------------------
// 10. GOROUTINES IN REAL BACKEND CODE
// -----------------------------------------------------------------------------
// common patterns you'll actually use:

// background job after response — don't make user wait for email
func handleSignup( /* w, r */ ) {
	// createUser(r)
	// respondWith201(w)

	go func() {
		// sendWelcomeEmail(user)   // runs after response is sent
		// sendSlackNotification()  // user never waits for this
	}()
}

// parallel service calls — hit multiple services simultaneously
func getProfilePage(userID int) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		// getUserData(userID)
	}()
	go func() {
		defer wg.Done()
		// getUserOrders(userID)
	}()
	go func() {
		defer wg.Done()
		// getUserNotifications(userID)
	}()

	wg.Wait()
	// combine results and respond
}

func main() {
	basicGoroutine()
	withWaitGroup()
	concurrentFetch()
}

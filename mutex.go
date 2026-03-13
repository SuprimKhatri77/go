package main

import (
	"fmt"
	"sync"
	"time"
)

// SafeCounter is a concurrency-safe counter using a mutex
type SafeCounter struct {
	mu sync.Mutex     // Mutex to protect access to the map
	v  map[string]int // The actual map storing counts
}

// Inc increments the counter for a given key
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()   // LOCK the mutex before accessing/modifying the map
	c.v[key]++    // Safely increment the value
	c.mu.Unlock() // UNLOCK the mutex after done
}

// Value returns the current value of the counter for a key
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()         // LOCK the mutex before reading the map
	defer c.mu.Unlock() // UNLOCK automatically when function returns
	return c.v[key]
}

func MutexExample() {
	// Create a SafeCounter with an initialized map
	counter := SafeCounter{v: make(map[string]int)}

	// -----------------------------
	// WHY WE NEED A MUTEX:
	// -----------------------------
	// In Go, multiple goroutines can access the same variable concurrently.
	// Without a mutex, if two goroutines try to increment the map at the same time,
	// the map can get corrupted, and we could lose updates.
	// This is called a "race condition".
	//
	// The mutex ensures that only ONE goroutine can access the map at a time.
	// Other goroutines trying to access it will wait (block) until the mutex is unlocked.
	// This guarantees safe, predictable updates to shared data.
	// -----------------------------

	// Spawn 1000 goroutines that increment the same key concurrently
	// The mutex inside Inc ensures these updates happen safely
	for i := 0; i < 1000; i++ {
		go counter.Inc("somekey")
	}

	// Sleep for 1 second to allow all goroutines to finish
	// (in real applications, use sync.WaitGroup to wait precisely)
	time.Sleep(time.Second)

	// Print the value of the counter safely
	fmt.Println("counter value: ", counter.Value("somekey"))
	// The expected output is 1000
	// Without the mutex, it would likely be less than 1000 due to race conditions
}

// ============================================================
// TOPIC 13: Interfaces
// ============================================================
// Interfaces define BEHAVIOR — a set of methods a type must have.
// They are contracts: implement the methods, satisfy the interface.
//
// KEY DIFFERENCE FROM JAVA/TS:
//   Java: class Rectangle implements Shape  ← explicit declaration
//   Go:   no declaration needed             ← implicit implementation
//
// If your struct has the methods the interface requires,
// it automatically satisfies the interface — no `implements` keyword.
// This means you can satisfy an interface defined in someone else's
// package without touching their code. This is Go's superpower.
//
// Real world use case (your Gin API):
//   Instead of coupling handlers directly to GORM:
//     type TodoRepository interface {
//         GetAll() ([]Todo, error)
//         Create(todo Todo) error
//     }
//   Now your handler accepts the interface — works with real DB
//   AND with a mock for testing. Swap implementations freely.
//   This is dependency injection — exactly what Gin closures do.
// ============================================================

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

// ============================================================
// 1. BASIC INTERFACE + IMPLICIT IMPLEMENTATION
// ============================================================
// Define the contract — what methods must be present
// Any struct with these methods satisfies Shape automatically
// ============================================================
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Length  float64
	Breadth float64
}

type Circle struct {
	Radius float64
}

func (r *Rectangle) Area() float64      { return r.Length * r.Breadth }
func (r *Rectangle) Perimeter() float64 { return 2 * (r.Length + r.Breadth) }
func (c *Circle) Area() float64         { return math.Pi * math.Pow(c.Radius, 2) }
func (c *Circle) Perimeter() float64    { return 2 * math.Pi * c.Radius }

// printShapeInfo is the "orchestrator" — it doesn't care if it gets
// a Rectangle, Circle, or any future shape you add.
// It just calls the interface methods — Go figures out which
// concrete implementation to call at runtime. This is POLYMORPHISM.
//
// If you add Color() string to Shape interface without implementing
// it on Rectangle/Circle → compile error immediately. Go enforces
// the full contract.
func printShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

// ============================================================
// 2. THE STRINGER INTERFACE — built into fmt package
// ============================================================
// fmt.Println checks if value implements Stringer interface:
//
//	type Stringer interface { String() string }
//
// If yes → calls .String() automatically
// This is implicit interface in action — fmt knows nothing about
// your structs, but if they have String(), fmt uses it.
// ============================================================
type ServerConfig struct {
	Host    string
	Port    int
	IsHTTPS bool
}

type User struct {
	ID        string
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
}

func (sc *ServerConfig) String() string {
	return fmt.Sprintf("%s:%d (HTTPS: %t)", sc.Host, sc.Port, sc.IsHTTPS)
}

func (u *User) String() string {
	return fmt.Sprintf("User{%s, %s, role: %s}", u.Name, u.Email, u.Role)
}

func (u *User) IsAdmin() bool { return u.Role == "admin" }

// ============================================================
// 3. INTERFACE AS FUNCTION PARAMETER
// ============================================================
// This is where interfaces shine — one function works for any
// type that satisfies the interface. Add new notifier types
// without changing notify() or notifyAll() at all.
// ============================================================
type Notifier interface {
	Send(message string) error
	Channel() string
}

type EmailNotifier struct{ recipient string }
type SMSNotifier struct{ phoneNumber string }
type SlackNotifier struct{ webhook string }

func (e *EmailNotifier) Send(msg string) error {
	fmt.Printf("[EMAIL] to %s: %s\n", e.recipient, msg)
	return nil
}
func (e *EmailNotifier) Channel() string { return "email" }

func (s *SMSNotifier) Send(msg string) error {
	fmt.Printf("[SMS] to %s: %s\n", s.phoneNumber, msg)
	return nil
}
func (s *SMSNotifier) Channel() string { return "sms" }

func (s *SlackNotifier) Send(msg string) error {
	fmt.Printf("[SLACK] to %s: %s\n", s.webhook, msg)
	return nil
}
func (s *SlackNotifier) Channel() string { return "slack" }

func notify(n Notifier, msg string) {
	n.Send(msg)
}

func notifyAll(notifiers []Notifier, msg string) {
	for _, n := range notifiers {
		n.Send(msg) // doesn't care which type — just calls Send
	}
}

// ============================================================
// 4. INTERFACE COMPOSITION
// ============================================================
// Interfaces can embed other interfaces — builds larger contracts
// from smaller ones. Same concept as struct embedding.
// ============================================================
type Reader interface {
	Read() string
}

type Writer interface {
	Write(data string) error
}

type Closer interface {
	Close()
}

// ReadWriteCloser requires ALL methods from Reader + Writer + Closer
type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

type FileStorage struct{ Name string }
type MemoryStorage struct{ data string } // stores data in memory, no file

func (fs *FileStorage) Read() string {
	return fmt.Sprintf("reading from file: %s", fs.Name)
}
func (fs *FileStorage) Write(data string) error {
	fmt.Printf("writing '%s' to file: %s\n", data, fs.Name)
	return nil
}
func (fs *FileStorage) Close() {
	fmt.Printf("closing file: %s\n", fs.Name)
}

func (ms *MemoryStorage) Read() string {
	return fmt.Sprintf("reading from memory: %s", ms.data)
}
func (ms *MemoryStorage) Write(data string) error {
	ms.data = data // store in memory
	fmt.Printf("writing '%s' to memory\n", data)
	return nil
}
func (ms *MemoryStorage) Close() {
	fmt.Println("clearing memory storage")
	ms.data = ""
}

func useStorage(rwc ReadWriteCloser, data string) {
	rwc.Write(data)
	fmt.Println(rwc.Read())
	rwc.Close()
}

// ============================================================
// 5. TYPE ASSERTION AND TYPE SWITCH
// ============================================================
// Sometimes you have an interface value and need the concrete type.
//
// Type assertion — single type:
//
//	value, ok := interfaceVar.(ConcreteType)
//	Always use comma ok pattern — panics without it if wrong type
//
// Type switch — multiple types:
//
//	switch v := n.(type) { case EmailNotifier: ... }
//	v is the CONCRETE TYPE inside each case
//	Use v to access type-specific fields
//
// ============================================================
func describeNotifier(n Notifier) {
	// ✅ Use v := n.(type) to get concrete type + access its fields
	switch v := n.(type) {
	case *EmailNotifier:
		fmt.Printf("Email notifier → recipient: %s\n", v.recipient)
	case *SMSNotifier:
		fmt.Printf("SMS notifier → phone: %s\n", v.phoneNumber)
	case *SlackNotifier:
		fmt.Printf("Slack notifier → webhook: %s\n", v.webhook)
	default:
		fmt.Println("unknown notifier type")
	}

	// Without v — you lose access to concrete fields:
	// switch n.(type) {
	// case *EmailNotifier:
	//     fmt.Println("email") // can't access .recipient here
	// }
}

// ============================================================
// 6. EMPTY INTERFACE — interface{}  or  any
// ============================================================
// interface{} accepts ANY type — like TypeScript's `any`
// Use sparingly — you lose type safety
// Common use: JSON Data field, generic response wrapper
// ============================================================
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // any type, omit if nil
}

func respond(success bool, message string, data interface{}) *APIResponse {
	return &APIResponse{
		Success: success,
		Message: message,
		Data:    data,
	}
}

func main() {

	// Shapes — polymorphism
	printShapeInfo(&Rectangle{Length: 4, Breadth: 5})
	printShapeInfo(&Circle{Radius: 4})

	// Stringer — fmt.Println calls String() automatically
	sc := &ServerConfig{Host: "localhost", Port: 8080, IsHTTPS: false}
	u := &User{Name: "Alice", Email: "alice@example.com", Role: "admin"}
	fmt.Println(sc) // calls sc.String() automatically
	fmt.Println(u)  // calls u.String() automatically

	// Notifier — all three types in one slice
	notifiers := []Notifier{
		&EmailNotifier{recipient: "alice@example.com"},
		&SMSNotifier{phoneNumber: "+9779800000000"},
		&SlackNotifier{webhook: "https://hooks.slack.com/abc"},
	}
	notifyAll(notifiers, "System maintenance at midnight")

	// Storage — interface composition
	useStorage(&FileStorage{Name: "data.txt"}, "hello file")
	useStorage(&MemoryStorage{}, "hello memory")

	// Type switch — with v to access concrete fields
	describeNotifier(&EmailNotifier{recipient: "bob@example.com"})
	describeNotifier(&SMSNotifier{phoneNumber: "+1234567890"})
	describeNotifier(&SlackNotifier{webhook: "https://hooks.slack.com/xyz"})

	// Single type assertion — comma ok pattern
	var n Notifier = &EmailNotifier{recipient: "test@example.com"}
	if email, ok := n.(*EmailNotifier); ok {
		fmt.Println("asserted email recipient:", email.recipient)
	}

	// APIResponse with interface{} Data field
	type Todo struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	r1 := respond(true, "todo fetched", &Todo{ID: "1", Title: "Learn Go"})
	r2 := respond(true, "routes fetched", []string{"/users", "/posts", "/auth"})
	r3 := respond(false, "not found", nil) // Data omitted — omitempty + nil

	for _, r := range []*APIResponse{r1, r2, r3} {
		j, _ := json.Marshal(r)
		fmt.Println(string(j))
	}

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1.  Interfaces define behavior — a set of method signatures
	// 2.  Implementation is IMPLICIT — no `implements` keyword
	//     If struct has the methods → it satisfies the interface
	// 3.  Java: explicit (implements Shape)
	//     Go:   implicit (just have the methods)
	// 4.  One function accepting interface → works for ALL types
	//     that satisfy it. Add new types without changing the function.
	// 5.  fmt.Stringer — implement String() string → fmt.Println uses it
	// 6.  Interface composition — embed interfaces in interfaces
	// 7.  Type assertion: v, ok := iface.(ConcreteType) — always comma ok
	// 8.  Type switch: switch v := n.(type) { case T: use v }
	//     v gives you the concrete type — access type-specific fields
	// 9.  interface{} / any — accepts any type, use sparingly
	// 10. Real world: dependency injection — handler accepts interface,
	//     swap real DB vs mock for testing without changing handler code
	// --------------------------------------------------------
}

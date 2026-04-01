// ============================================================
// TOPIC 12: Structs
// ============================================================
// Structs are Go's way of grouping related data together.
// No classes, no inheritance — just data + methods.
//
// Think of it like a TypeScript type but on steroids:
// TS:  type User = { name: string; email: string }
// Go:  type User struct { Name string; Email string }
//
// Key difference: capitalized fields = exported (public)
//                 lowercase fields  = unexported (private)
// This applies to EVERYTHING in Go — fields, functions, vars, types.
// Capitalized = accessible from other packages
// Lowercase   = only accessible within the same package
// ============================================================

package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// ============================================================
// STRUCT DECLARATION
// ============================================================
type ServerConfig struct {
	Host           string    // exported — accessible from other packages
	Port           int       // exported
	MaxConnections int       // exported
	IsHTTPS        bool      // exported
	Timeout        time.Time // exported
	secretKey      string    // unexported — private to this package
}

// ============================================================
// STRUCT TAGS — metadata for JSON, DB, validation
// ============================================================
// Syntax: `json:"fieldName"` after the type
// Used by encoding/json, GORM, validator packages
//
// json:"name"          → key name in JSON output
// json:"-"             → ALWAYS exclude from JSON (sensitive fields)
//
//	field is ignored regardless of its value
//
// json:"name,omitempty"→ exclude if value is zero/empty
//
//	like TypeScript's optional fields (name?: string)
//
// db:"column_name"     → maps to database column name (used by GORM/sqlx)
// ============================================================
type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`               // never sent in JSON — security
	Token    string `json:"token,omitempty"` // omitted if empty string
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

// ============================================================
// METHODS ON STRUCTS — receivers
// ============================================================
// Syntax: func (receiverName ReceiverType) MethodName() ReturnType
// This is how you attach functions to structs — no classes needed
// ============================================================

// ============================================================
// VALUE RECEIVER vs POINTER RECEIVER — the most important distinction
// ============================================================
//
// VALUE RECEIVER: func (s ServerConfig) Method()
//   - Go creates a COPY of the struct for this method call
//   - Changes inside the method do NOT affect the original
//   - Use when: reading data, not modifying, struct is small
//   - Real world: getters, String(), IsValid(), computed properties
//
// POINTER RECEIVER: func (s *ServerConfig) Method()
//   - Method works on the ORIGINAL struct in memory
//   - Changes inside the method DO affect the original
//   - Use when: modifying struct fields, struct is large (DB models etc)
//   - Real world: setters, Update(), any method that changes state
//
// REAL WORLD EXAMPLES:
//   ✅ Pointer receiver:
//      func (u *User) UpdateEmail(email string) { u.Email = email }
//      func (t *Todo) MarkComplete() { t.Completed = true }
//      func (s *Server) IncrementConnections() { s.connections++ }
//
//   ✅ Value receiver:
//      func (u User) IsAdmin() bool { return u.Role == "admin" }
//      func (t Todo) Summary() string { return t.Title }
//      func (s ServerConfig) Address() string { return s.Host + ":" + ... }
//
// RULE OF THUMB: if in doubt, use pointer receiver.
// Go will auto-dereference value variables when calling pointer methods:
//   c := Counter{}        // value
//   c.IncrementPointer()  // Go does (&c).IncrementPointer() automatically
// ============================================================

func (sc *ServerConfig) String() string {
	return fmt.Sprintf("%s:%d (HTTPS: %t)", sc.Host, sc.Port, sc.IsHTTPS)
}

func (sc *ServerConfig) IsValid() bool {
	if sc == nil { // ✅ always nil-check in pointer receivers
		return false
	}
	return sc.Host != "" && sc.Port > 0
}

func (sc *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}

// ============================================================
// VALUE vs POINTER — concrete demonstration
// ============================================================
type Counter struct {
	count int
}

func (c Counter) IncrementValue() {
	c.count++ // modifies the COPY — original unchanged
}

func (c *Counter) IncrementPointer() {
	c.count++ // modifies the ORIGINAL — count changes in memory
}

// ============================================================
// STRUCT EMBEDDING — Go's answer to inheritance
// ============================================================
// Embed one struct inside another — no field name needed.
// Embedded struct's fields and methods are PROMOTED to outer struct.
// Access them directly: post.ID instead of post.BaseModel.ID
// ============================================================
type BaseModel struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *BaseModel) Age() string {
	return b.CreatedAt.String()
}

type Post struct {
	BaseModel // embedded — fields and methods promoted
	Title     string
	Content   string
	AuthorID  int
}

func (p *Post) Summary() string {
	return fmt.Sprintf("[%d] %s by %d", p.ID, p.Title, p.AuthorID)
	// p.ID works because ID is promoted from BaseModel
}

// ============================================================
// & vs NO & — WHEN TO USE EACH
// ============================================================
//
// Use VALUE (no &) when:
//   - Struct is small (a few fields)
//   - Used temporarily, not passed around
//   - You want an isolated copy — changes don't affect original
//   - eg: a simple config read, a quick calculation struct
//
// Use POINTER (&) when:
//   - Struct is large (DB models, request/response with many fields)
//   - Passed through multiple layers (handler → service → repository)
//   - You need to modify the original
//   - Storing in slices/maps where consistent state matters
//   - Think: 1000 todos from DB — copying each = slow server
//             pointer = just pass the address = fast
//
// In real backends: DB models, request/response structs → always &
// ============================================================

// ============================================================
// DEREFERENCING — what it actually means
// ============================================================
//
// & = "give me the ADDRESS of this" (referencing)
// * = "give me the VALUE at this address" (dereferencing)
//
// res := &UserResponse{...}
//   → creates UserResponse in memory, res holds its ADDRESS
//   → res is *UserResponse (a pointer)
//   → this is REFERENCING — you're taking the address
//
// actual := *res
//   → go to the address res holds, give me the UserResponse there
//   → actual is UserResponse (the actual struct, not a pointer)
//   → this is explicit DEREFERENCING
//
// res.Name
//   → Go automatically does (*res).Name under the hood
//   → this is IMPLICIT/AUTOMATIC DEREFERENCING
//   → Go hides it so you don't write (*res).Name every time
//   → every time you do ptr.Field, Go is dereferencing silently
//
// So in your code — res.Name, todo.Title, post.ID on pointers
// were ALL dereferencing — Go just did it automatically for you.
// ============================================================

// ============================================================
// REAL WORLD: Todo API structs
// ============================================================
type Todo struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UserID      string
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Todo    *Todo  `json:"todo,omitempty"`
}

type User struct {
	ID        string
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin" // no if needed — boolean expression IS the return
}

func NewTodo(req *CreateTodoRequest, userID string) *Todo {
	return &Todo{
		ID:          "todo_1",
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}
}

func ToResponse(t *Todo) *TodoResponse {
	return &TodoResponse{
		Success: true,
		Message: "Todo created successfully",
		Todo:    t,
	}
}

func main() {

	// Three ways to initialize a struct
	named := &ServerConfig{
		Host:           "localhost",
		Port:           8080,
		MaxConnections: 100,
		IsHTTPS:        false,
	}

	positional := &ServerConfig{"localhost", 8080, 100, false, time.Now(), ""}

	var empty ServerConfig // zero values — all fields at zero value of their type
	empty.Host = "localhost"
	empty.Port = 9000

	fmt.Println(named.String())
	fmt.Println(positional.String())
	fmt.Println(empty.String())
	fmt.Println(named.IsValid())
	fmt.Println(named.Address())

	// Value vs pointer receiver
	c := Counter{}
	c.IncrementValue()
	fmt.Println("after value receiver:", c.count) // 0 — copy was modified

	cp := &Counter{}
	cp.IncrementPointer()
	fmt.Println("after pointer receiver:", cp.count) // 1 — original modified

	// Auto-dereference — Go handles this automatically
	c.IncrementPointer()                // Go does (&c).IncrementPointer() — works fine
	fmt.Println("auto-deref:", c.count) // 1

	// Embedding
	post := &Post{
		BaseModel: BaseModel{ // must name it when initializing
			ID:        77,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Title:    "My First Post",
		Content:  "Hello Go",
		AuthorID: 7,
	}
	fmt.Println(post.ID)        // promoted from BaseModel — no post.BaseModel.ID needed
	fmt.Println(post.Summary()) // [77] My First Post by 7
	fmt.Println(post.Age())     // method promoted from BaseModel

	// JSON tags
	res := &UserResponse{
		ID:       "xyz",
		Name:     "suprim77",
		Email:    "suprim@example.com",
		Success:  true,
		Message:  "User created",
		Password: "secret123", // json:"-" — will NOT appear in output
		Token:    "",          // omitempty — will NOT appear since empty
	}
	out, _ := json.Marshal(res)
	fmt.Println(string(out))

	// Todo API
	req := &CreateTodoRequest{
		Title:       "Learn Go Structs",
		Description: "Understand value vs pointer receivers",
	}
	todo := NewTodo(req, "user_123")
	response := ToResponse(todo)
	jsonOut, _ := json.Marshal(response)
	fmt.Println(string(jsonOut))

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1.  Capitalized = exported (public), lowercase = unexported (private)
	//     Applies to fields, methods, functions, types — everything
	// 2.  Three init ways: named fields ✅, positional, empty then assign
	//     Always use named fields — positional breaks when struct changes
	// 3.  Zero value of struct = zero values of all its fields
	// 4.  Value receiver → copy → changes don't affect original → use for reads
	// 5.  Pointer receiver → original → changes affect original → use for writes
	// 6.  Rule of thumb: when in doubt use pointer receiver
	// 7.  Go auto-dereferences: c.Method() works even if Method needs *Counter
	// 8.  Always nil-check at top of pointer receiver methods
	// 9.  Embedding promotes fields AND methods to outer struct
	// 10. When initializing embedded struct, name it explicitly: BaseModel: BaseModel{}
	// 11. json:"-"          → always exclude (passwords, internal fields)
	//     json:"omitempty"  → exclude if zero/empty (optional fields like Token)
	// 12. Use & (pointer) for large structs passed through multiple layers
	//     Use value for small, temporary, read-only structs
	// 13. & = referencing (take address), * = dereferencing (get value at address)
	//     res.Field on a pointer = Go auto-dereferences silently every time
	// --------------------------------------------------------
}

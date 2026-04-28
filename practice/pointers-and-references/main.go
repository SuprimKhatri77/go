// ============================================================
// TOPIC 16: Pointers and References
// ============================================================
// Two operators, one rule:
//   & — give me the ADDRESS of this value (referencing)
//   * — give me the VALUE at this address (dereferencing)
//
// You already use pointers — every &Todo{}, *User, return &User{}
// in your Gin API was pointers. This makes the mechanics explicit.
// ============================================================

package main

import "fmt"

type User struct {
	id    int
	email string
	name  string
	role  string
}

type Todo struct {
	title       string
	isCompleted bool
}

func main() {

	// --------------------------------------------------------
	// 1. THE BASICS
	// --------------------------------------------------------
	// & takes the address of a value — gives you a pointer
	// * dereferences a pointer — gives you the value at that address
	// --------------------------------------------------------
	x := 42
	p := &x // p is *int — holds the memory address of x

	fmt.Println(p)  // 0xc000... — the memory address
	fmt.Println(*p) // 42 — the value AT that address (dereferencing)

	*p = 100       // go to the address p holds, write 100 there
	fmt.Println(x) // 100 — x changed because p points to x

	// --------------------------------------------------------
	// 2. POINTER VS VALUE IN FUNCTIONS
	// --------------------------------------------------------
	// Value parameter — Go copies the value into the function.
	// Changes inside the function do NOT affect the original.
	//
	// Pointer parameter — Go passes the address.
	// Changes inside the function DO affect the original.
	// --------------------------------------------------------
	n := 7
	doubleValue(n) // works on a copy — n unchanged
	fmt.Println(n) // 7

	doublePointer(&n) // works on original — n changed
	fmt.Println(n)    // 14

	// Real backend version
	user := User{email: "before@example.com"}

	updateEmail(user, "after@example.com") // value — copy modified
	fmt.Println(user.email)                // "before@example.com" — unchanged

	updateEmailPtr(&user, "after@example.com") // pointer — original modified
	fmt.Println(user.email)                    // "after@example.com" — changed

	// --------------------------------------------------------
	// 3. NIL POINTERS
	// --------------------------------------------------------
	// var u *User — u is nil. pointer declared but no User in memory.
	// reading u.Name on a nil pointer = panic: nil pointer dereference
	// always nil-check before accessing fields on a pointer
	// --------------------------------------------------------
	var u *User
	fmt.Println(u) // <nil>
	// fmt.Println(u.name) // panic: nil pointer dereference

	fmt.Println(getUserName(nil))                  // "" — safe
	fmt.Println(getUserName(&User{name: "alice"})) // "alice"

	// --------------------------------------------------------
	// 4. POINTERS TO STRUCTS
	// --------------------------------------------------------
	// Two ways to create a struct:
	//   u1 := User{...}   — value, stored directly
	//   u2 := &User{...}  — pointer, stored as address to User in memory
	//
	// &User{} doesn't modify anything — it creates a User in memory
	// and gives you its address instead of copying the value.
	// There's no "original" floating separately — the struct exists
	// once, & just gives you the address to that one place.
	//
	// AUTO-DEREFERENCING:
	// When you do u2.Name on a *User, Go automatically does (*u2).Name
	// You never need to write (*u).field — Go handles it silently.
	// Every u.Field on a pointer is implicit dereferencing.
	// --------------------------------------------------------
	u2 := &User{name: "bob"}
	activateUser(u2)
	fmt.Println(u2.email) // "activated@example.com" — modified via pointer

	created := createUser("alice", "alice@example.com")
	fmt.Println(created.name) // "alice" — *User returned, auto-deref on .name

	// --------------------------------------------------------
	// 5. *[]Todo vs []Todo — when to use which
	// --------------------------------------------------------
	// []Todo is ALREADY a reference type.
	// It has a pointer inside it (pointer + length + cap = 24 bytes).
	// Passing []Todo to a function is cheap — you're not copying todos.
	//
	// *[]Todo = pointer TO the slice header itself.
	// Almost never needed EXCEPT when something needs to ASSIGN
	// a completely new slice to your variable from outside.
	//
	// The database/sql Scan case:
	//   var todos []Todo
	//   db.Scan(&todos)  // Scan needs *[]Todo to assign result to todos
	//                    // without &, Scan writes to a copy, todos stays empty
	//
	// In every other case — just use []Todo.
	// *[]Todo is not about memory efficiency — []Todo is already cheap.
	// --------------------------------------------------------

	// just use []Todo — clean and idiomatic
	todos := getTodoSlice()
	for _, v := range todos {
		fmt.Println(v.title)
	}

	// *[]Todo — must explicitly dereference to range over it
	todoPtr := getTodosPtr()
	for _, v := range *todoPtr { // must dereference with * to range
		fmt.Println(v.title)
	}

	// --------------------------------------------------------
	// 6. REAL WORLD POINTER PATTERNS
	// --------------------------------------------------------

	users := []User{
		{id: 1, name: "alice", email: "alice@example.com", role: "user"},
		{id: 2, name: "bob", email: "bob@example.com", role: "guest"},
		{id: 3, name: "charlie", email: "charlie@example.com", role: "user"},
	}

	// findUser returns *User — nil means not found
	found := findUser(2, users)
	if found == nil {
		fmt.Println("user not found")
		return
	}

	updateUserRole(found, "admin")
	fmt.Println(getUserEmail(found)) // "bob@example.com"

	// findUser bug — common Go gotcha
	// for _, v := range users { return &v }
	// v is a COPY created by range — temporary variable
	// &v returns address of the copy, not the actual slice element
	// after iteration, v may be overwritten — stale pointer
	//
	// fix: use index
	// for i := range users { return &users[i] } // address of actual element

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1.  & = take the address (referencing)   — result is *T
	//     * = get value at address (dereferencing) — result is T
	// 2.  Value param → copy → changes don't affect original
	//     Pointer param → address → changes affect original
	// 3.  var u *User → nil pointer. reading u.field = panic
	//     always nil-check: if u == nil { return }
	// 4.  u.Field on *User = Go auto-dereferences silently
	//     you never need to write (*u).Field
	// 5.  &User{} — creates User in memory, gives you address
	//     no "original" exists separately — struct is in one place
	// 6.  []Todo is already a reference type — cheap to pass
	//     *[]Todo needed only when something assigns a new slice to your var
	//     (database/sql Scan pattern) — not for memory efficiency
	// 7.  for _, v := range slice { return &v } — BUG
	//     v is a copy — use for i := range slice { return &slice[i] }
	// 8.  passing nil to a pointer receiver method = panic
	//     guard: if u == nil { return } at top of method
	// 9.  pointer size is always 8 bytes on 64-bit systems
	//     regardless of how large the struct is
	// 10. use *T for: large structs passed across layers, optional
	//     values (nil = not found), mutation across function boundaries
	//     use T  for: small structs, pure reads, no mutation needed
	// --------------------------------------------------------
}

func doubleValue(n int)    { n = n * 2 }   // copy — original unchanged
func doublePointer(n *int) { *n = *n * 2 } // original modified

func updateEmail(u User, email string)     { u.email = email } // copy
func updateEmailPtr(u *User, email string) { u.email = email } // original

func getUserName(u *User) string {
	if u == nil {
		return ""
	} // nil check — never panic
	return u.name
}

func getUserEmail(u *User) string {
	if u == nil {
		return ""
	}
	return u.email
}

func createUser(name, email string) *User {
	return &User{name: name, email: email}
}

func activateUser(u *User) {
	u.email = "activated@example.com" // auto-deref — Go does (*u).email
}

func getTodoSlice() []Todo {
	return []Todo{
		{title: "todo 1", isCompleted: false},
		{title: "todo 2", isCompleted: true},
	}
}

func getTodosPtr() *[]Todo {
	return &[]Todo{
		{title: "todo 3", isCompleted: false},
		{title: "todo 4", isCompleted: true},
	}
}

// BUG VERSION — don't do this
// func findUser(id int, users []User) *User {
//     for _, v := range users {
//         if v.id == id {
//             return &v  // &v = address of loop copy — stale pointer
//         }
//     }
//     return nil
// }

// CORRECT VERSION — address of actual slice element
func findUser(id int, users []User) *User {
	for i := range users {
		if users[i].id == id {
			return &users[i] // address of actual element in slice
		}
	}
	return nil
}

func updateUserRole(u *User, role string) {
	// updateUserRole(nil, "admin") would panic here
	// guard if needed: if u == nil { return }
	u.role = role
}

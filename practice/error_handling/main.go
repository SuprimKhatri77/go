// ============================================================
// TOPIC 14: Error Handling
// ============================================================
// In Go, errors are values — not exceptions.
// No try/catch, no throwing. You return errors and check them.
// Explicit, predictable, and impossible to accidentally ignore.
//
// The built-in error interface is literally just:
//   type error interface {
//       Error() string
//   }
// Any type with Error() string satisfies it — implicit implementation.
// ============================================================

package main

import (
	"errors"
	"fmt"
	"strings"
)

// ============================================================
// 1. RETURNING AND HANDLING ERRORS
// ============================================================
// Convention: error is always the last return value.
// Return nil for no error, fmt.Errorf for errors.
// Error messages lowercase — they get chained mid-sentence.
// ============================================================
func readConfig(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty") // lowercase
	}
	if !strings.HasSuffix(path, ".yaml") {
		return "", fmt.Errorf("invalid extension, expected .yaml got %s", path)
	}
	return "config data", nil
}

func connectDB(connString string) (string, error) {
	if !strings.HasPrefix(connString, "postgres://") {
		return "", fmt.Errorf("invalid connection string, must start with postgres://")
	}
	return "connected", nil
}

// ============================================================
// 2. CUSTOM ERROR TYPES
// ============================================================
// When you need structured errors — HTTP status codes, error codes,
// machine-readable context — define a custom type.
//
// To satisfy the error interface, implement Error() string.
// Use a pointer receiver (*AppError) — the pointer is what satisfies
// the interface, and it's what errors.As looks for in the chain.
// ============================================================
type AppError struct {
	Code       int
	Message    string
	StatusCode int
}

func (ae *AppError) Error() string {
	return fmt.Sprintf("[%d] %s (http %d)", ae.Code, ae.Message, ae.StatusCode)
}

type User struct{ ID int }

func getUser(id int) (*User, error) {
	switch {
	case id == 1 || id == 2:
		return &User{ID: id}, nil
	case id <= 0:
		return nil, &AppError{Code: 400, Message: "invalid id", StatusCode: 400}
	default:
		return nil, &AppError{Code: 404, Message: "user not found", StatusCode: 404}
	}
}

// ============================================================
// 3. ERROR WRAPPING WITH %w
// ============================================================
// Errors travel through layers — db → repository → service → handler.
// Each layer adds context without losing the original error.
//
// %w  = wraps the error, chain is preserved, can be unwrapped
// %v  = formats as string only, original is LOST, chain broken
//
// Use %w when callers need to check what's inside.
// Use %v when you just want the message as a string.
//
// ALWAYS guard before wrapping — %w on nil gives %!w(<nil>)
//
// HOW THE CHAIN WORKS:
//
//	dbQuery returns:    "record not found"
//	repository wraps:   "repository.GetUser: record not found"
//	service wraps:      "service.GetUser: repository.GetUser: record not found"
//
//	errors.Is walks the entire chain — finds the original inside
//	errors.Unwrap peels one layer at a time
//
// ============================================================
func dbQuery(id int) error {
	if id > 100 {
		return fmt.Errorf("record not found for id %d", id)
	}
	return nil
}

func repository(id int) error {
	err := dbQuery(id)
	if err == nil {
		return nil // guard — never wrap nil
	}
	return fmt.Errorf("repository.GetUser: %w", err) // %w preserves chain
}

func service(id int) error {
	err := repository(id)
	if err == nil {
		return nil
	}
	return fmt.Errorf("service.GetUser: %w", err) // wrap again, chain grows
}

// ============================================================
// 4. errors.As — EXTRACT TYPED ERROR FROM CHAIN
// ============================================================
// errors.As walks the error chain looking for a specific TYPE.
// When found, it writes the error into your target variable.
//
// HOW THE POINTER MECHANICS WORK:
//   var appErr *AppError     → appErr is nil right now, no AppError in memory
//   errors.As(err, &appErr)  → walks chain, finds *AppError, makes appErr POINT TO IT
//   appErr.StatusCode        → now accessible — appErr points to the struct in err
//
// The &AppError{} returned by getUser is one struct in memory.
// var appErr *AppError is a SEPARATE nil pointer.
// errors.As connects them — makes appErr point to what's in the chain.
// Nothing is overwritten — appErr just stops being nil.
//
// Why pointer receiver matters:
//   AppError{}  with value receiver → doesn't satisfy error interface → compile error
//   &AppError{} with pointer receiver → satisfies error → errors.As can find it
// ============================================================

// ============================================================
// 5. SENTINEL ERRORS
// ============================================================
// Predefined error values at package level.
// Check with errors.Is — never compare error strings directly.
// errors.Is walks the chain so it works even when wrapped with %w.
// ============================================================
var ErrUnauthorized = errors.New("unauthorized")
var ErrExpired = errors.New("token expired")
var ErrInvalidUserID = errors.New("invalid user id")

func authenticate(token string) error {
	switch token {
	case "valid-token":
		return nil
	case "expired-token":
		return ErrExpired
	case "":
		return ErrUnauthorized
	default:
		return fmt.Errorf("authenticate: %w", ErrUnauthorized) // wrapped sentinel
	}
}

// ============================================================
// 6. REAL WORLD — withErrorHandling MIDDLEWARE PATTERN
// ============================================================
// HandlerFunc is a type alias — same pattern as Topic 11 middleware.
// withErrorHandling wraps ANY HandlerFunc with error classification.
//
// HOW IT WORKS — same as logger(greetHandler) from Topic 11:
//  1. withErrorHandling takes a HandlerFunc (the real handler)
//  2. returns a NEW HandlerFunc that wraps the original
//  3. the wrapper calls the original, checks the error, classifies it
//  4. caller calls the wrapper, not the original
//
// CALLING IT:
//
//	wrong:   withErrorHandling(getProfile(25))  — passes (*User, error), not HandlerFunc
//	correct: wrapped := withErrorHandling(getProfile) — pass the func itself
//	         user, err := wrapped(25)            — then call the wrapper
//
// ============================================================
type HandlerFunc func(userID int) (*User, error)

func withErrorHandling(h HandlerFunc) HandlerFunc {
	return func(userID int) (*User, error) {
		user, err := h(userID) // call the actual handler
		if err == nil {
			return user, nil // no error — pass through
		}

		// classify the error — check most specific first
		var appErr *AppError
		if errors.As(err, &appErr) {
			// structured error — has status code, code, message
			fmt.Printf("app error — status %d: %s\n", appErr.StatusCode, appErr.Message)
		} else if errors.Is(err, ErrUnauthorized) {
			fmt.Println("sentinel: unauthorized")
		} else if errors.Is(err, ErrInvalidUserID) {
			fmt.Println("sentinel: invalid user id")
		} else {
			fmt.Println("generic error:", err)
		}

		return nil, err
	}
}

// handlers must match HandlerFunc signature: func(userID int) (*User, error)
func getProfile(userID int) (*User, error) {
	if userID > 20 {
		return nil, &AppError{Code: 400, Message: "bad request", StatusCode: 400}
	}
	return &User{ID: userID}, nil
}

func getPosts(userID int) (*User, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	return &User{ID: userID}, nil
}

func main() {

	// basic error handling
	config, err := readConfig("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config)

	status, err := connectDB("postgres://localhost/mydb")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(status)

	// custom error type
	user, err := getUser(5)
	if err != nil {
		fmt.Println(err.Error()) // calls AppError.Error()
	}
	fmt.Println(user)

	// error wrapping — %w
	err = service(150)
	fmt.Println(err)
	// "service.GetUser: repository.GetUser: record not found for id 150"

	// %w vs %v
	original := errors.New("database timeout")
	withW := fmt.Errorf("service: %w", original)
	withV := fmt.Errorf("service: %v", original)
	fmt.Println(errors.Is(withW, original)) // true  — chain intact
	fmt.Println(errors.Is(withV, original)) // false — chain broken, just a string

	// errors.As — extract typed error from chain
	var appErr *AppError
	_, err = getUser(5)
	wrapped := fmt.Errorf("handler: %w", err)
	if errors.As(wrapped, &appErr) {
		// appErr was nil, errors.As made it point to the *AppError in the chain
		fmt.Println("status code:", appErr.StatusCode) // 404
	}

	// sentinel errors
	err = authenticate("expired-token")
	if errors.Is(err, ErrExpired) {
		fmt.Println("token is expired — ask user to re-login")
	}

	err = authenticate("")
	if errors.Is(err, ErrUnauthorized) {
		fmt.Println("no token provided")
	}

	// withErrorHandling — middleware pattern
	profileHandler := withErrorHandling(getProfile) // pass func, not result
	_, err = profileHandler(25)                     // call the wrapper

	postsHandler := withErrorHandling(getPosts)
	_, err = postsHandler(0)

	// --------------------------------------------------------
	// KEY RULES TO REMEMBER
	// --------------------------------------------------------
	// 1.  error is a built-in interface: type error interface { Error() string }
	// 2.  any type with Error() string satisfies error — implicit implementation
	// 3.  always check: if err != nil { handle it }
	// 4.  error is always the last return value by convention
	// 5.  error messages lowercase — they get chained mid-sentence
	// 6.  fmt.Errorf("context: %w", err) — wraps, chain preserved
	//     fmt.Errorf("context: %v", err) — string only, chain broken
	// 7.  always guard before wrapping — %w on nil = %!w(<nil>)
	// 8.  errors.Is(err, target) — checks entire chain for exact value
	//     use for sentinel errors
	// 9.  errors.As(err, &target) — checks entire chain for a TYPE
	//     writes found error into target — target must be a pointer
	//     var appErr *AppError → nil, errors.As makes it point to found error
	// 10. sentinel errors — package level var, check with errors.Is
	//     var ErrNotFound = errors.New("not found")
	// 11. custom error types — implement Error() string with pointer receiver
	//     return &AppError{} not AppError{} — pointer satisfies the interface
	// 12. withErrorHandling pattern:
	//     pass the function not the result: withErrorHandling(getProfile)
	//     call the wrapper: wrapped(userID)
	// --------------------------------------------------------
}

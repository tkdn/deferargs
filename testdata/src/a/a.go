package a

import (
	"errors"
	"fmt"
	"os"
)

// Basic variable reference case
func fn() (err error) {
	defer dumpError(err) // want "defer with variable argument\\(s\\): value\\(s\\) are evaluated immediately; wrap in a closure if a later update is intended"
	return nil
}

// Multiple variable arguments
func multipleVars() {
	x := 1
	y := 2
	defer printf("x=%d, y=%d", x, y) // want "defer with variable argument\\(s\\): value\\(s\\) are evaluated immediately; wrap in a closure if a later update is intended"
	x = 10
	y = 20
}

// Channel variable
func channelVar() {
	ch := make(chan int)
	defer close(ch) // want "defer with variable argument\\(s\\): value\\(s\\) are evaluated immediately; wrap in a closure if a later update is intended"
}

// Struct field access
type MyStruct struct {
	Field string
}

func structField() {
	obj := &MyStruct{Field: "initial"}
	defer printf("field: %s", obj.Field) // want "defer with variable argument\\(s\\): value\\(s\\) are evaluated immediately; wrap in a closure if a later update is intended"
	obj.Field = "changed"
}

// Mixed arguments - variable and literal
func mixedArgs() {
	msg := "error"
	defer printf("status: %s, code: %d", msg, 404) // want "defer with variable argument\\(s\\): value\\(s\\) are evaluated immediately; wrap in a closure if a later update is intended"
	msg = "success"
}

// Variable from outer scope
var globalVar = "global"

func outerScopeVar() {
	localVar := "local"
	defer printf("vars: %s, %s", globalVar, localVar) // want "defer with variable argument\\(s\\): value\\(s\\) are evaluated immediately; wrap in a closure if a later update is intended"
}

// Function return value assigned to variable
func returnValue() {
	file, err := os.Open("test.txt")
	if file != nil {
		defer file.Close()
	}
	defer dumpError(err) // want "defer with variable argument\\(s\\): value\\(s\\) are evaluated immediately; wrap in a closure if a later update is intended"
}

// Cases that should NOT be reported

// No arguments
func noArgs() {
	defer cleanup()
}

// Function call as argument (not variable)
func functionCall() {
	defer dumpError(errors.New("boom"))
}

// Literal arguments
func literals() {
	defer printf("literal: %s", "hello")
	defer printf("number: %d", 42)
	defer dumpError(nil)
}

// Wrapped in closure
func wrappedInClosure() {
	err := errors.New("test")
	defer func() {
		dumpError(err) // This is fine - wrapped in closure
	}()
	err = errors.New("changed")
}

// Anonymous function with variable (should NOT be reported)
func anonymousFunc() {
	x := 1
	defer func(val int) {
		fmt.Println(val)
	}(x)
	x = 2
}

// Helper functions
func dumpError(err error) {
	fmt.Printf("err: %v", err)
}

func printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func cleanup() {
	fmt.Println("cleanup")
}

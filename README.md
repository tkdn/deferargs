# deferargs

<!-- [![Go Report Card](https://goreportcard.com/badge/github.com/yourname/deferargs)](https://goreportcard.com/report/github.com/yourname/deferargs) -->

`deferargs` is a static analysis tool to detect `defer` statements that call a function with variable arguments.  
Since arguments in `defer` calls are evaluated immediately, passing variables may lead to unexpected behavior if their values change afterward.

## Problem

In Go, the arguments to a `defer` call are evaluated **at the point of the `defer` statement**, not when the deferred function is executed.  
This can cause subtle bugs when passing a variable whose value is expected to change later.

For example:

```go
func fn() (err error) {
    defer dump(err) // ❗️ evaluated now (likely nil), not at the end
    err = errors.New("something failed")
    return
}
```

The call to `dump(err)` will receive `nil`, because `err` is evaluated at the point of the `defer`, before `err` is set.

## Recommendation

To avoid this issue, wrap the call in a closure:

```go
func f() (err error) {
    defer func() {
        dump(err) // ✅ evaluated at the time of execution
    }()
    err = errors.New("something failed")
    return
}
```

## Detected Code

This tool reports `defer` statements where:

- The deferred function is **not** an anonymous function (`func() { ... }`)
- At least one argument is a **variable reference** (`x`, `obj.Field`, etc.)

### Reported:

```go
defer fn(err)
defer close(ch)
defer logf("fail: %v", err)
```

### Not Reported:

```go
defer fn()                           // no arguments
defer fn(errors.New("boom"))        // function call
defer fn(nil)                       // literal
defer func() { fn(err) }()          // wrapped in closure
```

## Installation

```bash
go install github.com/tkdn/deferargs/cmd/deferargs@latest
```

## Usage

### As `go vet` plugin (with [staticcheck.io](https://staticcheck.io/)):

```bash
go vet -vettool=$(which deferargs) ./...
```

### Standalone

```bash
deferargs ./...
```

## License

MIT

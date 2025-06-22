package a

import "fmt"

func fn() (err error) {
	defer dumpError(err) // want "defer with variable argument\\(s\\): value\\(s\\) are evaluated immediately; wrap in a closure if a later update is intended"
	return nil
}

func dumpError(err error) {
	fmt.Printf("err: %v", err)
}

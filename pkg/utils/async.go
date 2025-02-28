package utils

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

func EnsureRunGoroutine(f func(), tryCount ...int) {
	try := 0
	if len(tryCount) > 0 {
		try = tryCount[0]
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic: %v", r)
				fmt.Printf("%s", string(debug.Stack()))
				time.Sleep(1 * time.Second)
				if try > 10 {
					os.Exit(1)
				}
				EnsureRunGoroutine(f, try+1)
			}
		}()

		f()
	}()
}

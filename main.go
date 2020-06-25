package main

import (
	fu "sock/funcs"
	"fmt"
)

func main() {
	opts := fu.Args()

	opts.CheckLock()
	pass := opts.Lock()

	if err := opts.Start(pass); err != nil {
		fmt.Println("sock: couldn't initialize screen")
	}
}

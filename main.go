package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) <= 2 {
		server(args)
		return
	} else if len(args) > 2 {
		client(args)
		return
	}
	fmt.Fprintln(os.Stderr, "Not enough arguments")
	os.Exit(1)
}
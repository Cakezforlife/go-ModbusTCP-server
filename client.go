package main

import (
	"os"
	"fmt"
)

func main() {
	args := os.Args
	if len(args) < 3{
		fmt.Println("Expected host and port")
		fmt.Println("Usage\n" +
					"./client ip port\n" +
					" - ip:\t\tIPv4 of Host\n" +
					" - port:\tPort of Host")
		os.Exit(1)
	}
	host := fmt.Sprintf("%s:%s", args[1], args[2])
	fmt.Println(host)
}

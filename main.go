package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"tcpserver/types"
	"tcpserver/helper"
)

func main() {
	args := os.Args
	hostport := ":502"
	if len(args) > 1{
		hostport = ":" + args[1]
	}
	l, err := net.Listen("tcp", hostport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Listening on Port: %s\n%v\n", hostport, err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error Accepting Connection\n%v\n", err)
			os.Exit(1)
		}
		for {
			data := make([]byte, 260)
			size, err := c.Read(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error Recieving Data\n%v\n", err)
				break
			}
			if size <= 0 {
				fmt.Fprintf(os.Stderr, "Recieved no Data\n%v\n", data)
				break
			}
	
			t := time.Now()
			time := t.Format(time.ANSIC)
			fmt.Printf("Recieved at %s\n%s\n", time, helper.FormatByteSliceAsHexSliceString(data))
	
			ADU, err := types.ParseModbusADU(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error Parsing Data: %v\n", err)
				break
			}
			fmt.Printf("%v\n", ADU.ToString())
			c.Write(data)
		}
	}
}
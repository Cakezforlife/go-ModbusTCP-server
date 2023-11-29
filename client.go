package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	"tcpserver/helper"
	"tcpserver/types"
)

var NUM_GOROUTINES = 50

var wg sync.WaitGroup

func dos(buffer []byte, c net.Conn){
	runtime.LockOSThread()
	for {
		//fmt.Fprintf(os.Stdout, "Printing!\n")
		m, err := c.Write(buffer)
		if err != nil {
			fmt.Printf("%d Error %s\n", m, err.Error())
			break
		}
		time.Sleep(1*time.Nanosecond)
	}
	runtime.UnlockOSThread()
	wg.Done();
}

func client(args []string) {
	if len(args) < 3{
		fmt.Println("Expected host and port")
		fmt.Println("Usage\n" +
					"./client ip port\n" +
					" - ip:\t\tIPv4 of Host\n" +
					" - port:\tPort of Host")
		os.Exit(1)
	}
	host := fmt.Sprintf("%s:%s", args[1], args[2])
	fmt.Println("Starting Client connection to", host)

	c, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to host %v\n", err)
		os.Exit(1)
	}

	responseMBAP := types.MBAPHeader{
		TransactionIdentifier: 1,
		ProtocolIdentifier: 0,
		Length: 6,
		UnitIdentifier: 0x01,
	}
	responsePDU := types.ModbusPDU {
		FunctionCode: 5,
		Data: []byte{0x40,0x00,0x00,0x00},
	}
	responseADU := types.ModbusADU{MBAP: responseMBAP, PDU: responsePDU}
	fmt.Printf("Sending: %s\n", responseADU.ToString())

	responseBytes, err := responseADU.ToBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ToBinary Failed %v", err)
	}
	fmt.Printf("%v\n\n", helper.FormatByteSliceAsHexSliceString(responseBytes))

	if args[3] == "dos" {
		wg.Add(NUM_GOROUTINES)
	
		oldsetting := runtime.GOMAXPROCS(12)
		fmt.Printf("OLD # of cores %d\n", oldsetting);

		for i := 0; i < NUM_GOROUTINES; i++ {
			go dos(responseBytes, c)
			fmt.Printf("Started Thread %d!!\n", i)
		}

		wg.Wait()
	} else {
		c.Write(responseBytes)

		data := make([]byte, 260) // ModbusTCP frames have a cap of 260 bytes
		size, err := c.Read(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error Recieving Data\n%v\n", err)
			os.Exit(1)
		}
		if size <= 0 {
			fmt.Fprintf(os.Stderr, "Recieved no Data\n%v\n", data)
			os.Exit(1)
		}

		t := time.Now()	// prints time to console
		time := t.Format(time.ANSIC)
		fmt.Printf("Recieved at %s\n%s\n", time, helper.FormatByteSliceAsHexSliceString(data))

		ADU, err := types.ParseModbusADU(data)	// sends data to function that returns decoded MODBUS frame
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error Parsing Data: %v\n", err)
		}
		fmt.Printf("%v\n", ADU.ToString())
	}


}

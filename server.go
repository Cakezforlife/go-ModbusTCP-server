package main

import (
	"fmt"
	"net"
	"os"
	"tcpserver/helper"
	"tcpserver/types"
	"time"
)

func server(args []string) {
	hostport := ":502"
	if len(args) > 1{
		hostport = ":" + args[1]
	}
	fmt.Printf("Starting Server on localhost:%s\n", hostport)

	// Listen on socket at port {hostport}
	l, err := net.Listen("tcp", hostport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Listening on Port: %s\n%v\n", hostport, err)
		os.Exit(1)
	}
	defer l.Close()

	// infinite loop to accept connections and split them into goroutines
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error Accepting Connection\n%v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Connection From: %v\n", c.RemoteAddr())
		go connectionRoutine(c)
	}
}

/*
	Function for use in goroutines
	Takes in a net connection and reads incoming data and sends either an echo or a legitimate response
*/
func connectionRoutine(c net.Conn) {
	for {
		data := make([]byte, 260) // ModbusTCP frames have a cap of 260 bytes
		size, err := c.Read(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error Recieving Data\n%v\n", err)
			break
		}
		if size <= 0 {
			fmt.Fprintf(os.Stderr, "Recieved no Data\n%v\n", data)
			break
		}

		t := time.Now()	// prints time to console
		time := t.Format(time.ANSIC)
		fmt.Printf("Recieved at %s\n%s\n", time, helper.FormatByteSliceAsHexSliceString(data))

		ADU, err := types.ParseModbusADU(data)	// sends data to function that returns decoded MODBUS frame
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error Parsing Data: %v\n", err)
			break
		}
		fmt.Printf("%v\n", ADU.ToString())

		// Checks if the FunctionCode is 2 -- Resquest Input
		if (ADU.PDU.FunctionCode == 2){
			responseMBAP := types.MBAPHeader{
				TransactionIdentifier: ADU.MBAP.TransactionIdentifier,
				ProtocolIdentifier: ADU.MBAP.ProtocolIdentifier,
				Length: 4,
				UnitIdentifier: 2,
			}
			responsePDU := types.ModbusPDU {
				FunctionCode: 2,
				Data: []byte{1,1},
			}
			responseADU := types.ModbusADU{MBAP: responseMBAP, PDU: responsePDU}
			fmt.Printf("Sending: %s\n", responseADU.ToString())

			responseBytes, err := responseADU.ToBinary()
			if err != nil {
				fmt.Fprintf(os.Stderr, "ToBinary Failed %v", err)
			}
			fmt.Printf("%v\n", helper.FormatByteSliceAsHexSliceString(responseBytes))

			c.Write(responseBytes)
		}else{
			c.Write(data)
		}
		
	}
}
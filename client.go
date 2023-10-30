package main

import (
	"os"
	"fmt"
	"net"
)

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

	c, err := net.Dial("TCP", host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to host %v\n", err)
	}

	responseMBAP := types.MBAPHeader{
		TransactionIdentifier: ADU.MBAP.TransactionIdentifier,
		ProtocolIdentifier: ADU.MBAP.ProtocolIdentifier,
		Length: 4,
		UnitIdentifier: 2,
	}
	responsePDU := types.ModbusPDU {
		FunctionCode: 6,
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
}

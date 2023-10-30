package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"tcpserver/helper"
)
 
type MBAPHeader struct {
	TransactionIdentifier int16
	ProtocolIdentifier    int16
	Length                int16
	UnitIdentifier        int8
}

type ModbusPDU struct {
	FunctionCode byte
	Data []byte
}

type ModbusADU struct {
	MBAP        MBAPHeader
	PDU			ModbusPDU
}

func (MBAP *MBAPHeader) ToString() string {
	var out string
	out += "MBAP Header\n"
	out += fmt.Sprintf("\tTransaction Identifier: %v\n", MBAP.TransactionIdentifier)
	out += fmt.Sprintf("\tProtocol Identifier: %v\n", MBAP.ProtocolIdentifier)
	out += fmt.Sprintf("\tLength: %v\n", MBAP.Length)
	out += fmt.Sprintf("\tUnit Identifier: %v\n", MBAP.UnitIdentifier)
	return out
}

func (PDU *ModbusPDU) ToString() string {
	var out string
	out += "Modbus PDU\n"
	out += fmt.Sprintf("\tFunction Code: %v\n", PDU.FunctionCode)
	out += fmt.Sprintf("\tData: %s\n", helper.FormatByteSliceAsHexSliceString(PDU.Data))
	return out
}

func (ADU *ModbusADU) ToString() string {
	var out string
	out += "Modbus ADU\n"
	out += ADU.MBAP.ToString()
	out += ADU.PDU.ToString()
	return out
}

/*
	Encodes the ModbusADU struct into the binary structure used in MODBUSTCP transactions
	Encodes individually the MBAP, PDU, and PDU.DATA
	Puts encodings into data buffer and returns it
	On error: returns empty []byte and passes through binary error
*/
func (ADU *ModbusADU) ToBinary() ([]byte, error) {
	data := &bytes.Buffer{}
	err := binary.Write(data, binary.BigEndian, &ADU.MBAP)
	if err != nil {
		return []byte{}, err
	}
	err = binary.Write(data, binary.BigEndian, &ADU.PDU.FunctionCode)
	if err != nil {
		return []byte{}, err
	}
	err = binary.Write(data, binary.BigEndian, ADU.PDU.Data[:(ADU.MBAP.Length-2)])
	if err != nil {
		return []byte{}, err
	}
	return data.Bytes(), nil
}

/*
	Reads bytes from buffer of a fixed length
	Parses MBAP Header and returns it
	on error: Returns empty MBAP and errors
*/
func ParseMBAPHeader(buffer []byte) (MBAPHeader, error) {
	if len(buffer) <= 0 || len(buffer) > 7 {
		return MBAPHeader{}, errors.New("incorrect buffer size")
	}
	var MBAP MBAPHeader
	reader := bytes.NewReader(buffer)
	err := binary.Read(reader, binary.BigEndian, &MBAP)
	if err != nil {
		return MBAP, err
	}
	return MBAP, nil
}

/*
	Reads butes from buffer using a []byte and predetermined length
	Parses FunctionCode and Data and puts it into a PDU
	on error: Returns empty ModbusPDU and errors
*/
func ParseModbusPDU(buffer []byte, length int16) (ModbusPDU, error){
	var PDU ModbusPDU
	reader := bytes.NewReader(buffer[7:])
	fc, err := reader.ReadByte()
	if err != nil {
		return ModbusPDU{}, err
	}
	PDU.FunctionCode = fc
	
	tempbuffer := make([]byte, length-2)
	_, err = reader.Read(tempbuffer)
	if err != nil {
		return ModbusPDU{}, err
	}
	PDU.Data = tempbuffer

	return PDU, nil
}

/*
	Creates a ModbusADU
	First parses an MBAP Header, then a Modbus PDU
	Takes the two and combines them
*/
func ParseModbusADU(buffer []byte) (ModbusADU, error) {
	if len(buffer) < 7 {
		return ModbusADU{}, errors.New("not enough data")
	}

	MBAP, err := ParseMBAPHeader(buffer[:7])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Parsing MBAP: %v", err)
	}
	if len(buffer) == 7 {
		return ModbusADU{MBAP,ModbusPDU{}}, nil
	}
	
	PDU, err := ParseModbusPDU(buffer[7:], MBAP.Length)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Parsing PDU: %v", err)
	}

	return ModbusADU{MBAP, PDU}, nil
}

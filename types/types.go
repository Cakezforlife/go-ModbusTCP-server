package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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

func ParseMBAPHeader(buffer []byte) (MBAPHeader, error) {
	if len(buffer) <= 0 || len(buffer) > 7 {
		return MBAPHeader{}, errors.New("incorrect buffer size")
	}
	var MBAP MBAPHeader
	reader := bytes.NewReader(buffer)
	err := binary.Read(reader, binary.BigEndian, &MBAP)
	if err != nil {
		return MBAP, errors.New("error parsing from buffer")
	}
	return MBAP, nil
}

func ParseModbusADU(buffer []byte) (ModbusADU, error) {
	if len(buffer) < 7 {
		return ModbusADU{}, errors.New("not enough data")
	}

	MBAP, _ := ParseMBAPHeader(buffer[:7])
	if len(buffer) == 7 {
		return ModbusADU{MBAP,ModbusPDU{}}, nil
	}
	
	var PDU ModbusPDU
	reader := bytes.NewReader(buffer[7:])
	fc, err := reader.ReadByte()
	if err != nil {
		return ModbusADU{MBAP,ModbusPDU{}}, err
	}
	PDU.FunctionCode = fc
	
	tempbuffer := make([]byte, MBAP.Length-2)
	_, err = reader.Read(tempbuffer)
	if err != nil {
		return ModbusADU{MBAP,ModbusPDU{}}, err
	}
	PDU.Data = tempbuffer

	ADU := ModbusADU{MBAP, PDU}

	return ADU, nil
}

package command

import (
	"encoding/binary"
	"log"
)

type LockStatus struct {
	GPSEnabled  bool
	IsCharching bool
	isLocked    bool
}

func ProcessHeartBeat(content []byte) {
	log.Printf("Lock Status: %#v", convertTerminalInformation(content[1]))
	log.Printf("Voltage: %d", convertVoltage(content[2:4]))
	log.Printf("GSM level: %d", content[5])
}

func GetAckHeartBeat() BL10Packet {
	result := BL10Packet{}
	result.protocolNumber = 0x23
	return result
}

func convertTerminalInformation(terminalInformationByte byte) LockStatus {
	result := LockStatus{}
	result.GPSEnabled = (terminalInformationByte>>5)&0x01 == 0x01
	result.IsCharching = terminalInformationByte>>2&0x01 == 0x01
	result.isLocked = terminalInformationByte&0x01 == 0x01
	return result
}

func convertVoltage(voltageLevelBytes []byte) uint16 {
	data := binary.BigEndian.Uint16(voltageLevelBytes)
	return data
}

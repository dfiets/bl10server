package command

import (
	bl10 "bl10server/bl10comms"
	"encoding/binary"
	"time"
)

type LockStatus struct {
	GPSEnabled  bool
	IsCharching bool
	isLocked    bool
}

func ProcessHeartBeat(content []byte, imei string) (lockStatus bl10.LockStatus) {
	lockStatus.HeartbeatPacket = extractHeartBeatData(content)
	lockStatus.Imei = imei
	lockStatus.Timestamp = time.Now().Unix()
	return lockStatus
}

func extractHeartBeatData(content []byte) *bl10.HeartBeatPacket {
	heartBeatPacket := convertTerminalInformation(content[1])
	heartBeatPacket.Voltage = int32(convertVoltage(content[2:4]))
	heartBeatPacket.SignalStrength = bl10.HeartBeatPacket_SignalStrength(content[4])
	return heartBeatPacket
}

func convertTerminalInformation(terminalInformationByte byte) *bl10.HeartBeatPacket {
	result := bl10.HeartBeatPacket{}
	result.GpsEnabled = (terminalInformationByte>>5)&0x01 == 0x01
	result.IsCharching = terminalInformationByte>>2&0x01 == 0x01
	result.IsLocked = terminalInformationByte&0x01 == 0x01
	return &result
}

func convertVoltage(voltageLevelBytes []byte) uint16 {
	data := binary.BigEndian.Uint16(voltageLevelBytes)
	return data
}

func GetAckHeartBeat() BL10Packet {
	result := BL10Packet{}
	result.protocolNumber = 0x23
	return result
}

package command

import (
	"bl10server/util"
	"encoding/binary"
)

// BL10Packet describes the parameters of a BL10Packet on a higher level.
type BL10Packet struct {
	protocolNumber byte
	content        []byte
	serialNumber   int
}

func (packet BL10Packet) CreatePacket() []byte {
	length := 1 + len(packet.content) + 2 + 2

	bPacketLength := make([]byte, 2)
	binary.BigEndian.PutUint16(bPacketLength, uint16(length))

	msg := []byte{}
	// When length byte is 1 byte
	if bPacketLength[0] == 0x0 {
		msg = append(msg, []byte{0x78, 0x78}...)
		msg = append(msg, bPacketLength[1])
	} else {
		msg = append(msg, []byte{0x79, 0x79}...)
		msg = append(msg, bPacketLength...)
	}
	msg = append(msg, packet.protocolNumber)
	msg = append(msg, packet.content...)

	bSerialNumber := make([]byte, 2)
	binary.BigEndian.PutUint16(bSerialNumber, uint16(packet.serialNumber))
	msg = append(msg, bSerialNumber...)

	// All bytes except start bytes and stop bytes
	msg = append(msg, util.CRC16Bytes(msg[2:])...)
	// Stop bytes
	msg = append(msg, []byte{0x0D, 0x0A}...)

	return msg
}

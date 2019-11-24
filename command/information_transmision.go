package command

func GetAckInformationTransmision() BL10Packet {
	result := BL10Packet{}
	result.protocolNumber = 0x98
	result.content = []byte{0x00}
	return result
}

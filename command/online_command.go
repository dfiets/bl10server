package command

import "log"

func ProcessOnlineCommand(content []byte) {
	onlineCommandResponse := string(content[6:])
	log.Println(onlineCommandResponse)
}

func GetOnlineCommand(cmd string) BL10Packet {
	result := BL10Packet{}

	result.protocolNumber = 0x80
	result.content = convertOnlineCommandContent(cmd)
	return result
}

func convertOnlineCommandContent(cmd string) []byte {
	result := []byte{}

	serverFlag := []byte{0x00, 0x00, 0x00, 0x00}
	cmdContent := []byte(cmd)
	length := byte(len(serverFlag) + len(cmdContent))

	result = append(result, length)
	result = append(result, serverFlag...)
	result = append(result, cmdContent...)
	return result
}

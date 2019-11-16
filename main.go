package main

import (
	"bl10server/command"
	"bl10server/util"
	"bytes"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	log.Println("started")
	startServer()
}

func startServer() {
	ln, err := net.Listen("tcp", ":9020")
	if err != nil {
		log.Print(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	serialNumber := 0
	for {
		readMessage(conn, serialNumber)
	}

}

func readMessage(conn net.Conn, serialNumber int) int {
	p := make([]byte, 2)
	for {
		_, err := conn.Read(p)

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		packageLength := 0
		packetLengthOneByte := bytes.Equal(p, []byte{0x78, 0x78})
		packetLengthTwoBytes := bytes.Equal(p, []byte{0x79, 0x79})

		if !(packetLengthOneByte || packetLengthTwoBytes) {
			continue
		} else if packetLengthOneByte {
			packageLength, err = getLength(conn, 1)
		} else {
			packageLength, err = getLength(conn, 2)
		}
		if err != nil {
			log.Print(err)
		}

		log.Println(packageLength)
		content := make([]byte, packageLength)
		_, err = conn.Read(content)
		if err != nil {
			log.Print(err)
		}

		responsePacket := processContent(content)
		if responsePacket.NotEmpty() {
			serialNumber += 1
			_, err := conn.Write(responsePacket.CreatePacket(serialNumber))
			if err != nil {
				log.Println(err)
			}
		}

		closeBytes := make([]byte, 2)
		_, err = conn.Read(closeBytes)
		if bytes.Equal(closeBytes, []byte{0x0D, 0x0A}) {
			log.Println("closeBytes")
		} else {
			log.Println("Something went wrong", closeBytes)
		}
	}
	return serialNumber

}

func processContent(content []byte) command.BL10Packet {
	switch content[0] {
	case 0x01:
		log.Println("LOGIN")
		command.ProcessLogin(content)
		return command.GetAckLogin(time.Now().UTC())
	case 0x23:
		log.Println("HEARTBEAT")
		command.ProcessHeartBeat(content)
		return command.GetAckHeartBeat()
	case 0x32:
		log.Println("GPS LOCATION")
	case 0x33:
		log.Println("LOCATION INFORMATION")
	case 0x80:
		log.Println("ONLINE COMMAND")
	case 0x98:
		log.Println("INFORMATION TRANSMISSION PACKET")
	default:
		log.Println("UNKNOWN protocolnumber: ERROR!!!")
	}
	return command.BL10Packet{}
}

func getLength(conn net.Conn, numberOfBytes int) (int, error) {
	bytesPacketLength := make([]byte, numberOfBytes)
	_, err := conn.Read(bytesPacketLength)
	if err != nil {
		return 0, err
	}
	return util.BytesToInt(bytesPacketLength), nil
}

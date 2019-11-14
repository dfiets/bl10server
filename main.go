package main

import (
	"bl10/command"
	"bl10/util"
	"bytes"
	"log"
	"net"
	"os"
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
	for {
		readMessage(conn)
	}

}

func readMessage(conn net.Conn) {
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
			packageLength = getLength(conn, 1)
		} else {
			packageLength = getLength(conn, 2)
		}

		content := make([]byte, packageLength)
		_, err = conn.Read(content)

		if err != nil {
			log.Print(err)
		}

		closeBytes := make([]byte, 2)
		_, err = conn.Read(closeBytes)
		if bytes.Equal(closeBytes, []byte{0x0D, 0x0A}) {
			log.Println("closeBytes")
		} else {
			log.Println("Something went wrong", closeBytes)
		}
	}

}

func processContent(content []byte) {
	switch content[0] {
	case 0x01:
		log.Println("LOGIN")
		command.ProcessLogin(con)
	case 0x23:
		log.Println("ONLINE COMMAND RESPONSE")
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
}

func getLength(conn net.Conn, numberOfBytes int) int {
	bytesPacketLength := make([]byte, numberOfBytes)
	return util.BytesToInt(bytesPacketLength)
}

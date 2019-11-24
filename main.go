package main

import (
	"bl10server/command"
	"bl10server/util"
	"bufio"
	"bytes"
	"fmt"
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
	ch := make(chan command.BL10Packet)
	go func() {

		for {
			log.Println("New message.")
			err := readMessage(conn, ch)
			if err != nil {
				log.Println("ERROR IN READ GOROUTINE")
				log.Println(err)
				return
			}
		}
	}()

	go func() {
		serialNumber := 0
		for {
			responsePacket := <-ch
			_, err := conn.Write(responsePacket.CreatePacket(serialNumber))
			if err != nil {
				log.Println("ERROR IN WRITE GOROUTINE")
				log.Println(err)
				return
			}
		}
	}()

	for {
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		switch input.Text() {
		case "u":
			fmt.Println("Unlock")
			ch <- command.GetOnlineCommand("UNLOCK#")

		case "s":
			fmt.Println("Status")
			ch <- command.GetOnlineCommand("STATUS#")
		}
	}

}

func readMessage(conn net.Conn, ch chan command.BL10Packet) error {
	p := make([]byte, 2)
	_, err := conn.Read(p)

	if err != nil {
		return err
	}

	packageLength := 0
	packetLengthOneByte := bytes.Equal(p, []byte{0x78, 0x78})
	packetLengthTwoBytes := bytes.Equal(p, []byte{0x79, 0x79})

	if !(packetLengthOneByte || packetLengthTwoBytes) {
		return nil
	} else if packetLengthOneByte {
		packageLength, err = getLength(conn, 1)
	} else {
		packageLength, err = getLength(conn, 2)
	}
	if err != nil {
		log.Print(err)
	}
	packageLength = packageLength - 4

	log.Println(packageLength)
	content := make([]byte, packageLength)
	_, err = conn.Read(content)
	if err != nil {
		log.Print(err)
	}

	responsePacket := processContent(content)
	if responsePacket.NotEmpty() {
		ch <- responsePacket
	}

	serialNumberBytes := make([]byte, 2)
	_, err = conn.Read(serialNumberBytes)
	errorCheckBytes := make([]byte, 2)
	_, err = conn.Read(errorCheckBytes)

	closeBytes := make([]byte, 2)
	_, err = conn.Read(closeBytes)
	if bytes.Equal(closeBytes, []byte{0x0D, 0x0A}) {
		log.Println("closeBytes")
	} else {
		log.Println("Something went wrong", closeBytes)
	}

	return nil

}

func processContent(content []byte) command.BL10Packet {
	switch content[0] {
	case 0x01:
		log.Println("LOGIN")
		command.ProcessLogin(content)
		return command.GetAckLogin(time.Now().UTC())
	case 0x21:
		log.Println("ONLINE COMMAND RESPONSE")
		command.ProcessOnlineCommand(content)
		return command.BL10Packet{}
	case 0x23:
		log.Println("HEARTBEAT")
		command.ProcessHeartBeat(content)
		return command.GetAckHeartBeat()
	case 0x32:
		log.Println("GPS LOCATION")
		command.ProcessGPS(content)
	case 0x33:
		log.Println("LOCATION INFORMATION")
		command.ProcessLocationAlarm(content)
	case 0x98:
		log.Println("INFORMATION TRANSMISSION PACKET")
		return command.GetAckInformationTransmision()
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

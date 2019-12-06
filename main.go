package main

import (
	"bl10server/command"
	"bl10server/util"
	"bytes"
	"errors"
	"log"
	"net"
	"time"
)

func main() {
	log.Println("started")

	go startServer()
	startGrpcServer()
}

var connections = map[int]bl10Connection{}
var imeiToConnection = map[string]int{}

type bl10Connection struct {
	conn      net.Conn
	commandCh chan command.BL10Packet
	connectCh chan confirmConnection
	connID    int
}

type confirmConnection struct {
	connID int
	imei   string
}

func Unlock(imei string) error {
	val, ok := imeiToConnection[imei]
	if !ok {
		return errors.New("This lock is not registered.")
	}

	bl10Connection, ok := connections[val]
	if !ok {
		return errors.New("Connection doesn't exist anymore.")
	}
	bl10Connection.commandCh <- command.GetOnlineCommand("SDFIND,ON,5,15,1#")
	bl10Connection.commandCh <- command.GetOnlineCommand("WIFION,30#")
	bl10Connection.commandCh <- command.GetOnlineCommand("LJDW#")
	return nil

}

func startServer() {
	ln, err := net.Listen("tcp", ":9020")
	connectionID := 1
	if err != nil {
		log.Print(err)
	}
	confirmCh := make(chan confirmConnection)

	go func() {
		for {
			confirmedConnection := <-confirmCh
			val, ok := imeiToConnection[confirmedConnection.imei]
			if ok {
				delete(connections, val)
			}
			imeiToConnection[confirmedConnection.imei] = confirmedConnection.connID
			log.Println("registered imei: ", confirmedConnection.imei)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
		ch := make(chan command.BL10Packet)
		bl10conn := bl10Connection{
			conn:      conn,
			commandCh: ch,
			connectCh: confirmCh,
			connID:    connectionID,
		}
		connections[connectionID] = bl10conn
		go bl10conn.handleConnection()
		connectionID++
	}
}

func (bl10conn bl10Connection) handleConnection() {
	stop := make(chan bool)
	go func() {

		for {
			log.Println("New message.")
			err := bl10conn.readMessage()
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
			select {
			case responsePacket := <-bl10conn.commandCh:
				_, err := bl10conn.conn.Write(responsePacket.CreatePacket(serialNumber))
				if err != nil {
					log.Println("ERROR IN WRITE GOROUTINE")
					log.Println(err)
					return
				}
			case stopNow := <-stop:
				if stopNow {
					return
				}
			}
		}
	}()

	// for {
	// 	input := bufio.NewScanner(os.Stdin)
	// 	input.Scan()
	// 	switch input.Text() {
	// 	case "u":
	// 		fmt.Println("Unlock")
	// 		ch <- command.GetOnlineCommand("UNLOCK#")

	// 	case "s":
	// 		fmt.Println("Status")
	// 		ch <- command.GetOnlineCommand("STATUS#")
	// 	}
	// }

}

func (bl10conn bl10Connection) readMessage() error {
	conn := bl10conn.conn
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

	responsePacket := bl10conn.processContent(content)
	if responsePacket.NotEmpty() {
		bl10conn.commandCh <- responsePacket
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

func (bl10conn bl10Connection) processContent(content []byte) command.BL10Packet {
	switch content[0] {
	case 0x01:
		log.Println("LOGIN")
		imei := command.ProcessLogin(content)
		bl10conn.connectCh <- confirmConnection{connID: bl10conn.connID, imei: imei}
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
	case 0x2C:
		log.Println("WIFI")
		log.Println(content)
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

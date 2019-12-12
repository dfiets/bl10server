package main

import (
	bl10 "bl10server/bl10comms"
	"bl10server/command"
	"bl10server/util"
	"bytes"
	"errors"
	"log"
	"net"
	"time"
)

func main() {
	log.Println("Started bl10server")

	go startServer()
	startGrpcServer()
}

var lockConnections = map[int]bl10Connection{}
var imeiToConnection = map[string]int{}
var serverConnections []bl10.BL10Lock_StatusUpdatesServer

type bl10Connection struct {
	conn                  net.Conn
	commandCh             chan command.BL10Packet
	connectCh             chan confirmConnection
	lockStatusBroadcastCh chan bl10.LockStatus
	connID                int
	serialNumber          int
	imei                  string
}

type confirmConnection struct {
	connID int
	imei   string
}

// SendCommandToLock this function will try to lookup the lock
// and send a command to that lock.
func SendCommandToLock(imei string, commandStr string) error {
	val, ok := imeiToConnection[imei]
	if !ok {
		return errors.New("This lock is not registered")
	}

	bl10Connection, ok := lockConnections[val]
	if !ok {
		return errors.New("Connection doesn't exist anymore")
	}
	log.Printf("Send command %s to %s", commandStr, imei)
	bl10Connection.commandCh <- command.GetOnlineCommand(commandStr)
	return nil
}

func addConsumer(stream bl10.BL10Lock_StatusUpdatesServer) {
	serverConnections = append(serverConnections, stream)
}

func startServer() {
	ln, err := net.Listen("tcp", ":9020")
	connectionID := 1
	if err != nil {
		log.Print(err)
	}
	confirmCh := make(chan confirmConnection, 100)

	go func() {
		for {
			confirmedConnection := <-confirmCh
			val, ok := imeiToConnection[confirmedConnection.imei]
			if ok {
				delete(lockConnections, val)
			}
			imeiToConnection[confirmedConnection.imei] = confirmedConnection.connID
			log.Println("Registered imei: ", confirmedConnection.imei)
		}
	}()

	lockStatusCh := make(chan bl10.LockStatus, 100)
	// This function should be made save.
	go func() {
		for {
			lockStatus := <-lockStatusCh
			index := 0
			for _, serverConn := range serverConnections {
				err := serverConn.Send(&lockStatus)
				// Clean up not working connections.
				if err == nil {
					serverConnections[index] = serverConn
					index++
				}
			}
			serverConnections = serverConnections[:index]
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
		ch := make(chan command.BL10Packet, 10)
		bl10conn := bl10Connection{
			conn:                  conn,
			commandCh:             ch,
			connectCh:             confirmCh,
			connID:                connectionID,
			lockStatusBroadcastCh: lockStatusCh,
			serialNumber:          -1,
		}
		lockConnections[connectionID] = bl10conn
		go bl10conn.handleConnection()
		connectionID++
	}
}

func (bl10conn *bl10Connection) handleConnection() {
	stop := make(chan bool)
	go func() {

		for {
			err := bl10conn.readMessage()
			if err != nil {
				log.Println("ERROR IN READ GOROUTINE")
				log.Println(err)
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case responsePacket := <-bl10conn.commandCh:
				(*bl10conn).serialNumber++
				_, err := bl10conn.conn.Write(responsePacket.CreatePacket(bl10conn.serialNumber))
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
}

func (bl10conn *bl10Connection) readMessage() error {
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
	bl10conn.serialNumber = util.BytesToInt(serialNumberBytes)
	errorCheckBytes := make([]byte, 2)
	_, err = conn.Read(errorCheckBytes)

	closeBytes := make([]byte, 2)
	_, err = conn.Read(closeBytes)
	if !bytes.Equal(closeBytes, []byte{0x0D, 0x0A}) {
		log.Printf("%s Something went wrong %+v", bl10conn.imei, closeBytes)
	}

	return nil

}

func (bl10conn *bl10Connection) processContent(content []byte) command.BL10Packet {
	switch content[0] {
	case 0x01:
		imei := command.ProcessLogin(content)
		bl10conn.imei = imei
		confirm := confirmConnection{connID: bl10conn.connID, imei: imei}
		log.Printf("%s LOGIN %+v", imei, confirm)
		bl10conn.connectCh <- confirm
		return command.GetAckLogin(time.Now().UTC())
	case 0x21:
		response := command.ProcessOnlineCommand(content)
		log.Printf("%s ONLINE COMMAND RESPONSE '%s'", bl10conn.imei, response)
		return command.BL10Packet{}
	case 0x23:
		heartbeatContent := command.ProcessHeartBeat(content, bl10conn.imei)
		log.Printf("%s HEARTBEAT %+v", bl10conn.imei, heartbeatContent)
		bl10conn.lockStatusBroadcastCh <- heartbeatContent
		return command.GetAckHeartBeat()
	case 0x32:
		gpsData := command.ProcessGPS(content, bl10conn.imei)
		bl10conn.lockStatusBroadcastCh <- gpsData
		log.Printf("%s GPS LOCATION %+v", bl10conn.imei, gpsData)
	case 0x33:
		gpsData := command.ProcessGPS(content, bl10conn.imei)
		bl10conn.lockStatusBroadcastCh <- gpsData
		log.Printf("%s LOCATION INFORMATION %+v", bl10conn.imei, gpsData)
	case 0x98:
		log.Printf("%s INFORMATION TRANSMISSION PACKET, not implemented", bl10conn.imei)
		return command.GetAckInformationTransmision()
	default:
		log.Printf("%s UNKNOWN protocolnumber: ERROR!!!", bl10conn.imei)
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

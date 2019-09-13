package main

import (
	"log"
	"net"
	"os"
	"bytes"
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
	p := make([]byte, 4)
	for {
		_, err := conn.Read(p)

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		startBytes := bytes.Compare(p, []byte{0x78})
		if startBytes == 0 {
			log.Print("Start byte 0x78")
			p2 := make([]byte, 8)
			_, err := conn.Read(p2)
			if err != nil {
				log.Print(err)
			}
			startByte2 := bytes.Compare(p[0:4], []byte{0x78})
			if startByte2 == 0 {
				log.Print("also correct, now lenght")
			}
			log.Print(p[4:8])
		}
	}

}

package command

import (
	"bl10server/util"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

func ProcessGPS(content []byte) {
	processLocation(content[1:], false)

}

func ProcessLocationAlarm(content []byte) {
	processLocation(content[1:], true)

}

func processLocation(content []byte, isAlarm bool) int {
	year := int(content[0]) + 2000
	month := int(content[1])
	day := int(content[2])
	hour := int(content[3])
	minute := int(content[4])
	second := int(content[5])

	timestamp := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
	log.Printf("timestamp %s", timestamp)

	log.Printf("GPS_INFORMATION %d", content[6])
	log.Printf("Number of satelites %d", content[7])
	startIndex := 8
	if content[6] != 0x0 {
		endIndex := startIndex + int(content[6])
		processGpsInformation(content[startIndex:endIndex])
		startIndex = endIndex
	}

	mainBaseStationStatusLength := int(content[startIndex])
	startIndex++
	if mainBaseStationStatusLength > 0 {
		endIndex := startIndex + int(mainBaseStationStatusLength)
		processBaseStationInformation(content[startIndex:endIndex])
		startIndex = endIndex
	}

	subBaseStationLength := int(content[startIndex])
	startIndex++
	if subBaseStationLength > 0 {
		startIndex = startIndex + subBaseStationLength
	}

	wifiMessageLength := int(content[startIndex])
	startIndex++
	if wifiMessageLength > 0 {
		endIndex := startIndex + 7*wifiMessageLength
		processWifiMessage(content[startIndex:endIndex], wifiMessageLength)
		startIndex = endIndex
	}

	if isAlarm {
		processStatusAlarm(content[startIndex])
	} else {
		processStatusLocation(content[startIndex])
	}

	return 0

}

func processGpsInformation(data []byte) {
	if len(data) != 12 {
		log.Printf("processGpsInformation data length not long enough, length is %d.", len(data))
		return
	}
	log.Println(util.BytesToInt(data[0:4]))
	latitude := float64(util.BytesToInt(data[0:4])) / 1800000
	log.Println(util.BytesToInt(data[4:8]))
	longitude := float64(util.BytesToInt(data[4:8])) / 18000000
	log.Printf("Location %.7f,%.7f", latitude, longitude)
}

func processBaseStationInformation(data []byte) {
	if len(data) != 9 {
		log.Printf("processBaseStationInformation data length not long enough, length is %d.", len(data))
		return
	}

	log.Println(util.BytesToInt(data[0:2]))
	log.Printf("MNC %d", data[2])
	log.Printf("LAC %d", util.BytesToInt(data[3:5]))
	log.Printf("CelltowerID", util.BytesToInt(data[5:8]))
	log.Printf("RSSI: %d", data[8])
}

func processSubBaseStationInformation() {

}

func processWifiMessage(data []byte, numberOfStations int) {
	for i := 0; i < numberOfStations; i++ {
		start := i * 7
		end := start + 6
		fmt.Println("Data: ", hex.EncodeToString(data[start:end]))
		fmt.Printf("Strength %d", data[end])
	}
}

func processStatusLocation(data byte) {
	log.Println("location")
	switch data {
	case 0x00:
		log.Println("timing report")
	case 0x01:
		log.Println("Report in fixed distance")
	case 0x02:
		log.Println("Reuplaod gps data.")
	case 0x0B:
		log.Println("LJDW report.")
	}
}

func processStatusAlarm(data byte) {
	log.Println("alarm")
	switch data {
	case 0xA0:
		log.Println("Lock report")
	case 0xA1:
		log.Println("Unlock report")
	case 0xA2:
		log.Println("Low internal battery alarm.")
	case 0xA3:
		log.Println("Low battery and shutdown.")
	case 0xA4:
		log.Println("Abnormal alarm.")
	case 0xA5:
		log.Println("Abnormal unlocking alarm.")
	}

}

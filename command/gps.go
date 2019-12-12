package command

import (
	bl10 "bl10server/bl10comms"
	"bl10server/util"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

func ProcessGPS(content []byte, imei string) (lockStatus bl10.LockStatus) {
	lockStatus.LocationPacket = processLocation(content[1:])
	lockStatus.Imei = imei
	lockStatus.Timestamp = time.Now().Unix()
	return lockStatus
}

func processLocation(content []byte) *bl10.LocationPacket {
	locationPacket := bl10.LocationPacket{}
	year := int(content[0]) + 2000
	month := int(content[1])
	day := int(content[2])
	hour := int(content[3])
	minute := int(content[4])
	second := int(content[5])
	timestamp := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
	locationPacket.LockTimestamp = timestamp.Unix()

	startIndex := 7
	if content[6] != 0x0 {
		endIndex := startIndex + int(content[6])
		locationPacket.Location = processGpsInformation(content[startIndex:endIndex])
		startIndex = endIndex
	}

	mainBaseStationStatusLength := int(content[startIndex])
	startIndex++
	if mainBaseStationStatusLength > 0 {
		endIndex := startIndex + int(mainBaseStationStatusLength)
		locationPacket.BaseStation = processBaseStationInformation(content[startIndex:endIndex])
		startIndex = endIndex
	}

	// Experimental
	subBaseStationLength := int(content[startIndex])
	startIndex++
	if subBaseStationLength > 0 {
		startIndex = startIndex + subBaseStationLength
	}

	// Experimental
	wifiMessageLength := int(content[startIndex])
	startIndex++
	if wifiMessageLength > 0 {
		endIndex := startIndex + 7*wifiMessageLength
		processWifiMessage(content[startIndex:endIndex], wifiMessageLength)
		startIndex = endIndex
	}

	locationPacket.Status = processStatus(content[startIndex])
	return &locationPacket

}

func processGpsInformation(data []byte) *bl10.PositionPackage {
	positionPackage := bl10.PositionPackage{}
	if len(data) != 12 {
		log.Printf("processGpsInformation data length not long enough, length is %d.", len(data))
		return &positionPackage
	}
	positionPackage.NumberOfSatelites = int32(data[0])
	positionPackage.Latitude = float32(util.BytesToInt(data[1:5])) / 1800000
	positionPackage.Longitude = float32(util.BytesToInt(data[5:9])) / 1800000
	positionPackage.Speed = float32(data[10])
	positionPackage.Course = int32(0x3F & binary.BigEndian.Uint16(data[11:13]))
	return &positionPackage
}

func processBaseStationInformation(data []byte) *bl10.BaseStation {
	baseStation := bl10.BaseStation{}
	if len(data) != 9 {
		log.Printf("processBaseStationInformation data length not long enough, length is %d.", len(data))
		return &baseStation
	}

	baseStation.Mcc = int32(util.BytesToInt(data[0:2]))
	baseStation.Mnc = int32(data[2])
	baseStation.Lac = int32(util.BytesToInt(data[3:5]))
	baseStation.Ci = int32(util.BytesToInt(data[5:8]))
	baseStation.Rssi = int32(data[8])
	return &baseStation
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

func processStatus(data byte) bl10.LocationPacket_Status {
	switch data {
	case 0x00:
		return bl10.LocationPacket_TIMING_REPORT
	case 0x01:
		return bl10.LocationPacket_FIXED_DISTANCE_REPORT
	case 0x02:
		return bl10.LocationPacket_GPS_REUPLOAD
	case 0x0B:
		return bl10.LocationPacket_LJDW_REPORT
	case 0xA0:
		return bl10.LocationPacket_LOCK_REPORT
	case 0xA1:
		return bl10.LocationPacket_UNLOCK_REPORT
	case 0xA2:
		return bl10.LocationPacket_LOW_INTERNAL_BATTERY_ALARM
	case 0xA3:
		return bl10.LocationPacket_LOW_BATTERY_SHUTDOWN
	case 0xA4:
		return bl10.LocationPacket_ABNORMAL_ALARM
	case 0xA5:
		return bl10.LocationPacket_ABNORMAL_UNLOCKING_ALARM
	}
	return 0xFF
}

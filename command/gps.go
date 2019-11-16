package command

import (
	"bl10server/util"
	"log"
	"time"
)

func ProcessGPS(content []byte) {
	processLocation(content[1:])

}

func ProcessLocationAlarm(content []byte) {
	processLocation(content[1:])

}

func processLocation(content []byte) int {
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
	latitude := float64(util.BytesToInt(content[8:12])) / 1800000
	longitude := float64(util.BytesToInt(content[12:16])) / 18000000
	log.Printf("Location %.7f,%.7f", latitude, longitude)

	return 16

}

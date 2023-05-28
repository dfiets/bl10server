package command

import (
	"strconv"
	"time"
)

func ProcessLogin(content []byte) string {
	return convertBytesToIMEI(content[1:9])
}

func GetAckLogin(now time.Time) BL10Packet {
	content := []byte{byte(now.Year() - 2000), byte(now.Month()), byte(now.Day()),
		byte(now.Hour()), byte(now.Minute()), byte(now.Second()), 0x00}

	packet := BL10Packet{
		protocolNumber: 0x01,
		content:        content,
	}
	return packet

}

func convertBytesToIMEI(imeiBytes []byte) string {
	imei := ""
	for index, imeiByte := range imeiBytes {
		if index != 0 {
			firstDigit := imeiByte >> 4
			imei += strconv.Itoa(int(firstDigit))
		}
		secondDigit := 0x0F & imeiByte
		imei += strconv.Itoa(int(secondDigit))
	}
	return imei
}

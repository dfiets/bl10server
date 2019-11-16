package command

import (
	"bl10server/util"
	"log"
	"strconv"
)

func ProcessLogin(content []byte) {
	log.Println("IMEI ", convertBytesToIMEI(content[1:9]))
	log.Println("Sequence number: ", util.BytesToInt(content[13:15]))
}

func GetAckLogin() {

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

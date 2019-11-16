package command

import (
	"bytes"
	"testing"
	"time"
)

func TestConvertBytesToIMEI(t *testing.T) {
	expectedOutput := "123456789123456"
	input := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x12, 0x34, 0x56}
	imei := convertBytesToIMEI(input)
	if imei != expectedOutput {
		t.Errorf("convertBytesToIMEI was incorrect, got: %s, want: %s.", imei, expectedOutput)
	}
}

func TestResponseLogin(t *testing.T) {
	expectedOutput := []byte{0x78, 0x78, 0x0C, 0x01, 0x11, 0x03, 0x14, 0x08, 0x38, 0x39, 0x00, 0x00, 0x39, 0x95, 0x70, 0x0D, 0x0A}
	input := time.Date(
		2017, 3, 20, 8, 56, 57, 0, time.UTC)
	serialNumber := 57
	result := GetAckLogin(input).CreatePacket(serialNumber)
	if !bytes.Equal(result, expectedOutput) {
		t.Errorf("GetAckLogin was incorrect, got: % x, want: % x.", result, expectedOutput)
	}

}

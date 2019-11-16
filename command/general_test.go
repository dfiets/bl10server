package command

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBL10CreatePacket(t *testing.T) {
	expectedOutput := []byte{0x78, 0x78, 0x0C, 0x01, 0x11, 0x03, 0x14, 0x08, 0x38, 0x39, 0x00, 0x00, 0x39, 0x95, 0x70, 0x0D, 0x0A}

	// Setup input
	contentBytes := []byte{0x11, 0x03, 0x14, 0x08, 0x38, 0x39, 0x00}
	protocolNumber := byte(0x01)
	serialNumber := 57
	input := BL10Packet{
		protocolNumber: protocolNumber,
		serialNumber:   serialNumber,
		content:        contentBytes,
	}

	result := input.CreatePacket()

	fmt.Printf("%s", result)
	fmt.Printf("%s", expectedOutput)
	if !bytes.Equal(result, expectedOutput) {
		t.Errorf("CreatePacket was incorrect, got: % x, want: % x.", result, expectedOutput)
	}
}

package command

import (
	"bytes"
	"testing"
)

func TestConvertOnlineCommandContent(t *testing.T) {
	expectedOutput := []byte{0x0B, 0x00, 0x00, 0x00, 0x00, 0x55, 0x4E, 0x4C, 0x4F, 0x43, 0x4B, 0x23}
	input := "UNLOCK#"
	result := convertOnlineCommandContent(input)
	if !bytes.Equal(expectedOutput, result) {
		t.Errorf("%s: Expected % x, got % x", input, expectedOutput, result)
	}
}

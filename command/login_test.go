package command

import "testing"

func TestConvertBytesToIMEI(t *testing.T) {
	expectedOutput := "123456789123456"
	input := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x12, 0x34, 0x56}
	imei := convertBytesToIMEI(input)
	if imei != expectedOutput {
		t.Errorf("convertBytesToIMEI was incorrect, got: %s, want: %s.", imei, expectedOutput)
	}
}

package command

import (
	"testing"
)

type TestCase struct {
	input    byte
	expected LockStatus
}

func TestConvertTerminalInformation(t *testing.T) {
	cases := []TestCase{
		TestCase{
			input: 0x00,
			expected: LockStatus{
				GPSEnabled:  false,
				IsCharching: false,
				isLocked:    false,
			},
		},
		TestCase{
			input: 0x01,
			expected: LockStatus{
				GPSEnabled:  false,
				IsCharching: false,
				isLocked:    true,
			},
		},
		TestCase{
			input: 0x01 << 2,
			expected: LockStatus{
				GPSEnabled:  false,
				IsCharching: true,
				isLocked:    false,
			},
		},
		TestCase{
			input: 0x01 << 5,
			expected: LockStatus{
				GPSEnabled:  true,
				IsCharching: false,
				isLocked:    false,
			},
		},
		TestCase{
			input: 0x01<<5 | 0x01<<2 | 0x01,
			expected: LockStatus{
				GPSEnabled:  true,
				IsCharching: true,
				isLocked:    true,
			},
		},
		TestCase{
			input: 0x01<<5 | 0x01<<2,
			expected: LockStatus{
				GPSEnabled:  true,
				IsCharching: true,
				isLocked:    false,
			},
		},
	}

	for _, tc := range cases {
		got := convertTerminalInformation(tc.input)
		if got != tc.expected {
			t.Errorf("% x: Expected %#v, got %#v", tc.input, tc.expected, got)
		}
	}
}

func TestConvertVoltage(t *testing.T) {
	expectedOutput := 415
	input := []byte{0x01, 0x9F}
	voltage := convertVoltage(input)
	if expectedOutput != int(voltage) {
		t.Errorf("convertBytesVoltage was incorrect, got: %d, want: %d.", voltage, expectedOutput)
	}
}

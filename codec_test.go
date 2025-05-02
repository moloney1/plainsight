package main

import (
	"bytes"
	"testing"
)

var bytesIn []byte = []byte{
	50, 51, 52, 52, 55, 56, 56, 56, // H
	56, 57, 57, 58, 58, 59, 60, 61, // e
	62, 63, 63, 64, 65, 65, 66, 66, // l
	62, 63, 63, 64, 65, 65, 66, 66, // l
	66, 67, 67, 68, 69, 69, 69, 69, // o
	70, 70, 71, 72, 72, 72, 72, 72, // space
	70, 71, 72, 73, 72, 73, 73, 73, // W
	66, 67, 67, 68, 69, 69, 69, 69, // o
	70, 71, 71, 71, 72, 72, 73, 74, // r
	62, 63, 63, 64, 65, 65, 66, 66, // l
	66, 67, 67, 68, 68, 69, 70, 70, // d
}

func TestReadMessagePositive(t *testing.T) {
	var tests = []struct {
		name       string
		toBytesArg int
		expected   string
	}{
		{"ReadMessageP1", 0, ""},
		{"ReadMessageP2", 1, "H"},
		{"ReadMessageP3", 5, "Hello"},
		{"ReadMessageP4", 11, "Hello World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ReadMessage(bytesIn, tt.toBytesArg)
			if err != nil {
				t.Errorf("unexpected error %v\n", err)
			}
			if res != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, res)
			}
		})
	}
}

func TestReadMessageNegative(t *testing.T) {
	var tests = []struct {
		name       string
		toBytesArg int
	}{
		{"ReadMessageN1", 42},
		{"ReadMessageN2", -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadMessage(bytesIn, tt.toBytesArg)
			if err == nil {
				t.Errorf("expected err")
			}
		})
	}
}

func TestWriteMessagePositive(t *testing.T) {

	var tests = []struct {
		name         string
		message      string
		fromBytesArg int
	}{
		{"WriteMessageP1", "Change!", 5},
		{"WriteMessageP2", "Howdy Earth", 0}, // exactly enough space
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesInCopy := make([]byte, len(bytesIn))
			copy(bytesInCopy, bytesIn)

			bytesInCopy, err := WriteMessage(tt.message, bytesInCopy, tt.fromBytesArg)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if bytes.Equal(bytesInCopy, bytesIn) {
				t.Error("expected bytesInCopy to have changed")
			}
		})
	}

}

func TestWriteMessageNegative(t *testing.T) {
	var tests = []struct {
		name        string
		message     string
		fromByteArg int
	}{
		{"WriteMessageN1", "", 2319},
		{"WriteMessageN2", "", -5},
		{"WriteMessageN3", "Howdy Earths", 0}, // One character too large
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesInCopy := make([]byte, len(bytesIn))
			copy(bytesInCopy, bytesIn)

			_, err := WriteMessage(tt.message, bytesInCopy, tt.fromByteArg)

			if err == nil {
				t.Errorf("expected err")
			}
		})
	}
}

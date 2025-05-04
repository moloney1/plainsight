package codec

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const bitsPerByte = 8

// Encode message contained in 'message' to string and return it
func encodeMessage(message string) string {
	var builder strings.Builder

	for i := range len(message) {
		builder.WriteString(fmt.Sprintf("%08b", message[i]))
	}

	return builder.String()
}

// Decode message from binary and return it as plaintext string
func decodeMessage(message string) string {
	var builder strings.Builder

	bytesRead, start := 0, 0

	for bytesRead < int(len(message)/bitsPerByte) {
		i, _ := strconv.ParseInt(message[start:start+bitsPerByte], 2, 32)
		builder.WriteByte(byte(i))
		bytesRead += 1
		start += bitsPerByte
	}

	return builder.String()
}

// Read the bits up to bytesToRead, return them decoded to plaintext
func ReadMessage(bytes []byte, bytesToRead int) (string, error) {
	if bytesToRead > len(bytes)/bitsPerByte || bytesToRead < 0 {
		return "", fmt.Errorf("value %d for bytesToRead is invalid", bytesToRead)
	}

	var builder strings.Builder

	for i := range bitsPerByte * bytesToRead {
		builder.WriteString(
			fmt.Sprintf("%v", bytes[i]&1), // &1 grabs the LSB
		)
	}

	return decodeMessage(builder.String()), nil
}

// Encode message to binary and write it to the byte slice (from byte 0) via LSB steganography, return modified byte slice
// Note fromByte needs to be 0 or multiple of 8 to play nicely
func WriteMessage(message string, bytes []byte, fromByte int) ([]byte, error) { // TODO bounds
	if fromByte+(len(message)*8) > len(bytes) || fromByte < 0 {
		return []byte{}, errors.New("invalid args: out of range")
	}

	messageBin := encodeMessage(message)

	for i := range len(messageBin) {
		// messageBin[i] is either '0' (ASCII 48) or '1' (49) so subtracting '0' gives us the bit value
		messageBit := messageBin[i] - '0'

		idx := i + fromByte
		if messageBit == 0 {
			bytes[idx] = bytes[idx] &^ 1 // &^ 1 sets LSB to 0
		} else {
			bytes[idx] = bytes[idx] | 1 // | 1 sets LSB to 1
		}
	}
	return bytes, nil
}

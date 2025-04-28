package main

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"
)

const inputFile = "test_img.png"
const outputFile = "output.png"
const bitsPerByte = 8

func main() {

	img, err := readImage(inputFile)
	if err != nil {
		panic(err)
	}

	if err = writeMessageToImage("Hello World!", img, outputFile); err != nil {
		panic(err)
	}

	img, err = readImage(outputFile)

	decodedMessage, err := readMessageFromImage(*img, 12)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read following message from %s: '%s'", outputFile, decodedMessage)

}

// read image from 'filename' and decode to PNG, convert from image.RGBA to image.NRGBA if necessary
func readImage(filename string) (*image.NRGBA, error) {
	imgFile, err := os.Open(filename)
	if err != nil {
		return &image.NRGBA{}, err
	}
	defer imgFile.Close()

	reader := bufio.NewReader(imgFile)

	pngImg, err := png.Decode(reader)
	if err != nil {
		return &image.NRGBA{}, err
	}

	nrgbaImg, ok := pngImg.(*image.NRGBA)
	if !ok {
		fmt.Printf("Image %s is not of type *image.NRGBA. Attempting conversion\n", filename)
		// Thanks, https://go.dev/blog/image-draw#converting-an-image-to-rgba
		bounds := pngImg.Bounds()
		nrgbaImg = image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
		draw.Draw(nrgbaImg, nrgbaImg.Bounds(), pngImg, bounds.Min, draw.Src)
	}

	return nrgbaImg, nil
}

// Encode new PNG file and write to disk at location 'filename'
func writeImage(filename string, img image.Image) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)
	err = png.Encode(writer, img)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil

}

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
func readMessageFromImage(img image.NRGBA, bytesToRead int) (string, error) {
	var builder strings.Builder

	for i := range bitsPerByte * bytesToRead {
		builder.WriteString(
			fmt.Sprintf("%v", img.Pix[i]&1), // &1 grabs the LSB
		)
	}

	return decodeMessage(builder.String()), nil
}

// Encode message to binary and write it to the image (from pixel 0) via LSB steganography
func writeMessageToImage(message string, img *image.NRGBA, outputFilename string) error {
	messageBin := encodeMessage(message)

	for i := range len(messageBin) {
		// messageBin[i] is either '0' (ASCII 48) or '1' (49) so subtracting '0' gives us the bit value
		messageBit := messageBin[i] - '0'

		if messageBit == 0 {
			img.Pix[i] = img.Pix[i] &^ 1 // &^ 1 sets LSB to 0
		} else {
			img.Pix[i] = img.Pix[i] | 1 // | 1 sets LSB to 1
		}
	}

	return writeImage(outputFilename, img)
}

package main

import "fmt"

const inputFile = "test_img.png"
const outputFile = "output.png"

func main() {

	img, err := ReadImage(inputFile)
	if err != nil {
		panic(err)
	}

	if err = WriteMessageToImage("Hello World!", img, outputFile); err != nil {
		panic(err)
	}

	img, err = ReadImage(outputFile)

	decodedMessage, err := ReadMessageFromImage(*img, 12)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read following message from %s: '%s'", outputFile, decodedMessage)

}

package main

import "fmt"

const inputFile = "test_img.png"
const outputFile = "output.png"

func main() {
	testCodec()

}

func testCodec() {
	img, err := ReadImage(inputFile)
	if err != nil {
		panic(err)
	}

	WriteMessage("Hello World!", img.Pix)

	if err = WriteImageFile(outputFile, img); err != nil {
		panic(err)
	}

	img, err = ReadImage(outputFile)
	if err != nil {
		panic(err)
	}

	decodedMessage, err := ReadMessage(img.Pix, 12)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read following message from %s: '%s'", outputFile, decodedMessage)

}

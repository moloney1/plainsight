package main

import "fmt"

const inputFile = "test_img.png"
const outputFile = "output.png"

func main() {
	// testCodec()
	testTable()

}

func testTable() {

	img, err := ReadImage(inputFile)
	if err != nil {
		panic(err)
	}

	t := NewTable(img.Pix)
	fmt.Println(t.Meta.Cap)
	fmt.Println(t.Meta.Size)
	fmt.Println(t.Meta.Keys)
	fmt.Println(len(t.Data))

	img.Pix = t.Data
	WriteImageFile(outputFile, img)

	newImg, _ := ReadImage(outputFile)
	decodedMessage, err := ReadMessage(newImg.Pix, 44)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read following message from %s: '%s'", outputFile, decodedMessage)

}

func testCodec() {
	img, err := ReadImage(inputFile)
	if err != nil {
		panic(err)
	}

	img.Pix = WriteMessage("Hello World!", img.Pix)

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

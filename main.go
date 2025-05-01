package main

import (
	"fmt"
	"hash/fnv"
)

const inputFile = "test_img.png"
const outputFile = "output.png"

func main() {
	// testCodec()
	//testTable()
	testHashAddAndRead()
}

func testHashAddAndRead() {
	h := fnv.New64a()
	h.Write([]byte("passwordorsomething"))
	fmt.Printf("%v\n", h.Sum64())

	g := fnv.New64a()
	g.Write([]byte("hello"))
	fmt.Printf("%v\n", g.Sum64())

	img, _ := ReadImage(inputFile)

	someValidJson := "{\"test\": \"yes\"}"
	t := NewTable(img.Pix)
	t.Add("myKey", someValidJson)
	WriteImageFile(outputFile, img)

	newImg, _ := ReadImage(outputFile)
	newT, err := TableFromBytes(newImg.Pix)
	if err != nil {
		panic(err)
	}
	fmt.Println(newT.Meta)
	msg, err := ReadMessage(t.Data[73088:], 15)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got message: %s\n", msg)

	s, err := t.Read("myKey")
	if err != nil {
		panic(err)
	}
	fmt.Printf("reading myKey: %s\n", s)

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
	// decodedMessage, err := ReadMessage(newImg.Pix, 44)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Read following message from %s: '%s'", outputFile, decodedMessage)
	newT, err := TableFromBytes(newImg.Pix)
	if err != nil {
		panic(err)
	}
	fmt.Println(newT.Meta)

}

func testCodec() {
	img, err := ReadImage(inputFile)
	if err != nil {
		panic(err)
	}

	img.Pix = WriteMessage("Hello World!", img.Pix, 8)

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

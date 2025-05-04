package main

import (
	"fmt"
	"hash/fnv"

	"github.com/moloney1/plainsight/cmd"
	"github.com/moloney1/plainsight/internal/codec"
	"github.com/moloney1/plainsight/internal/imageio"
	"github.com/moloney1/plainsight/internal/table"
)

const inputFile = "test_img.png"
const outputFile = "output.png"

func main() {
	// 	testCodec()
	// 	fmt.Println()
	//
	// 	testTable()
	// 	fmt.Println()
	//
	// 	testHashAddAndRead()
	// 	fmt.Println()
	//
	cmd.Execute()

}

func testHashAddAndRead() {
	h := fnv.New64a()
	h.Write([]byte("passwordorsomething"))
	fmt.Printf("%v\n", h.Sum64())

	g := fnv.New64a()
	g.Write([]byte("hello"))
	fmt.Printf("%v\n", g.Sum64())

	img, _ := imageio.ReadImage(inputFile)

	someValidJson := "{\"test\": \"yes\"}"
	t, _ := table.NewTable(img.Pix, fnv.New64a())
	t.Add("myKey", someValidJson)
	imageio.WriteImageFile(outputFile, img)

	newImg, _ := imageio.ReadImage(outputFile)
	newT, err := table.TableFromBytes(newImg.Pix, fnv.New64a())
	if err != nil {
		panic(err)
	}
	fmt.Println(newT.Meta)
	msg, err := codec.ReadMessage(t.Data[73088:], 15)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got message: %s\n", msg)

	s, err := newT.Read("myKey")
	if err != nil {
		panic(err)
	}
	fmt.Printf("reading myKey: %s\n", s)

	// try that again to be sure hash.Reset() worked
	s, err = newT.Read("myKey")
	if err != nil {
		panic(err)
	}
	fmt.Printf("reading myKey: %s\n", s)

}

func testTable() {

	img, err := imageio.ReadImage(inputFile)
	if err != nil {
		panic(err)
	}

	t, _ := table.NewTable(img.Pix, fnv.New64a())
	fmt.Println(t.Meta.Cap)
	fmt.Println(t.Meta.Size)
	fmt.Println(t.Meta.Keys)
	fmt.Println(len(t.Data))

	img.Pix = t.Data
	imageio.WriteImageFile(outputFile, img)

	newImg, _ := imageio.ReadImage(outputFile)
	// decodedMessage, err := ReadMessage(newImg.Pix, 44)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Read following message from %s: '%s'", outputFile, decodedMessage)
	newT, err := table.TableFromBytes(newImg.Pix, fnv.New64a())
	if err != nil {
		panic(err)
	}
	fmt.Println(newT.Meta)

}

func testCodec() {
	img, err := imageio.ReadImage(inputFile)
	if err != nil {
		panic(err)
	}

	// img.Pix, _ = WriteMessage("Hello World!", img.Pix, 8)
	img.Pix, _ = codec.WriteMessage("Hello World!", img.Pix, 0)

	if err = imageio.WriteImageFile(outputFile, img); err != nil {
		panic(err)
	}

	img, err = imageio.ReadImage(outputFile)
	if err != nil {
		panic(err)
	}

	decodedMessage, err := codec.ReadMessage(img.Pix, 12)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read following message from %s: '%s'", outputFile, decodedMessage)

}

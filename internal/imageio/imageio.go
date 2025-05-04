package imageio

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

// Read image from 'filename' and decode to PNG, convert from image.RGBA to image.NRGBA if necessary
func ReadImage(filename string) (*image.NRGBA, error) {
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
func WriteImageFile(filename string, img image.Image) error {
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

package cmd

import (
	"fmt"
	"hash/fnv"
	"image"
	"os"

	"github.com/moloney1/plainsight/internal/imageio"
	"github.com/moloney1/plainsight/internal/table"
)

// openTable reads an existing Table from a file. It should be used when a table is expected to exist already.
func openTable(imageFile string) (*image.NRGBA, *table.Table) {
	img, err := imageio.ReadImage(imageFile)
	if err != nil {
		fmt.Printf("Error reading image file: %s", err)
		os.Exit(1)
	}
	t, err := table.TableFromBytes(img.Pix, fnv.New64a())
	if err != nil {
		fmt.Printf("Plainsight does not recognize the file %s", imageFile)
		os.Exit(1)
	}
	return img, t
}

// openOrCreateTable reads an existing Table from a file or creates one if it doesn't exist. It should be used when
// a table not existing already would not constitute an error.
func openOrCreateTable(imageFile string) (*image.NRGBA, *table.Table) {
	img, err := imageio.ReadImage(imageFile)
	if err != nil {
		fmt.Printf("Error reading image file: %s", err)
		os.Exit(1)
	}
	hasher := fnv.New64a()
	t, err := table.TableFromBytes(img.Pix, hasher)
	if err != nil {
		fmt.Printf("Plainsight does not recognize the file. Attempting to register %s\n", imageFile)
		t, err = table.NewTable(img.Pix, hasher)
		if err != nil {
			fmt.Printf("Error registering image file %s: %v", imageFile, err)
			os.Exit(1)
		}
	}
	return img, t
}

// saveOutput writes the new data to a file.
func saveOutput(img image.Image) {
	imageio.WriteImageFile("new.png", img) // for now, always create a new file instead of overwriting.
	fmt.Print("Output written to file new.png")
}

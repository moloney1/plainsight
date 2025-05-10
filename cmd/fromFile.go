package cmd

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/moloney1/plainsight/internal/imageio"
	"github.com/moloney1/plainsight/internal/table"
)

func init() {
	addCmd.AddCommand(fromFileCmd)

	fromFileCmd.PersistentFlags().StringVarP(&sourceFile, "source-file", "s", "", "path to JSON file")
	fromFileCmd.MarkPersistentFlagRequired("source-file")

}

var fromFileCmd = &cobra.Command{
	Use:     "from-file",
	Aliases: []string{""},
	Short:   "Add arbitrary JSON data from file to your image",
	Long:    `Add arbitrary JSON data from file to your image under supplied key. Example: 'plainsight add --file myImageFile.png --key myKey from-file --source-file myData.json`,
	Run: func(cmd *cobra.Command, args []string) {
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

		// TODO this logic should be moved to table.go
		if slices.Contains(t.Meta.Keys, key) {
			fmt.Printf("Image %s already has data stored under key %s. Taking no action", imageFile, key)
			os.Exit(1)
		}

		file, err := os.Open(sourceFile)
		if err != nil {
			fmt.Printf("Error reading source file %s: %v", sourceFile, err)
			os.Exit(1)
		}
		defer file.Close()
		jsonBytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Error reading source file %s: %v", sourceFile, err)
			os.Exit(1)
		}

		var jsonMap map[string]any

		err = json.Unmarshal(jsonBytes, &jsonMap)
		if err != nil {
			fmt.Printf("Cannot add data from file %s. Likely invalid JSON. Error: %v", sourceFile, err)
			os.Exit(1)
		}

		if err = t.Add(key, strings.TrimSpace(string(jsonBytes))); err != nil {
			fmt.Printf("Failed to write data to file new.png. Error: %s", err)
			os.Exit(1)
		}

		imageio.WriteImageFile("new.png", img) // for now, always create a new file instead of overwriting.
		fmt.Print("Output written to file new.png")
	},
}

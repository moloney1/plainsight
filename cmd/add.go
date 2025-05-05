package cmd

import (
	"fmt"
	"hash/fnv"
	"os"

	"github.com/spf13/cobra"

	"github.com/moloney1/plainsight/internal/imageio"
	"github.com/moloney1/plainsight/internal/table"
)

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&imageFile, "file", "f", "", "image file path")
	addCmd.MarkPersistentFlagRequired("file")

	addCmd.PersistentFlags().StringVarP(&keyToSearch, "key", "k", "", "key to read")
	addCmd.MarkPersistentFlagRequired("key")

	addCmd.PersistentFlags().StringVarP(&jsonData, "data", "d", "", "data to add")
	addCmd.MarkPersistentFlagRequired("data")

}

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add data to your image",
	Long:    `Add specified JSON data to image under supplied key`,
	Run: func(cmd *cobra.Command, args []string) {
		img, err := imageio.ReadImage(imageFile)
		if err != nil {
			fmt.Printf("Error reading image file: %s\n", err)
			os.Exit(1)
		}
		hasher := fnv.New64a()
		t, err := table.TableFromBytes(img.Pix, hasher)
		if err != nil {
			fmt.Printf("Plainsight does not recognize the file %s\n", imageFile)
			t, err = table.NewTable(img.Pix, hasher)
			if err != nil {
				fmt.Printf("Error initializing image file %s: %v", imageFile, err)
				os.Exit(1)
			}
		}

		if err := t.Add(keyToSearch, jsonData); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		imageio.WriteImageFile("new.png", img) // for now, always create a new file instead of overwriting.
		fmt.Println("successful!")
	},
}

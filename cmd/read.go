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
	rootCmd.AddCommand(readCmd)

	readCmd.PersistentFlags().StringVarP(&imageFile, "file", "f", "", "image file path")
	readCmd.MarkPersistentFlagRequired("file")

	readCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "key to read")
	readCmd.MarkPersistentFlagRequired("key")

}

var readCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Read the data at the supplied key",
	Long:    `Read data from file at key if it exists. Data is returned as JSON string. Example: 'plainsight list --file myImageFile.png --key myKey'`,
	Run: func(cmd *cobra.Command, args []string) {
		img, err := imageio.ReadImage(imageFile)
		if err != nil {
			fmt.Printf("Error reading image file: %s", err)
			os.Exit(1)
		}
		table, err := table.TableFromBytes(img.Pix, fnv.New64a())
		if err != nil {
			fmt.Printf("Plainsight does not recognize the file %s", imageFile)
			os.Exit(1)
		}

		result, err := table.Read(key)
		if err != nil {
			fmt.Printf("Error reading data for key %s: %v", key, err)
			os.Exit(1)
		}
		fmt.Printf("Found the following data for key %s: %s", key, result)
	},
}

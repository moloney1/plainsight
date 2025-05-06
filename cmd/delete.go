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
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.PersistentFlags().StringVarP(&imageFile, "file", "f", "", "image file path")
	deleteCmd.MarkPersistentFlagRequired("file")

	deleteCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "key to delete")
	deleteCmd.MarkPersistentFlagRequired("key")

}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del"},
	Short:   "Delete data stored at key from your image",
	Long:    `Delete the data stored under supplied key, if it exists. Example: 'plainsight delete --file myImageFile.png --key myKey'`,
	Run: func(cmd *cobra.Command, args []string) {
		img, err := imageio.ReadImage(imageFile)
		if err != nil {
			fmt.Printf("Error reading image file: %s", err)
			os.Exit(1)
		}
		hasher := fnv.New64a()
		t, err := table.TableFromBytes(img.Pix, hasher)
		if err != nil {
			fmt.Printf("Plainsight does not recognize the file %s", imageFile)
			os.Exit(1)
		}

		if err := t.Delete(key); err != nil {
			fmt.Printf("Error deleting data under key %s: %v", key, err)
			os.Exit(1)
		}

		imageio.WriteImageFile("new.png", img) // for now, always create a new file instead of overwriting.
		fmt.Print("Output written to file new.png")
	},
}

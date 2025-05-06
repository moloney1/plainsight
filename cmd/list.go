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
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringVarP(&imageFile, "file", "f", "", "image file path")
	listCmd.MarkPersistentFlagRequired("file")
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the keys of the data stored in file",
	Long:    `List the keys of the data stored in the specified PNG file. Example: 'plainsight list --file myImageFile.png'`,
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

		list := table.List()
		if len(list) == 0 {
			fmt.Printf("No data found in %s", imageFile)
			os.Exit(0)
		}
		fmt.Printf("Found the following data in %s: ", imageFile)
		for i := range list {
			fmt.Printf("%s ", list[i])
		}
	},
}

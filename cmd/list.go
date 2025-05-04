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
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the keys of the data stored in file",
	Long:    `List the keys of the data stored in file`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("The file you want to look at is %v\n", rootCmd.PersistentFlags().Lookup("file").Value)

		imgPath := rootCmd.PersistentFlags().Lookup("file").Value.String()
		img, err := imageio.ReadImage(imgPath)
		if err != nil {
			fmt.Printf("Error reading image file: %s\n", err)
			os.Exit(1)
		}
		table, err := table.TableFromBytes(img.Pix, fnv.New64a())
		if err != nil {
			fmt.Printf("Plainsight does not recognize the file %s\n", imgPath)
		}
		fmt.Println(table.List())
	},
}

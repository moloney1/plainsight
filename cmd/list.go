package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringVarP(&imageFile, "file", "f", "", "image file path")
	if err := listCmd.MarkPersistentFlagRequired("file"); err != nil {
		panic(err)
	}
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the keys of the data stored in file",
	Long:    `List the keys of the data stored in the specified PNG file. Example: 'plainsight list --file myImageFile.png'`,
	Run: func(cmd *cobra.Command, args []string) {
		_, tbl := openTable(imageFile)

		list := tbl.List()
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

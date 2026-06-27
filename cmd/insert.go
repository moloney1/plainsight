package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(insertCmd)

	insertCmd.PersistentFlags().StringVarP(&imageFile, "file", "f", "", "image file path")
	if err := insertCmd.MarkPersistentFlagRequired("file"); err != nil {
		panic(err)
	}

	insertCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "key to add data under")
	if err := insertCmd.MarkPersistentFlagRequired("key"); err != nil {
		panic(err)
	}

}

var insertCmd = &cobra.Command{
	Use:     "insert",
	Aliases: []string{"i"},
	Short:   "Add data to your image",
	Long:    `Add specified data to image under supplied key. See 'plainsight help insert' for subcommand options.`,
}

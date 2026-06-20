package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&imageFile, "file", "f", "", "image file path")
	if err := addCmd.MarkPersistentFlagRequired("file"); err != nil {
		panic(err)
	}

	addCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "key to add data under")
	if err := addCmd.MarkPersistentFlagRequired("key"); err != nil {
		panic(err)
	}

}

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add data to your image",
	Long:    `Add specified data to image under supplied key. See 'plainsight help add' for subcommand options.`,
}

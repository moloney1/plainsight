package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.PersistentFlags().StringVarP(&imageFile, "file", "f", "", "image file path")
	if err := removeCmd.MarkPersistentFlagRequired("file"); err != nil {
		panic(err)
	}

	removeCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "key to delete")
	if err := removeCmd.MarkPersistentFlagRequired("key"); err != nil {
		panic(err)
	}
}

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove data stored at key from your image",
	Long:    `Remove the data stored under supplied key, if it exists. Example: 'plainsight remove --file myImageFile.png --key myKey'`,
	Run: func(cmd *cobra.Command, args []string) {
		img, t := openTable(imageFile)

		if err := t.Remove(key); err != nil {
			fmt.Printf("Error removing data under key %s: %v", key, err)
			os.Exit(1)
		}

		saveOutput(img, imageFile)
	},
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
		img, t := openTable(imageFile)

		if err := t.Delete(key); err != nil {
			fmt.Printf("Error deleting data under key %s: %v", key, err)
			os.Exit(1)
		}

		saveOutput(img)
	},
}

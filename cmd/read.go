package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

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
		_, tbl := openTable(imageFile)

		jsonResult, err := tbl.Read(key)
		if err != nil {
			fmt.Printf("Error reading data for key %s: %v", key, err)
			os.Exit(1)
		}

		var userPassPair table.UserPass
		if err = json.Unmarshal([]byte(jsonResult), &userPassPair); err != nil {
			fmt.Printf("Error reading data for key %s: %v", key, err)
			os.Exit(1)
		}

		if userPassPair.User == "" && userPassPair.Pass == "" {
			// TODO if stored arbitrary JSON has user or pass key then...
			fmt.Printf("Found the following data for key %s: %s", key, jsonResult)
			os.Exit(0)
		}

		fmt.Printf("Found the following data for key %s:\n", key)
		fmt.Printf("Username: %s\n", userPassPair.User)
		fmt.Printf("Password: %s\n", userPassPair.Pass)
	},
}

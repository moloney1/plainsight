package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	insertCmd.AddCommand(fromFileCmd)

	fromFileCmd.PersistentFlags().StringVarP(&sourceFile, "source-file", "s", "", "path to JSON file")
	if err := fromFileCmd.MarkPersistentFlagRequired("source-file"); err != nil {
		panic(err)
	}
}

var fromFileCmd = &cobra.Command{
	Use:     "from-file",
	Aliases: []string{""},
	Short:   "Add arbitrary JSON data from file to your image",
	Long:    `Add arbitrary JSON data from file to your image under supplied key. Example: 'plainsight insert --file myImageFile.png --key myKey from-file --source-file myData.json`,
	Run: func(cmd *cobra.Command, args []string) {
		img, t := openOrCreateTable(imageFile)

		// TODO this logic should be moved to table.go
		if slices.Contains(t.Meta.Keys, key) {
			fmt.Printf("Image %s already has data stored under key %s. Taking no action", imageFile, key)
			os.Exit(1)
		}

		file, err := os.Open(sourceFile)
		if err != nil {
			fmt.Printf("Error reading source file %s: %v", sourceFile, err)
			os.Exit(1)
		}
		defer file.Close() //nolint
		jsonBytes, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Error reading source file %s: %v", sourceFile, err)
			os.Exit(1)
		}

		var jsonMap map[string]any
		if err = json.Unmarshal(jsonBytes, &jsonMap); err != nil {
			fmt.Printf("Cannot add data from file %s. Likely invalid JSON. Error: %v", sourceFile, err)
			os.Exit(1)
		}

		if err = t.Insert(key, strings.TrimSpace(string(jsonBytes))); err != nil {
			fmt.Printf("Failed to write data to file new.png. Error: %s", err)
			os.Exit(1)
		}

		saveOutput(img)
	},
}

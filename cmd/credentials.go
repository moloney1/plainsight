package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"
)

func init() {
	addCmd.AddCommand(credentialsCmd)

	credentialsCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username to store")
	if err := credentialsCmd.MarkPersistentFlagRequired("username"); err != nil {
		panic(err)
	}

	credentialsCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password to store")
	if err := credentialsCmd.MarkPersistentFlagRequired("password"); err != nil {
		panic(err)
	}
}

var credentialsCmd = &cobra.Command{
	Use:     "credentials",
	Aliases: []string{"cred"},
	Short:   "Add a username/password pair to your image",
	Long:    `Add specified username/password pair to image under supplied key. Example: 'plainsight add --file myImageFile.png --key myKey credentials --username "inigo" --password "prepare2die"'`,
	Run: func(cmd *cobra.Command, args []string) {
		img, t := openOrCreateTable(imageFile)

		// TODO this logic should be moved to table.go
		if slices.Contains(t.Meta.Keys, key) {
			fmt.Printf("Image %s already has data stored under key %s. Taking no action", imageFile, key)
			os.Exit(1)
		}

		if err := t.AddUsernamePasswordPair(key, username, password); err != nil {
			fmt.Printf("Error adding data under key %s: %v", key, err)
			os.Exit(1)
		}

		saveOutput(img)
	},
}

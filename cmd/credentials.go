package cmd

import (
	"fmt"
	"hash/fnv"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/moloney1/plainsight/internal/imageio"
	"github.com/moloney1/plainsight/internal/table"
)

func init() {
	addCmd.AddCommand(credentialsCmd)

	credentialsCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username to store")
	credentialsCmd.MarkPersistentFlagRequired("username")

	credentialsCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password to store")
	credentialsCmd.MarkPersistentFlagRequired("password")

}

var credentialsCmd = &cobra.Command{
	Use:     "credentials",
	Aliases: []string{"cred"},
	Short:   "Add a username/password pair to your image",
	Long:    `Add specified username/password pair to image under supplied key. Example: 'plainsight add --file myImageFile.png --key myKey credentials --username "inigo" --password "prepare2die"'`,
	Run: func(cmd *cobra.Command, args []string) {
		img, err := imageio.ReadImage(imageFile)
		if err != nil {
			fmt.Printf("Error reading image file: %s", err)
			os.Exit(1)
		}
		hasher := fnv.New64a()
		t, err := table.TableFromBytes(img.Pix, hasher)
		if err != nil {
			fmt.Printf("Plainsight does not recognize the file. Attempting to register %s\n", imageFile)
			t, err = table.NewTable(img.Pix, hasher)
			if err != nil {
				fmt.Printf("Error registering image file %s: %v", imageFile, err)
				os.Exit(1)
			}
		}

		// TODO this logic should be moved to table.go
		if slices.Contains(t.Meta.Keys, key) {
			fmt.Printf("Image %s already has data stored under key %s. Taking no action", imageFile, key)
			os.Exit(1)
		}

		if err := t.AddUsernamePasswordPair(key, username, password); err != nil {
			fmt.Printf("Error adding data under key %s: %v", key, err)
			os.Exit(1)
		}

		imageio.WriteImageFile("new.png", img) // for now, always create a new file instead of overwriting.
		fmt.Print("Output written to file new.png")
	},
}

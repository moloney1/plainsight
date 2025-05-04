package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(queryCmd)
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Get the keys of the data stored in file",
	Long:  `Get the keys of the data stored in file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("The file you want to look at is %v\n", rootCmd.PersistentFlags().Lookup("file").Value)
	},
}

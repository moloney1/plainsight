package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var imageFile string
var key string
var jsonData string

var username string
var password string

var rootCmd = &cobra.Command{
	Use:   "plainsight",
	Short: "plainsight hides data in your favourite PNG files",
	Long:  `A tool for embedding data into PNG images and managing the embedded data`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

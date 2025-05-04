package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "plainsight",
	Short: "plainsight hides data in your favourite PNG files",
	Long:  `A tool for embedding data into PNG images and managing the embedded data`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("uh, hello?")
	},
}

func init() {
	var imageFile string
	rootCmd.PersistentFlags().StringVarP(&imageFile, "file", "f", "", "image file path")
	rootCmd.MarkPersistentFlagRequired("file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

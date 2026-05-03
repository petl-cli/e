package commands

import "github.com/spf13/cobra"

var libraryCmd = &cobra.Command{
	Use:   "library",
	Short: "",
}

func init() {
	rootCmd.AddCommand(libraryCmd)
}

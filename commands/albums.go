package commands

import "github.com/spf13/cobra"

var albumsCmd = &cobra.Command{
	Use:   "albums",
	Short: "",
}

func init() {
	rootCmd.AddCommand(albumsCmd)
}

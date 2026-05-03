package commands

import "github.com/spf13/cobra"

var artistsCmd = &cobra.Command{
	Use:   "artists",
	Short: "",
}

func init() {
	rootCmd.AddCommand(artistsCmd)
}

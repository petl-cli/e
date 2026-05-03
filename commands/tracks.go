package commands

import "github.com/spf13/cobra"

var tracksCmd = &cobra.Command{
	Use:   "tracks",
	Short: "",
}

func init() {
	rootCmd.AddCommand(tracksCmd)
}

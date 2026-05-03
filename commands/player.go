package commands

import "github.com/spf13/cobra"

var playerCmd = &cobra.Command{
	Use:   "player",
	Short: "",
}

func init() {
	rootCmd.AddCommand(playerCmd)
}

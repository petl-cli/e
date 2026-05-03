package commands

import "github.com/spf13/cobra"

var marketsCmd = &cobra.Command{
	Use:   "markets",
	Short: "",
}

func init() {
	rootCmd.AddCommand(marketsCmd)
}

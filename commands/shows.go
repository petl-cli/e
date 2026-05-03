package commands

import "github.com/spf13/cobra"

var showsCmd = &cobra.Command{
	Use:   "shows",
	Short: "",
}

func init() {
	rootCmd.AddCommand(showsCmd)
}

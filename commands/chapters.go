package commands

import "github.com/spf13/cobra"

var chaptersCmd = &cobra.Command{
	Use:   "chapters",
	Short: "",
}

func init() {
	rootCmd.AddCommand(chaptersCmd)
}

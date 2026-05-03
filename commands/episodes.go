package commands

import "github.com/spf13/cobra"

var episodesCmd = &cobra.Command{
	Use:   "episodes",
	Short: "",
}

func init() {
	rootCmd.AddCommand(episodesCmd)
}

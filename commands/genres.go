package commands

import "github.com/spf13/cobra"

var genresCmd = &cobra.Command{
	Use:   "genres",
	Short: "",
}

func init() {
	rootCmd.AddCommand(genresCmd)
}

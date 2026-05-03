package commands

import "github.com/spf13/cobra"

var audiobooksCmd = &cobra.Command{
	Use:   "audiobooks",
	Short: "",
}

func init() {
	rootCmd.AddCommand(audiobooksCmd)
}

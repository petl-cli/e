package commands

import "github.com/spf13/cobra"

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "",
}

func init() {
	rootCmd.AddCommand(usersCmd)
}

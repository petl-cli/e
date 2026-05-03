package commands

import "github.com/spf13/cobra"

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "",
}

func init() {
	rootCmd.AddCommand(categoriesCmd)
}

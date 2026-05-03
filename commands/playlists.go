package commands

import "github.com/spf13/cobra"

var playlistsCmd = &cobra.Command{
	Use:   "playlists",
	Short: "",
}

func init() {
	rootCmd.AddCommand(playlistsCmd)
}

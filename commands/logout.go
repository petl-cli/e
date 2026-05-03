package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored OAuth 2.0 credentials",
	Long:  `Deletes the locally stored OAuth token from ~/.config/spotify-web-api/token.json.`,
	RunE:  runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) error {
	cfg, err := rootConfig()
	if err != nil {
		return err
	}

	store := cfg.OAuthTokenStore()
	if store == nil {
		return fmt.Errorf("cannot determine token storage directory")
	}

	if err := store.Delete(); err != nil {
		return fmt.Errorf("removing token: %w", err)
	}

	fmt.Println("Logged out. Token removed.")
	return nil
}

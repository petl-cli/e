package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure authentication and default settings",
	Long: `Stores credentials and preferences in ~/.config/spotify-web-api/config.yaml.

Environment variables always take precedence over this file.
Supported keys: bearer_token, api_key, base_url, output_format`,
	RunE: runConfigure,
}

func init() {
	rootCmd.AddCommand(configureCmd)
}

func runConfigure(cmd *cobra.Command, args []string) error {
	path, err := _configLoader.ConfigFilePath()
	if err != nil {
		return err
	}
	fmt.Printf("Config file: %s\n\n", path)
	fmt.Println("Set credentials via environment variables or edit the config file directly.")
	fmt.Println()
	fmt.Printf("  OAuth2:       Run 'spotify-web-api login' to authenticate (client ID: export %s=<id>)\n", "SPOTIFY_WEB_API_OAUTH_CLIENT_ID")
	if cfg, cfgErr := rootConfig(); cfgErr == nil {
		if store := cfg.OAuthTokenStore(); store != nil {
			if tok := store.Load(); tok != nil {
				if tok.IsExpired() {
					fmt.Println("                Status: token expired — run 'login' to re-authenticate")
				} else {
					fmt.Printf("                Status: authenticated (token expires %s)\n", tok.ExpiresAt.Format("2006-01-02 15:04:05"))
				}
			} else {
				fmt.Println("                Status: not logged in")
			}
		}
	}
	fmt.Printf("  Base URL:     export %s_BASE_URL=<url>\n", "SPOTIFY_WEB_API")
	return nil
}

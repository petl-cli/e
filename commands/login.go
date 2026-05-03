package commands

import (
	"fmt"

	"github.com/rishimantri795/CLICreator/runtime/auth"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate via OAuth 2.0 (opens browser)",
	Long: `Starts an OAuth 2.0 Authorization Code flow with PKCE.

Opens your browser to the authorization page. After you approve access,
the CLI captures the token and stores it locally.

Tokens are stored in ~/.config/spotify-web-api/token.json (permissions 0600).
Token refresh is automatic — you should not need to re-login unless the
refresh token is revoked.`,
	RunE: runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	cfg, err := rootConfig()
	if err != nil {
		return err
	}

	clientID := cfg.OAuthClientID
	if clientID == "" {
		return fmt.Errorf("OAuth client ID is required. Set via --client-id flag or SPOTIFY_WEB_API_OAUTH_CLIENT_ID env var")
	}

	store := cfg.OAuthTokenStore()
	if store == nil {
		return fmt.Errorf("cannot determine token storage directory")
	}

	loginCfg := auth.LoginConfig{
		AuthorizeURL: "https://accounts.spotify.com/authorize",
		TokenURL:     "https://accounts.spotify.com/api/token",
		ClientID:     clientID,
		Scopes: []string{
			"app-remote-control",
			"playlist-modify-private",
			"playlist-modify-public",
			"playlist-read-collaborative",
			"playlist-read-private",
			"streaming",
			"ugc-image-upload",
			"user-follow-modify",
			"user-follow-read",
			"user-library-modify",
			"user-library-read",
			"user-modify-playback-state",
			"user-read-currently-playing",
			"user-read-email",
			"user-read-playback-position",
			"user-read-playback-state",
			"user-read-private",
			"user-read-recently-played",
			"user-top-read",
		},
	}

	tok, err := auth.Login(loginCfg)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	if err := store.Save(tok); err != nil {
		return fmt.Errorf("saving token: %w", err)
	}

	fmt.Println("Login successful! Token saved.")
	return nil
}

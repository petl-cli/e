package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var playlistsUsersCheckIfFollowsCmd = &cobra.Command{
	Use:   "users-check-if-follows",
	Short: "Check if Users Follow Playlist ",
	RunE:  withTelemetry(runPlaylistsUsersCheckIfFollows),
}

var playlistsUsersCheckIfFollowsFlags struct {
	playlistId string
	ids        string
}

func init() {
	playlistsUsersCheckIfFollowsCmd.Flags().StringVar(&playlistsUsersCheckIfFollowsFlags.playlistId, "playlist-id", "", "")
	playlistsUsersCheckIfFollowsCmd.MarkFlagRequired("playlist-id")
	playlistsUsersCheckIfFollowsCmd.Flags().StringVar(&playlistsUsersCheckIfFollowsFlags.ids, "ids", "", "")
	playlistsUsersCheckIfFollowsCmd.MarkFlagRequired("ids")

	playlistsCmd.AddCommand(playlistsUsersCheckIfFollowsCmd)
}

func runPlaylistsUsersCheckIfFollows(cmd *cobra.Command, args []string) error {
	// --schema: print full input/output type contract without making any network call.
	if rootFlags.schema {
		type flagSchema struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Required    bool   `json:"required"`
			Location    string `json:"location"`
			Description string `json:"description,omitempty"`
		}
		var flags []flagSchema
		flags = append(flags, flagSchema{
			Name:        "playlist-id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "ids",
			Type:        "string",
			Required:    true,
			Location:    "query",
			Description: "",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "200",
			ContentType: "application/json",
			Description: "Array of booleans",
		})
		responses = append(responses, responseSchema{
			Status:      "401",
			ContentType: "application/json",
			Description: "Bad or expired token. This can happen if the user revoked a token or the access token has expired. You should re-authenticate the user. ",
		})
		responses = append(responses, responseSchema{
			Status:      "403",
			ContentType: "application/json",
			Description: "Bad OAuth request (wrong consumer key, bad nonce, expired timestamp...). Unfortunately, re-authenticating the user won't help here. ",
		})
		responses = append(responses, responseSchema{
			Status:      "429",
			ContentType: "application/json",
			Description: "The app has exceeded its rate limits. ",
		})

		schema := map[string]any{
			"command":     "users-check-if-follows",
			"description": "Check if Users Follow Playlist ",
			"http": map[string]any{
				"method": "GET",
				"path":   "/playlists/{playlist_id}/followers/contains",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     false,
				"body_required": false,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         true,
				"idempotent":   true,
				"reversible":   true,
				"side_effects": []string{},
				"impact":       "low",
			},
			"requires_auth": true,
		}
		data, _ := json.MarshalIndent(schema, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	cfg, err := rootConfig()
	if err != nil {
		e := output.NetworkError(err)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	client := httpclient.New(cfg.BaseURL, cfg.AuthProvider())
	client.Debug = rootFlags.debug
	client.DryRun = rootFlags.dryRun
	if rootFlags.noRetries {
		client.RetryConfig.MaxRetries = 0
	}

	// Build path params
	pathParams := map[string]string{}
	pathParams["playlist_id"] = fmt.Sprintf("%v", playlistsUsersCheckIfFollowsFlags.playlistId)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/playlists/{playlist_id}/followers/contains", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("ids") {
		req.QueryParams["ids"] = fmt.Sprintf("%v", playlistsUsersCheckIfFollowsFlags.ids)
	}

	// Header parameters

	resp, err := client.Do(req)
	if err != nil {
		e := output.NetworkError(err)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	if resp.StatusCode >= 400 {
		e := output.HTTPError(resp.StatusCode, resp.Body)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	if rootFlags.jq != "" {
		return output.JQFilter(os.Stdout, resp.Body, rootFlags.jq)
	}
	return output.Print(os.Stdout, resp.Body, output.Format(cfg.OutputFormat))
}

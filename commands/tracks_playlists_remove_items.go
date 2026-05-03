package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var tracksPlaylistsRemoveItemsCmd = &cobra.Command{
	Use:   "playlists-remove-items",
	Short: "Remove Playlist Items ",
	RunE:  withTelemetry(runTracksPlaylistsRemoveItems),
}

var tracksPlaylistsRemoveItemsFlags struct {
	playlistId string
	tracks     []string
	snapshotId string
	body       string
}

func init() {
	tracksPlaylistsRemoveItemsCmd.Flags().StringVar(&tracksPlaylistsRemoveItemsFlags.playlistId, "playlist-id", "", "")
	tracksPlaylistsRemoveItemsCmd.MarkFlagRequired("playlist-id")
	tracksPlaylistsRemoveItemsCmd.Flags().StringSliceVar(&tracksPlaylistsRemoveItemsFlags.tracks, "tracks", nil, "An array of objects containing [Spotify URIs](/documentation/web-api/concepts/spotify-uris-ids) of the tracks or episodes to remove. For example: `{ \"tracks\": [{ \"uri\": \"spotify:track:4iV5W9uYEdYUVa79Axb7Rh\" },{ \"uri\": \"spotify:track:1301WleyT98MSxVHPZCA6M\" }] }`. A maximum of 100 objects can be sent at once. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	tracksPlaylistsRemoveItemsCmd.Flags().StringVar(&tracksPlaylistsRemoveItemsFlags.snapshotId, "snapshot-id", "", "The playlist's snapshot ID against which you want to make the changes. The API will validate that the specified items exist and in the specified positions and make the changes, even if more recent changes have been made to the playlist. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	tracksPlaylistsRemoveItemsCmd.Flags().StringVar(&tracksPlaylistsRemoveItemsFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	tracksCmd.AddCommand(tracksPlaylistsRemoveItemsCmd)
}

func runTracksPlaylistsRemoveItems(cmd *cobra.Command, args []string) error {
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
			Name:        "tracks",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "An array of objects containing [Spotify URIs](/documentation/web-api/concepts/spotify-uris-ids) of the tracks or episodes to remove. For example: `{ \"tracks\": [{ \"uri\": \"spotify:track:4iV5W9uYEdYUVa79Axb7Rh\" },{ \"uri\": \"spotify:track:1301WleyT98MSxVHPZCA6M\" }] }`. A maximum of 100 objects can be sent at once. ",
		})
		flags = append(flags, flagSchema{
			Name:        "snapshot-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The playlist's snapshot ID against which you want to make the changes. The API will validate that the specified items exist and in the specified positions and make the changes, even if more recent changes have been made to the playlist. ",
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
			Description: "A snapshot ID for the playlist",
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
			"command":     "playlists-remove-items",
			"description": "Remove Playlist Items ",
			"http": map[string]any{
				"method": "DELETE",
				"path":   "/playlists/{playlist_id}/tracks",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     true,
				"body_required": false,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         false,
				"idempotent":   true,
				"reversible":   false,
				"side_effects": []string{"destroys_resource"},
				"impact":       "high",
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
	pathParams["playlist_id"] = fmt.Sprintf("%v", tracksPlaylistsRemoveItemsFlags.playlistId)

	req := &httpclient.Request{
		Method:      "DELETE",
		Path:        httpclient.SubstitutePath("/playlists/{playlist_id}/tracks", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if tracksPlaylistsRemoveItemsFlags.body != "" {
		if err := json.Unmarshal([]byte(tracksPlaylistsRemoveItemsFlags.body), &bodyMap); err != nil {
			cliErr := &output.CLIError{
				Error:    true,
				Code:     "validation_error",
				Message:  fmt.Sprintf("invalid JSON in --body: %v", err),
				ExitCode: output.ExitValidation,
			}
			cliErr.Write(os.Stderr)
			return output.NewExitError(cliErr)
		}
	}
	// Individual flags overlay onto body (flags take precedence over --body JSON)
	if cmd.Flags().Changed("tracks") {
		bodyMap["tracks"] = tracksPlaylistsRemoveItemsFlags.tracks
	}
	if cmd.Flags().Changed("snapshot-id") {
		bodyMap["snapshot_id"] = tracksPlaylistsRemoveItemsFlags.snapshotId
	}
	req.Body = bodyMap

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

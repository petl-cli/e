package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var playlistsUpdatePlaylistItemsCmd = &cobra.Command{
	Use:   "update-playlist-items",
	Short: "Update Playlist Items ",
	RunE:  withTelemetry(runPlaylistsUpdatePlaylistItems),
}

var playlistsUpdatePlaylistItemsFlags struct {
	playlistId   string
	uris         string
	rangeStart   int
	insertBefore int
	rangeLength  int
	snapshotId   string
	body         string
}

func init() {
	playlistsUpdatePlaylistItemsCmd.Flags().StringVar(&playlistsUpdatePlaylistItemsFlags.playlistId, "playlist-id", "", "")
	playlistsUpdatePlaylistItemsCmd.MarkFlagRequired("playlist-id")
	playlistsUpdatePlaylistItemsCmd.Flags().StringVar(&playlistsUpdatePlaylistItemsFlags.uris, "uris", "", "")
	playlistsUpdatePlaylistItemsCmd.Flags().IntVar(&playlistsUpdatePlaylistItemsFlags.rangeStart, "range-start", 0, "The position of the first item to be reordered. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	playlistsUpdatePlaylistItemsCmd.Flags().IntVar(&playlistsUpdatePlaylistItemsFlags.insertBefore, "insert-before", 0, "The position where the items should be inserted.<br/>To reorder the items to the end of the playlist, simply set _insert_before_ to the position after the last item.<br/>Examples:<br/>To reorder the first item to the last position in a playlist with 10 items, set _range_start_ to 0, and _insert_before_ to 10.<br/>To reorder the last item in a playlist with 10 items to the start of the playlist, set _range_start_ to 9, and _insert_before_ to 0. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	playlistsUpdatePlaylistItemsCmd.Flags().IntVar(&playlistsUpdatePlaylistItemsFlags.rangeLength, "range-length", 0, "The amount of items to be reordered. Defaults to 1 if not set.<br/>The range of items to be reordered begins from the _range_start_ position, and includes the _range_length_ subsequent items.<br/>Example:<br/>To move the items at index 9-10 to the start of the playlist, _range_start_ is set to 9, and _range_length_ is set to 2. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	playlistsUpdatePlaylistItemsCmd.Flags().StringVar(&playlistsUpdatePlaylistItemsFlags.snapshotId, "snapshot-id", "", "The playlist's snapshot ID against which you want to make the changes. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	playlistsUpdatePlaylistItemsCmd.Flags().StringVar(&playlistsUpdatePlaylistItemsFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	playlistsCmd.AddCommand(playlistsUpdatePlaylistItemsCmd)
}

func runPlaylistsUpdatePlaylistItems(cmd *cobra.Command, args []string) error {
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
			Name:        "uris",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "range-start",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "The position of the first item to be reordered. ",
		})
		flags = append(flags, flagSchema{
			Name:        "insert-before",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "The position where the items should be inserted.<br/>To reorder the items to the end of the playlist, simply set _insert_before_ to the position after the last item.<br/>Examples:<br/>To reorder the first item to the last position in a playlist with 10 items, set _range_start_ to 0, and _insert_before_ to 10.<br/>To reorder the last item in a playlist with 10 items to the start of the playlist, set _range_start_ to 9, and _insert_before_ to 0. ",
		})
		flags = append(flags, flagSchema{
			Name:        "range-length",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "The amount of items to be reordered. Defaults to 1 if not set.<br/>The range of items to be reordered begins from the _range_start_ position, and includes the _range_length_ subsequent items.<br/>Example:<br/>To move the items at index 9-10 to the start of the playlist, _range_start_ is set to 9, and _range_length_ is set to 2. ",
		})
		flags = append(flags, flagSchema{
			Name:        "snapshot-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The playlist's snapshot ID against which you want to make the changes. ",
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
			"command":     "update-playlist-items",
			"description": "Update Playlist Items ",
			"http": map[string]any{
				"method": "PUT",
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
				"reversible":   true,
				"side_effects": []string{"mutates_resource"},
				"impact":       "medium",
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
	pathParams["playlist_id"] = fmt.Sprintf("%v", playlistsUpdatePlaylistItemsFlags.playlistId)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/playlists/{playlist_id}/tracks", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("uris") {
		req.QueryParams["uris"] = fmt.Sprintf("%v", playlistsUpdatePlaylistItemsFlags.uris)
	}

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if playlistsUpdatePlaylistItemsFlags.body != "" {
		if err := json.Unmarshal([]byte(playlistsUpdatePlaylistItemsFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("range-start") {
		bodyMap["range_start"] = playlistsUpdatePlaylistItemsFlags.rangeStart
	}
	if cmd.Flags().Changed("insert-before") {
		bodyMap["insert_before"] = playlistsUpdatePlaylistItemsFlags.insertBefore
	}
	if cmd.Flags().Changed("range-length") {
		bodyMap["range_length"] = playlistsUpdatePlaylistItemsFlags.rangeLength
	}
	if cmd.Flags().Changed("snapshot-id") {
		bodyMap["snapshot_id"] = playlistsUpdatePlaylistItemsFlags.snapshotId
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

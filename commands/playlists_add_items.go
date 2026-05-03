package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var playlistsAddItemsCmd = &cobra.Command{
	Use:   "add-items",
	Short: "Add Items to Playlist ",
	RunE:  withTelemetry(runPlaylistsAddItems),
}

var playlistsAddItemsFlags struct {
	playlistId string
	position   int
	uris       string
	body       string
}

func init() {
	playlistsAddItemsCmd.Flags().StringVar(&playlistsAddItemsFlags.playlistId, "playlist-id", "", "")
	playlistsAddItemsCmd.MarkFlagRequired("playlist-id")
	playlistsAddItemsCmd.Flags().IntVar(&playlistsAddItemsFlags.position, "position", 0, "")
	playlistsAddItemsCmd.Flags().StringVar(&playlistsAddItemsFlags.uris, "uris", "", "")
	playlistsAddItemsCmd.Flags().StringVar(&playlistsAddItemsFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	playlistsCmd.AddCommand(playlistsAddItemsCmd)
}

func runPlaylistsAddItems(cmd *cobra.Command, args []string) error {
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
			Name:        "position",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "uris",
			Type:        "string",
			Required:    false,
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
			Status:      "201",
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
			"command":     "add-items",
			"description": "Add Items to Playlist ",
			"http": map[string]any{
				"method": "POST",
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
				"idempotent":   false,
				"reversible":   true,
				"side_effects": []string{"creates_resource"},
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
	pathParams["playlist_id"] = fmt.Sprintf("%v", playlistsAddItemsFlags.playlistId)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/playlists/{playlist_id}/tracks", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("position") {
		req.QueryParams["position"] = fmt.Sprintf("%v", playlistsAddItemsFlags.position)
	}
	if cmd.Flags().Changed("uris") {
		req.QueryParams["uris"] = fmt.Sprintf("%v", playlistsAddItemsFlags.uris)
	}

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if playlistsAddItemsFlags.body != "" {
		if err := json.Unmarshal([]byte(playlistsAddItemsFlags.body), &bodyMap); err != nil {
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

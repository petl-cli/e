package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var libraryPlaylistsUpdateDetailsCmd = &cobra.Command{
	Use:   "playlists-update-details",
	Short: "Change Playlist Details ",
	RunE:  withTelemetry(runLibraryPlaylistsUpdateDetails),
}

var libraryPlaylistsUpdateDetailsFlags struct {
	playlistId    string
	description   string
	name          string
	public        bool
	collaborative bool
	body          string
}

func init() {
	libraryPlaylistsUpdateDetailsCmd.Flags().StringVar(&libraryPlaylistsUpdateDetailsFlags.playlistId, "playlist-id", "", "")
	libraryPlaylistsUpdateDetailsCmd.MarkFlagRequired("playlist-id")
	libraryPlaylistsUpdateDetailsCmd.Flags().StringVar(&libraryPlaylistsUpdateDetailsFlags.description, "description", "", "Value for playlist description as displayed in Spotify Clients and in the Web API. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	libraryPlaylistsUpdateDetailsCmd.Flags().StringVar(&libraryPlaylistsUpdateDetailsFlags.name, "name", "", "The new name for the playlist, for example `\"My New Playlist Title\"` ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	libraryPlaylistsUpdateDetailsCmd.Flags().BoolVar(&libraryPlaylistsUpdateDetailsFlags.public, "public", false, "If `true` the playlist will be public, if `false` it will be private. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	libraryPlaylistsUpdateDetailsCmd.Flags().BoolVar(&libraryPlaylistsUpdateDetailsFlags.collaborative, "collaborative", false, "If `true`, the playlist will become collaborative and other users will be able to modify the playlist in their Spotify client. <br/> _**Note**: You can only set `collaborative` to `true` on non-public playlists._ ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	libraryPlaylistsUpdateDetailsCmd.Flags().StringVar(&libraryPlaylistsUpdateDetailsFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	libraryCmd.AddCommand(libraryPlaylistsUpdateDetailsCmd)
}

func runLibraryPlaylistsUpdateDetails(cmd *cobra.Command, args []string) error {
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
			Name:        "description",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Value for playlist description as displayed in Spotify Clients and in the Web API. ",
		})
		flags = append(flags, flagSchema{
			Name:        "name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The new name for the playlist, for example `\"My New Playlist Title\"` ",
		})
		flags = append(flags, flagSchema{
			Name:        "public",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "If `true` the playlist will be public, if `false` it will be private. ",
		})
		flags = append(flags, flagSchema{
			Name:        "collaborative",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "If `true`, the playlist will become collaborative and other users will be able to modify the playlist in their Spotify client. <br/> _**Note**: You can only set `collaborative` to `true` on non-public playlists._ ",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "200",
			ContentType: "",
			Description: "Playlist updated",
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
			"command":     "playlists-update-details",
			"description": "Change Playlist Details ",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/playlists/{playlist_id}",
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
	pathParams["playlist_id"] = fmt.Sprintf("%v", libraryPlaylistsUpdateDetailsFlags.playlistId)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/playlists/{playlist_id}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if libraryPlaylistsUpdateDetailsFlags.body != "" {
		if err := json.Unmarshal([]byte(libraryPlaylistsUpdateDetailsFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("description") {
		bodyMap["description"] = libraryPlaylistsUpdateDetailsFlags.description
	}
	if cmd.Flags().Changed("name") {
		bodyMap["name"] = libraryPlaylistsUpdateDetailsFlags.name
	}
	if cmd.Flags().Changed("public") {
		bodyMap["public"] = libraryPlaylistsUpdateDetailsFlags.public
	}
	if cmd.Flags().Changed("collaborative") {
		bodyMap["collaborative"] = libraryPlaylistsUpdateDetailsFlags.collaborative
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

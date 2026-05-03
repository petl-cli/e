package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var libraryPlaylistsCreatePlaylistCmd = &cobra.Command{
	Use:   "playlists-create-playlist",
	Short: "Create Playlist ",
	RunE:  withTelemetry(runLibraryPlaylistsCreatePlaylist),
}

var libraryPlaylistsCreatePlaylistFlags struct {
	userId        string
	description   string
	name          string
	public        bool
	collaborative bool
	body          string
}

func init() {
	libraryPlaylistsCreatePlaylistCmd.Flags().StringVar(&libraryPlaylistsCreatePlaylistFlags.userId, "user-id", "", "")
	libraryPlaylistsCreatePlaylistCmd.MarkFlagRequired("user-id")
	libraryPlaylistsCreatePlaylistCmd.Flags().StringVar(&libraryPlaylistsCreatePlaylistFlags.description, "description", "", "value for playlist description as displayed in Spotify Clients and in the Web API. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	libraryPlaylistsCreatePlaylistCmd.Flags().StringVar(&libraryPlaylistsCreatePlaylistFlags.name, "name", "", "The name for the new playlist, for example `\"Your Coolest Playlist\"`. This name does not need to be unique; a user may have several playlists with the same name. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	libraryPlaylistsCreatePlaylistCmd.Flags().BoolVar(&libraryPlaylistsCreatePlaylistFlags.public, "public", false, "Defaults to `true`. If `true` the playlist will be public, if `false` it will be private. To be able to create private playlists, the user must have granted the `playlist-modify-private` [scope](/documentation/web-api/concepts/scopes/#list-of-scopes) ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	libraryPlaylistsCreatePlaylistCmd.Flags().BoolVar(&libraryPlaylistsCreatePlaylistFlags.collaborative, "collaborative", false, "Defaults to `false`. If `true` the playlist will be collaborative. _**Note**: to create a collaborative playlist you must also set `public` to `false`. To create collaborative playlists you must have granted `playlist-modify-private` and `playlist-modify-public` [scopes](/documentation/web-api/concepts/scopes/#list-of-scopes)._ ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	libraryPlaylistsCreatePlaylistCmd.Flags().StringVar(&libraryPlaylistsCreatePlaylistFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	libraryCmd.AddCommand(libraryPlaylistsCreatePlaylistCmd)
}

func runLibraryPlaylistsCreatePlaylist(cmd *cobra.Command, args []string) error {
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
			Name:        "user-id",
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
			Description: "value for playlist description as displayed in Spotify Clients and in the Web API. ",
		})
		flags = append(flags, flagSchema{
			Name:        "name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The name for the new playlist, for example `\"Your Coolest Playlist\"`. This name does not need to be unique; a user may have several playlists with the same name. ",
		})
		flags = append(flags, flagSchema{
			Name:        "public",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Defaults to `true`. If `true` the playlist will be public, if `false` it will be private. To be able to create private playlists, the user must have granted the `playlist-modify-private` [scope](/documentation/web-api/concepts/scopes/#list-of-scopes) ",
		})
		flags = append(flags, flagSchema{
			Name:        "collaborative",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Defaults to `false`. If `true` the playlist will be collaborative. _**Note**: to create a collaborative playlist you must also set `public` to `false`. To create collaborative playlists you must have granted `playlist-modify-private` and `playlist-modify-public` [scopes](/documentation/web-api/concepts/scopes/#list-of-scopes)._ ",
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
			Description: "A playlist",
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
			"command":     "playlists-create-playlist",
			"description": "Create Playlist ",
			"http": map[string]any{
				"method": "POST",
				"path":   "/users/{user_id}/playlists",
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
	pathParams["user_id"] = fmt.Sprintf("%v", libraryPlaylistsCreatePlaylistFlags.userId)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/users/{user_id}/playlists", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if libraryPlaylistsCreatePlaylistFlags.body != "" {
		if err := json.Unmarshal([]byte(libraryPlaylistsCreatePlaylistFlags.body), &bodyMap); err != nil {
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
		bodyMap["description"] = libraryPlaylistsCreatePlaylistFlags.description
	}
	if cmd.Flags().Changed("name") {
		bodyMap["name"] = libraryPlaylistsCreatePlaylistFlags.name
	}
	if cmd.Flags().Changed("public") {
		bodyMap["public"] = libraryPlaylistsCreatePlaylistFlags.public
	}
	if cmd.Flags().Changed("collaborative") {
		bodyMap["collaborative"] = libraryPlaylistsCreatePlaylistFlags.collaborative
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

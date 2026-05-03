package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var artistsUsersUnfollowArtistsUsersCmd = &cobra.Command{
	Use:   "users-unfollow-artists-users",
	Short: "Unfollow Artists or Users ",
	RunE:  withTelemetry(runArtistsUsersUnfollowArtistsUsers),
}

var artistsUsersUnfollowArtistsUsersFlags struct {
	type_ string
	ids   string
	body  string
}

func init() {
	artistsUsersUnfollowArtistsUsersCmd.Flags().StringVar(&artistsUsersUnfollowArtistsUsersFlags.type_, "type", "", "")
	artistsUsersUnfollowArtistsUsersCmd.MarkFlagRequired("type")
	artistsUsersUnfollowArtistsUsersCmd.Flags().StringVar(&artistsUsersUnfollowArtistsUsersFlags.ids, "ids", "", "")
	artistsUsersUnfollowArtistsUsersCmd.MarkFlagRequired("ids")
	artistsUsersUnfollowArtistsUsersCmd.Flags().StringVar(&artistsUsersUnfollowArtistsUsersFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	artistsCmd.AddCommand(artistsUsersUnfollowArtistsUsersCmd)
}

func runArtistsUsersUnfollowArtistsUsers(cmd *cobra.Command, args []string) error {
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
			Name:        "type",
			Type:        "string",
			Required:    true,
			Location:    "query",
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
			ContentType: "",
			Description: "Artist or user unfollowed",
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
			"command":     "users-unfollow-artists-users",
			"description": "Unfollow Artists or Users ",
			"http": map[string]any{
				"method": "DELETE",
				"path":   "/me/following",
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

	req := &httpclient.Request{
		Method:      "DELETE",
		Path:        httpclient.SubstitutePath("/me/following", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("type") {
		req.QueryParams["type"] = fmt.Sprintf("%v", artistsUsersUnfollowArtistsUsersFlags.type_)
	}
	if cmd.Flags().Changed("ids") {
		req.QueryParams["ids"] = fmt.Sprintf("%v", artistsUsersUnfollowArtistsUsersFlags.ids)
	}

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if artistsUsersUnfollowArtistsUsersFlags.body != "" {
		if err := json.Unmarshal([]byte(artistsUsersUnfollowArtistsUsersFlags.body), &bodyMap); err != nil {
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

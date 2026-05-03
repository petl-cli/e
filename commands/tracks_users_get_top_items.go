package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var tracksUsersGetTopItemsCmd = &cobra.Command{
	Use:   "users-get-top-items",
	Short: "Get User's Top Items ",
	RunE:  withTelemetry(runTracksUsersGetTopItems),
}

var tracksUsersGetTopItemsFlags struct {
	type_     string
	timeRange string
	limit     int
	offset    int
}

func init() {
	tracksUsersGetTopItemsCmd.Flags().StringVar(&tracksUsersGetTopItemsFlags.type_, "type", "", "")
	tracksUsersGetTopItemsCmd.MarkFlagRequired("type")
	tracksUsersGetTopItemsCmd.Flags().StringVar(&tracksUsersGetTopItemsFlags.timeRange, "time-range", "", "")
	tracksUsersGetTopItemsCmd.Flags().IntVar(&tracksUsersGetTopItemsFlags.limit, "limit", 0, "")
	tracksUsersGetTopItemsCmd.Flags().IntVar(&tracksUsersGetTopItemsFlags.offset, "offset", 0, "")

	tracksCmd.AddCommand(tracksUsersGetTopItemsCmd)
}

func runTracksUsersGetTopItems(cmd *cobra.Command, args []string) error {
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
			Location:    "path",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "time-range",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "offset",
			Type:        "integer",
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
			Status:      "200",
			ContentType: "application/json",
			Description: "Pages of artists or tracks",
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
			"command":     "users-get-top-items",
			"description": "Get User's Top Items ",
			"http": map[string]any{
				"method": "GET",
				"path":   "/me/top/{type}",
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
	pathParams["type"] = fmt.Sprintf("%v", tracksUsersGetTopItemsFlags.type_)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/me/top/{type}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("time-range") {
		req.QueryParams["time_range"] = fmt.Sprintf("%v", tracksUsersGetTopItemsFlags.timeRange)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", tracksUsersGetTopItemsFlags.limit)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", tracksUsersGetTopItemsFlags.offset)
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

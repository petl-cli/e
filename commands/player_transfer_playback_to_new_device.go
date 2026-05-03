package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var playerTransferPlaybackToNewDeviceCmd = &cobra.Command{
	Use:   "transfer-playback-to-new-device",
	Short: "Transfer Playback ",
	RunE:  withTelemetry(runPlayerTransferPlaybackToNewDevice),
}

var playerTransferPlaybackToNewDeviceFlags struct {
	deviceIds []string
	play      bool
	body      string
}

func init() {
	playerTransferPlaybackToNewDeviceCmd.Flags().StringSliceVar(&playerTransferPlaybackToNewDeviceFlags.deviceIds, "device-ids", nil, "A JSON array containing the ID of the device on which playback should be started/transferred.<br/>For example:`{device_ids:[\"74ASZWbe4lXaubB36ztrGX\"]}`<br/>_**Note**: Although an array is accepted, only a single device_id is currently supported. Supplying more than one will return `400 Bad Request`_ ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	playerTransferPlaybackToNewDeviceCmd.Flags().BoolVar(&playerTransferPlaybackToNewDeviceFlags.play, "play", false, "**true**: ensure playback happens on new device.<br/>**false** or not provided: keep the current playback state. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	playerTransferPlaybackToNewDeviceCmd.Flags().StringVar(&playerTransferPlaybackToNewDeviceFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	playerCmd.AddCommand(playerTransferPlaybackToNewDeviceCmd)
}

func runPlayerTransferPlaybackToNewDevice(cmd *cobra.Command, args []string) error {
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
			Name:        "device-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "A JSON array containing the ID of the device on which playback should be started/transferred.<br/>For example:`{device_ids:[\"74ASZWbe4lXaubB36ztrGX\"]}`<br/>_**Note**: Although an array is accepted, only a single device_id is currently supported. Supplying more than one will return `400 Bad Request`_ ",
		})
		flags = append(flags, flagSchema{
			Name:        "play",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "**true**: ensure playback happens on new device.<br/>**false** or not provided: keep the current playback state. ",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "204",
			ContentType: "",
			Description: "Playback transferred",
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
			"command":     "transfer-playback-to-new-device",
			"description": "Transfer Playback ",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/me/player",
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

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/me/player", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if playerTransferPlaybackToNewDeviceFlags.body != "" {
		if err := json.Unmarshal([]byte(playerTransferPlaybackToNewDeviceFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("device-ids") {
		bodyMap["device_ids"] = playerTransferPlaybackToNewDeviceFlags.deviceIds
	}
	if cmd.Flags().Changed("play") {
		bodyMap["play"] = playerTransferPlaybackToNewDeviceFlags.play
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

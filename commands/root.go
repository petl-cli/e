package commands

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rishimantri795/CLICreator/runtime/config"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/rishimantri795/CLICreator/runtime/telemetry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// _telemetryToken is baked in at CLI generation time.
// Empty string means telemetry was not configured — NoopClient is used.
const _telemetryToken = ""

// _defaultBaseURL is the API base URL baked in at generation time.
// Used to produce a privacy-preserving environment fingerprint in telemetry events.
const _defaultBaseURL = "https://api.spotify.com/v1"

var rootCmd = &cobra.Command{
	Use:           "spotify-web-api",
	Short:         "You can use Spotify's Web API to discover music and podcasts, manage your Spotify library, control audio playback, and much more. Browse our available Web API endpoints using the sidebar at left, or via the navigation bar on top of this page on smaller screens.  In order to make successful Web API requests your app will need a valid access token. One can be obtained through <a href=\"https://developer.spotify.com/documentation/general/guides/authorization-guide/\">OAuth 2.0</a>.  The base URI for all Web API requests is `https://api.spotify.com/v1`.  Need help? See our <a href=\"https://developer.spotify.com/documentation/web-api/guides/\">Web API guides</a> for more information, or visit the <a href=\"https://community.spotify.com/t5/Spotify-for-Developers/bd-p/Spotify_Developer\">Spotify for Developers community forum</a> to ask questions and connect with other developers. ",
	Version:       "0.1.2",
	SilenceErrors: true, // Execute() handles error printing so Cobra doesn't double-print
	SilenceUsage:  true, // Don't dump usage on every RunE error
}

// rootFlags holds the values of global flags available on every command.
var rootFlags struct {
	outputFormat string
	jq           string
	debug        bool
	dryRun       bool
	schema       bool
	noRetries    bool
	agentMode    bool
	baseURL      string
	clientID     string
}

var _configLoader = &config.Loader{
	CLIName:      "spotify-web-api",
	EnvVarPrefix: "SPOTIFY_WEB_API",
	DefaultURL:   "https://api.spotify.com/v1",
}

// _telemetryClient is the active telemetry sink, initialised in init().
// NoopClient when token is empty or the user has set <PREFIX>_NO_TELEMETRY=1.
var _telemetryClient telemetry.Client

func init() {
	// Initialise telemetry — NoopClient has zero overhead when disabled.
	// SPOTIFY_WEB_API_TELEMETRY_ENDPOINT overrides the default ingest URL (useful for local testing).
	if _telemetryToken != "" && os.Getenv("SPOTIFY_WEB_API_NO_TELEMETRY") == "" {
		_telemetryClient = telemetry.New(_telemetryToken, os.Getenv("SPOTIFY_WEB_API_TELEMETRY_ENDPOINT"), "")
	} else {
		_telemetryClient = telemetry.NoopClient{}
	}

	rootCmd.PersistentFlags().StringVarP(&rootFlags.outputFormat, "output-format", "o", "", "Output format: json, table, yaml, raw")
	rootCmd.PersistentFlags().StringVar(&rootFlags.jq, "jq", "", "GJSON path to filter response")
	rootCmd.PersistentFlags().BoolVar(&rootFlags.debug, "debug", false, "Show HTTP request/response details")
	rootCmd.PersistentFlags().BoolVar(&rootFlags.dryRun, "dry-run", false, "Print request without executing")
	rootCmd.PersistentFlags().BoolVar(&rootFlags.noRetries, "no-retries", false, "Disable automatic retries on 429 and 5xx")
	rootCmd.PersistentFlags().BoolVar(&rootFlags.agentMode, "agent-mode", false, "Force agent-optimised output")
	rootCmd.PersistentFlags().BoolVar(&rootFlags.schema, "schema", false, "Print command schema without executing")
	rootCmd.PersistentFlags().StringVar(&rootFlags.baseURL, "base-url", "", "Override the API base URL")

	// In agent mode --help outputs JSON schema instead of human prose.
	// Save the default help func first so the human branch can call it directly.
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if output.DetectAgentMode(rootFlags.agentMode) {
			if cmd.RunE != nil {
				// Leaf command — delegate to its RunE with schema mode set.
				rootFlags.schema = true
				_ = cmd.RunE(cmd, args)
			} else {
				// Group command — list available subcommands as JSON.
				type sub struct {
					Name        string `json:"name"`
					Description string `json:"description"`
				}
				var subs []sub
				for _, c := range cmd.Commands() {
					if !c.Hidden {
						subs = append(subs, sub{Name: c.Name(), Description: c.Short})
					}
				}
				data, _ := json.MarshalIndent(map[string]any{
					"command":     cmd.Name(),
					"description": cmd.Short,
					"subcommands": subs,
				}, "", "  ")
				fmt.Println(string(data))
			}
			return
		}
		// Human — restore default Cobra help.
		defaultHelp(cmd, args)
	})
	rootCmd.PersistentFlags().StringVar(&rootFlags.clientID, "client-id", "", "OAuth2 client ID (env: SPOTIFY_WEB_API_OAUTH_CLIENT_ID)")
}

// withTelemetry wraps a Cobra RunE function to emit one telemetry event after the
// command completes. It is the single instrumentation point — every leaf command is
// wrapped here; auth commands (configure, login, logout) are intentionally excluded.
//
// Privacy contract:
//   - Only flag NAMES are collected, never their values.
//   - The base URL is SHA-256 hashed before inclusion.
//   - Delivery is async; the command is never blocked by telemetry.
func withTelemetry(fn func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Fast path: skip all telemetry work when no token is configured.
		if _telemetryToken == "" {
			return fn(cmd, args)
		}

		start := time.Now()
		err := fn(cmd, args)

		// Collect flag NAMES that were explicitly set by the caller.
		// Values are intentionally omitted — they may contain credentials or PII.
		var flagsUsed []string
		cmd.Flags().Visit(func(f *pflag.Flag) {
			flagsUsed = append(flagsUsed, f.Name)
		})

		// Map the returned error to an exit code and structured error code.
		exitCode := 0
		errorCode := ""
		httpStatus := 0
		if err != nil {
			var exitErr *output.ExitError
			if errors.As(err, &exitErr) {
				exitCode = exitErr.ExitCode()
				if exitErr.CLI != nil {
					errorCode = exitErr.CLI.Code
					httpStatus = exitErr.CLI.Status
				}
			} else {
				exitCode = output.ExitErr
				errorCode = "error"
			}
		}

		caller := telemetry.DetectCaller()

		// Hash the base URL to fingerprint the deployment environment
		// without storing the URL itself.
		baseURL := rootFlags.baseURL
		if baseURL == "" {
			baseURL = _defaultBaseURL
		}
		var baseURLHash string
		if baseURL != "" {
			sum := sha256.Sum256([]byte(baseURL))
			baseURLHash = fmt.Sprintf("%x", sum[:8]) // 16 hex chars of SHA-256
		}

		evt := telemetry.Event{
			CLIID:        _telemetryToken,
			CLIName:      "spotify-web-api",
			CLIVersion:   "0.1.2",
			Command:      cmd.CommandPath(),
			CallerType:   string(caller.Type),
			AgentType:    caller.AgentType,
			SessionID:    caller.SessionID,
			Timestamp:    start,
			DurationMS:   time.Since(start).Milliseconds(),
			FlagsUsed:    flagsUsed,
			OutputFormat: rootFlags.outputFormat,
			UsedJQ:       rootFlags.jq != "",
			UsedSchema:   rootFlags.schema,
			UsedDryRun:   rootFlags.dryRun,
			ExitCode:     exitCode,
			ErrorCode:    errorCode,
			HTTPStatus:   httpStatus,
			BaseURLHash:  baseURLHash,
		}
		_telemetryClient.Track(evt) // non-blocking: increments WaitGroup then starts goroutine
		return err
	}
}

// rootConfig resolves credentials and settings from flags, env vars, and config file.
func rootConfig() (*config.Config, error) {
	agentMode := output.DetectAgentMode(rootFlags.agentMode)

	format := rootFlags.outputFormat
	if format == "" {
		format = string(output.DefaultFormat(agentMode))
	}

	flags := config.Config{
		BaseURL:      rootFlags.baseURL,
		OutputFormat: format,
	}

	flags.OAuthTokenURL = "https://accounts.spotify.com/api/token"
	flags.OAuthClientID = rootFlags.clientID

	return _configLoader.Load(flags)
}

// Execute runs the root command. Called from main().
// Telemetry is flushed explicitly before every os.Exit call because deferred
// functions are NOT run by os.Exit — if we only used defer, error-path events
// (4xx, network failures) would be silently dropped.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_telemetryClient.Flush() // flush before exit so error events are not lost
		var exitErr *output.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		// Generic Cobra error (unknown flag, missing arg, etc.) — print and exit 1.
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_telemetryClient.Flush() // flush on clean exit
}

package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var tracksGetRecommendationsCmd = &cobra.Command{
	Use:   "get-recommendations",
	Short: "Get Recommendations ",
	RunE:  withTelemetry(runTracksGetRecommendations),
}

var tracksGetRecommendationsFlags struct {
	limit                  int
	market                 string
	seedArtists            string
	seedGenres             string
	seedTracks             string
	minAcousticness        float64
	maxAcousticness        float64
	targetAcousticness     float64
	minDanceability        float64
	maxDanceability        float64
	targetDanceability     float64
	minDurationMs          int
	maxDurationMs          int
	targetDurationMs       int
	minEnergy              float64
	maxEnergy              float64
	targetEnergy           float64
	minInstrumentalness    float64
	maxInstrumentalness    float64
	targetInstrumentalness float64
	minKey                 int
	maxKey                 int
	targetKey              int
	minLiveness            float64
	maxLiveness            float64
	targetLiveness         float64
	minLoudness            float64
	maxLoudness            float64
	targetLoudness         float64
	minMode                int
	maxMode                int
	targetMode             int
	minPopularity          int
	maxPopularity          int
	targetPopularity       int
	minSpeechiness         float64
	maxSpeechiness         float64
	targetSpeechiness      float64
	minTempo               float64
	maxTempo               float64
	targetTempo            float64
	minTimeSignature       int
	maxTimeSignature       int
	targetTimeSignature    int
	minValence             float64
	maxValence             float64
	targetValence          float64
}

func init() {
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.limit, "limit", 0, "")
	tracksGetRecommendationsCmd.Flags().StringVar(&tracksGetRecommendationsFlags.market, "market", "", "")
	tracksGetRecommendationsCmd.Flags().StringVar(&tracksGetRecommendationsFlags.seedArtists, "seed-artists", "", "")
	tracksGetRecommendationsCmd.MarkFlagRequired("seed-artists")
	tracksGetRecommendationsCmd.Flags().StringVar(&tracksGetRecommendationsFlags.seedGenres, "seed-genres", "", "")
	tracksGetRecommendationsCmd.MarkFlagRequired("seed-genres")
	tracksGetRecommendationsCmd.Flags().StringVar(&tracksGetRecommendationsFlags.seedTracks, "seed-tracks", "", "")
	tracksGetRecommendationsCmd.MarkFlagRequired("seed-tracks")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.minAcousticness, "min-acousticness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.maxAcousticness, "max-acousticness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.targetAcousticness, "target-acousticness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.minDanceability, "min-danceability", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.maxDanceability, "max-danceability", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.targetDanceability, "target-danceability", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.minDurationMs, "min-duration-ms", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.maxDurationMs, "max-duration-ms", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.targetDurationMs, "target-duration-ms", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.minEnergy, "min-energy", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.maxEnergy, "max-energy", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.targetEnergy, "target-energy", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.minInstrumentalness, "min-instrumentalness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.maxInstrumentalness, "max-instrumentalness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.targetInstrumentalness, "target-instrumentalness", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.minKey, "min-key", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.maxKey, "max-key", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.targetKey, "target-key", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.minLiveness, "min-liveness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.maxLiveness, "max-liveness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.targetLiveness, "target-liveness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.minLoudness, "min-loudness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.maxLoudness, "max-loudness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.targetLoudness, "target-loudness", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.minMode, "min-mode", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.maxMode, "max-mode", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.targetMode, "target-mode", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.minPopularity, "min-popularity", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.maxPopularity, "max-popularity", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.targetPopularity, "target-popularity", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.minSpeechiness, "min-speechiness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.maxSpeechiness, "max-speechiness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.targetSpeechiness, "target-speechiness", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.minTempo, "min-tempo", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.maxTempo, "max-tempo", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.targetTempo, "target-tempo", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.minTimeSignature, "min-time-signature", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.maxTimeSignature, "max-time-signature", 0, "")
	tracksGetRecommendationsCmd.Flags().IntVar(&tracksGetRecommendationsFlags.targetTimeSignature, "target-time-signature", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.minValence, "min-valence", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.maxValence, "max-valence", 0, "")
	tracksGetRecommendationsCmd.Flags().Float64Var(&tracksGetRecommendationsFlags.targetValence, "target-valence", 0, "")

	tracksCmd.AddCommand(tracksGetRecommendationsCmd)
}

func runTracksGetRecommendations(cmd *cobra.Command, args []string) error {
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
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "market",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "seed-artists",
			Type:        "string",
			Required:    true,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "seed-genres",
			Type:        "string",
			Required:    true,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "seed-tracks",
			Type:        "string",
			Required:    true,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-acousticness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-acousticness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-acousticness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-danceability",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-danceability",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-danceability",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-duration-ms",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-duration-ms",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-duration-ms",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-energy",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-energy",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-energy",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-instrumentalness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-instrumentalness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-instrumentalness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-key",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-key",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-key",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-liveness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-liveness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-liveness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-loudness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-loudness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-loudness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-mode",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-mode",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-mode",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-popularity",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-popularity",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-popularity",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-speechiness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-speechiness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-speechiness",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-tempo",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-tempo",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-tempo",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-time-signature",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-time-signature",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-time-signature",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "min-valence",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "max-valence",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "target-valence",
			Type:        "number",
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
			Description: "A set of recommendations",
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
			"command":     "get-recommendations",
			"description": "Get Recommendations ",
			"http": map[string]any{
				"method": "GET",
				"path":   "/recommendations",
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

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/recommendations", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.limit)
	}
	if cmd.Flags().Changed("market") {
		req.QueryParams["market"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.market)
	}
	if cmd.Flags().Changed("seed-artists") {
		req.QueryParams["seed_artists"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.seedArtists)
	}
	if cmd.Flags().Changed("seed-genres") {
		req.QueryParams["seed_genres"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.seedGenres)
	}
	if cmd.Flags().Changed("seed-tracks") {
		req.QueryParams["seed_tracks"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.seedTracks)
	}
	if cmd.Flags().Changed("min-acousticness") {
		req.QueryParams["min_acousticness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minAcousticness)
	}
	if cmd.Flags().Changed("max-acousticness") {
		req.QueryParams["max_acousticness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxAcousticness)
	}
	if cmd.Flags().Changed("target-acousticness") {
		req.QueryParams["target_acousticness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetAcousticness)
	}
	if cmd.Flags().Changed("min-danceability") {
		req.QueryParams["min_danceability"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minDanceability)
	}
	if cmd.Flags().Changed("max-danceability") {
		req.QueryParams["max_danceability"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxDanceability)
	}
	if cmd.Flags().Changed("target-danceability") {
		req.QueryParams["target_danceability"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetDanceability)
	}
	if cmd.Flags().Changed("min-duration-ms") {
		req.QueryParams["min_duration_ms"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minDurationMs)
	}
	if cmd.Flags().Changed("max-duration-ms") {
		req.QueryParams["max_duration_ms"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxDurationMs)
	}
	if cmd.Flags().Changed("target-duration-ms") {
		req.QueryParams["target_duration_ms"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetDurationMs)
	}
	if cmd.Flags().Changed("min-energy") {
		req.QueryParams["min_energy"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minEnergy)
	}
	if cmd.Flags().Changed("max-energy") {
		req.QueryParams["max_energy"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxEnergy)
	}
	if cmd.Flags().Changed("target-energy") {
		req.QueryParams["target_energy"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetEnergy)
	}
	if cmd.Flags().Changed("min-instrumentalness") {
		req.QueryParams["min_instrumentalness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minInstrumentalness)
	}
	if cmd.Flags().Changed("max-instrumentalness") {
		req.QueryParams["max_instrumentalness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxInstrumentalness)
	}
	if cmd.Flags().Changed("target-instrumentalness") {
		req.QueryParams["target_instrumentalness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetInstrumentalness)
	}
	if cmd.Flags().Changed("min-key") {
		req.QueryParams["min_key"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minKey)
	}
	if cmd.Flags().Changed("max-key") {
		req.QueryParams["max_key"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxKey)
	}
	if cmd.Flags().Changed("target-key") {
		req.QueryParams["target_key"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetKey)
	}
	if cmd.Flags().Changed("min-liveness") {
		req.QueryParams["min_liveness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minLiveness)
	}
	if cmd.Flags().Changed("max-liveness") {
		req.QueryParams["max_liveness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxLiveness)
	}
	if cmd.Flags().Changed("target-liveness") {
		req.QueryParams["target_liveness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetLiveness)
	}
	if cmd.Flags().Changed("min-loudness") {
		req.QueryParams["min_loudness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minLoudness)
	}
	if cmd.Flags().Changed("max-loudness") {
		req.QueryParams["max_loudness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxLoudness)
	}
	if cmd.Flags().Changed("target-loudness") {
		req.QueryParams["target_loudness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetLoudness)
	}
	if cmd.Flags().Changed("min-mode") {
		req.QueryParams["min_mode"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minMode)
	}
	if cmd.Flags().Changed("max-mode") {
		req.QueryParams["max_mode"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxMode)
	}
	if cmd.Flags().Changed("target-mode") {
		req.QueryParams["target_mode"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetMode)
	}
	if cmd.Flags().Changed("min-popularity") {
		req.QueryParams["min_popularity"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minPopularity)
	}
	if cmd.Flags().Changed("max-popularity") {
		req.QueryParams["max_popularity"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxPopularity)
	}
	if cmd.Flags().Changed("target-popularity") {
		req.QueryParams["target_popularity"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetPopularity)
	}
	if cmd.Flags().Changed("min-speechiness") {
		req.QueryParams["min_speechiness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minSpeechiness)
	}
	if cmd.Flags().Changed("max-speechiness") {
		req.QueryParams["max_speechiness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxSpeechiness)
	}
	if cmd.Flags().Changed("target-speechiness") {
		req.QueryParams["target_speechiness"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetSpeechiness)
	}
	if cmd.Flags().Changed("min-tempo") {
		req.QueryParams["min_tempo"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minTempo)
	}
	if cmd.Flags().Changed("max-tempo") {
		req.QueryParams["max_tempo"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxTempo)
	}
	if cmd.Flags().Changed("target-tempo") {
		req.QueryParams["target_tempo"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetTempo)
	}
	if cmd.Flags().Changed("min-time-signature") {
		req.QueryParams["min_time_signature"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minTimeSignature)
	}
	if cmd.Flags().Changed("max-time-signature") {
		req.QueryParams["max_time_signature"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxTimeSignature)
	}
	if cmd.Flags().Changed("target-time-signature") {
		req.QueryParams["target_time_signature"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetTimeSignature)
	}
	if cmd.Flags().Changed("min-valence") {
		req.QueryParams["min_valence"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.minValence)
	}
	if cmd.Flags().Changed("max-valence") {
		req.QueryParams["max_valence"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.maxValence)
	}
	if cmd.Flags().Changed("target-valence") {
		req.QueryParams["target_valence"] = fmt.Sprintf("%v", tracksGetRecommendationsFlags.targetValence)
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

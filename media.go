package promptrails

import "context"

// MediaService handles media generation API calls.
type MediaService struct {
	http *httpClient
}

// GenerateMediaParams are parameters for generating media.
type GenerateMediaParams struct {
	Provider   string         `json:"provider"`
	MediaType  string         `json:"media_type"`
	Model      string         `json:"model"`
	Prompt     string         `json:"prompt,omitempty"`
	InputURL   string         `json:"input_url,omitempty"`
	Config     map[string]any `json:"config,omitempty"`
	Guardrails map[string]any `json:"guardrails,omitempty"`
}

// GenerateMediaResponse is the response from generating media.
type GenerateMediaResponse struct {
	Status      string         `json:"status"`
	JobID       string         `json:"job_id,omitempty"`
	AssetURL    string         `json:"asset_url,omitempty"`
	TextOutput  string         `json:"text_output,omitempty"`
	ContentType string         `json:"content_type,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// Generate creates media content using the specified provider and model.
func (s *MediaService) Generate(ctx context.Context, params *GenerateMediaParams) (*GenerateMediaResponse, error) {
	var result GenerateMediaResponse
	err := s.http.post(ctx, "/api/v1/media/generate", params, &result)
	return &result, err
}

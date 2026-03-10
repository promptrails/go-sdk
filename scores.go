package promptrails

import (
	"context"
	"fmt"
)

// ScoresService handles score-related API calls.
type ScoresService struct {
	http *httpClient
}

// ListScoresParams are parameters for listing scores.
type ListScoresParams struct {
	Page          int
	Limit         int
	ExecutionID   string
	ScoreConfigID string
}

// CreateScoreParams are parameters for creating a score.
type CreateScoreParams struct {
	ExecutionID   string         `json:"execution_id"`
	ScoreConfigID string         `json:"score_config_id"`
	Value         float64        `json:"value"`
	Comment       string         `json:"comment,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

// UpdateScoreParams are parameters for updating a score.
type UpdateScoreParams struct {
	Value    *float64       `json:"value,omitempty"`
	Comment  *string        `json:"comment,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// CreateScoreConfigParams are parameters for creating a score config.
type CreateScoreConfigParams struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	DataType    string   `json:"data_type"`
	MinValue    *float64 `json:"min_value,omitempty"`
	MaxValue    *float64 `json:"max_value,omitempty"`
}

// UpdateScoreConfigParams are parameters for updating a score config.
type UpdateScoreConfigParams struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	MinValue    *float64 `json:"min_value,omitempty"`
	MaxValue    *float64 `json:"max_value,omitempty"`
}

// List returns a paginated list of scores.
func (s *ScoresService) List(ctx context.Context, params *ListScoresParams) (*PaginatedResponse[Score], error) {
	if params == nil {
		params = &ListScoresParams{}
	}
	p := &ListParams{Page: params.Page, Limit: params.Limit}
	p.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", p.Page),
		"limit": fmt.Sprintf("%d", p.Limit),
	}
	if params.ExecutionID != "" {
		qp["execution_id"] = params.ExecutionID
	}
	if params.ScoreConfigID != "" {
		qp["score_config_id"] = params.ScoreConfigID
	}
	var result PaginatedResponse[Score]
	err := s.http.get(ctx, "/api/v1/scores", qp, &result)
	return &result, err
}

// Get returns a single score by ID.
func (s *ScoresService) Get(ctx context.Context, id string) (*Score, error) {
	var result Score
	err := s.http.get(ctx, "/api/v1/scores/"+id, nil, &result)
	return &result, err
}

// Create creates a new score.
func (s *ScoresService) Create(ctx context.Context, params *CreateScoreParams) (*Score, error) {
	var result Score
	err := s.http.post(ctx, "/api/v1/scores", params, &result)
	return &result, err
}

// Update updates an existing score.
func (s *ScoresService) Update(ctx context.Context, id string, params *UpdateScoreParams) (*Score, error) {
	var result Score
	err := s.http.patch(ctx, "/api/v1/scores/"+id, params, &result)
	return &result, err
}

// Delete removes a score.
func (s *ScoresService) Delete(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/scores/"+id)
}

// Aggregates returns aggregated score data.
func (s *ScoresService) Aggregates(ctx context.Context, params *ListParams) ([]ScoreAggregate, error) {
	if params == nil {
		params = &ListParams{}
	}
	params.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", params.Page),
		"limit": fmt.Sprintf("%d", params.Limit),
	}
	var result []ScoreAggregate
	err := s.http.get(ctx, "/api/v1/scores/aggregates", qp, &result)
	return result, err
}

// ListConfigs returns a paginated list of score configs.
func (s *ScoresService) ListConfigs(ctx context.Context, params *ListParams) (*PaginatedResponse[ScoreConfig], error) {
	if params == nil {
		params = &ListParams{}
	}
	params.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", params.Page),
		"limit": fmt.Sprintf("%d", params.Limit),
	}
	var result PaginatedResponse[ScoreConfig]
	err := s.http.get(ctx, "/api/v1/scores/configs", qp, &result)
	return &result, err
}

// GetConfig returns a single score config.
func (s *ScoresService) GetConfig(ctx context.Context, id string) (*ScoreConfig, error) {
	var result ScoreConfig
	err := s.http.get(ctx, "/api/v1/scores/configs/"+id, nil, &result)
	return &result, err
}

// CreateConfig creates a new score config.
func (s *ScoresService) CreateConfig(ctx context.Context, params *CreateScoreConfigParams) (*ScoreConfig, error) {
	var result ScoreConfig
	err := s.http.post(ctx, "/api/v1/scores/configs", params, &result)
	return &result, err
}

// UpdateConfig updates an existing score config.
func (s *ScoresService) UpdateConfig(ctx context.Context, id string, params *UpdateScoreConfigParams) (*ScoreConfig, error) {
	var result ScoreConfig
	err := s.http.patch(ctx, "/api/v1/scores/configs/"+id, params, &result)
	return &result, err
}

// DeleteConfig removes a score config.
func (s *ScoresService) DeleteConfig(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/scores/configs/"+id)
}

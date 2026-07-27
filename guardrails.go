package promptrails

import "context"

// GuardrailsService handles workspace-level guardrail API calls. Guardrails
// are attached to an agent version (see CreateVersionParams.Guardrails) or an
// agent (see AgentsService.CreateGuardrail); this service manages the
// available scanner catalog and individual attached guardrails.
type GuardrailsService struct {
	http *httpClient
}

// UpdateGuardrailParams are parameters for updating an attached guardrail.
type UpdateGuardrailParams struct {
	Action    *string        `json:"action,omitempty"`
	Config    map[string]any `json:"config,omitempty"`
	IsActive  *bool          `json:"is_active,omitempty"`
	SortOrder *int           `json:"sort_order,omitempty"`
}

// ListScanners lists the guardrail scanners available in this workspace.
func (s *GuardrailsService) ListScanners(ctx context.Context) ([]ScannerMeta, error) {
	var result []ScannerMeta
	err := s.http.get(ctx, "/api/v1/guardrails/scanners", nil, &result)
	return result, err
}

// Update modifies an attached guardrail.
func (s *GuardrailsService) Update(ctx context.Context, id string, params *UpdateGuardrailParams) (*Guardrail, error) {
	var result Guardrail
	err := s.http.patch(ctx, "/api/v1/guardrails/"+id, params, &result)
	return &result, err
}

// Delete removes an attached guardrail.
func (s *GuardrailsService) Delete(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/guardrails/"+id)
}

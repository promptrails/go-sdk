package promptrails

import "context"

// CostsService handles cost-related API calls.
type CostsService struct {
	http *httpClient
}

// CostParams are parameters for cost queries.
type CostParams struct {
	From string
	To   string
}

// GetSummary returns cost summary for the workspace.
func (s *CostsService) GetSummary(ctx context.Context, params *CostParams) (*CostSummary, error) {
	qp := map[string]string{}
	if params != nil {
		if params.From != "" {
			qp["from"] = params.From
		}
		if params.To != "" {
			qp["to"] = params.To
		}
	}
	var result CostSummary
	err := s.http.get(ctx, "/api/v1/costs/summary", qp, &result)
	return &result, err
}

// GetAgentSummary returns cost summary for a specific agent.
func (s *CostsService) GetAgentSummary(ctx context.Context, agentID string, params *CostParams) (*CostSummary, error) {
	qp := map[string]string{}
	if params != nil {
		if params.From != "" {
			qp["from"] = params.From
		}
		if params.To != "" {
			qp["to"] = params.To
		}
	}
	var result CostSummary
	err := s.http.get(ctx, "/api/v1/costs/agents/"+agentID+"/summary", qp, &result)
	return &result, err
}

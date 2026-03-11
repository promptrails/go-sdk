package promptrails

import (
	"context"
	"fmt"
)

// AssetsService handles asset-related API calls.
type AssetsService struct {
	http *httpClient
}

// ListAssetsParams are parameters for listing assets.
type ListAssetsParams struct {
	Page      int
	Limit     int
	MediaType string
	Provider  string
}

// List returns a paginated list of assets.
func (s *AssetsService) List(ctx context.Context, params *ListAssetsParams) (*PaginatedResponse[Asset], error) {
	if params == nil {
		params = &ListAssetsParams{}
	}
	p := &ListParams{Page: params.Page, Limit: params.Limit}
	p.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", p.Page),
		"limit": fmt.Sprintf("%d", p.Limit),
	}
	if params.MediaType != "" {
		qp["media_type"] = params.MediaType
	}
	if params.Provider != "" {
		qp["provider"] = params.Provider
	}
	var result PaginatedResponse[Asset]
	err := s.http.get(ctx, "/api/v1/assets", qp, &result)
	return &result, err
}

// Get returns a single asset by ID.
func (s *AssetsService) Get(ctx context.Context, id string) (*Asset, error) {
	var result Asset
	err := s.http.get(ctx, "/api/v1/assets/"+id, nil, &result)
	return &result, err
}

// Delete removes an asset.
func (s *AssetsService) Delete(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/assets/"+id)
}

// GetSignedURL returns a signed URL for accessing an asset.
func (s *AssetsService) GetSignedURL(ctx context.Context, id string) (*SignedURLResponse, error) {
	var result SignedURLResponse
	err := s.http.get(ctx, "/api/v1/assets/"+id+"/signed-url", nil, &result)
	return &result, err
}

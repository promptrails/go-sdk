package promptrails

import (
	"context"
	"fmt"
)

// DataSourcesService handles data source API calls.
type DataSourcesService struct {
	http *httpClient
}

// ListDataSourcesParams are parameters for listing data sources.
type ListDataSourcesParams struct {
	Page  int
	Limit int
}

// CreateDataSourceParams are parameters for creating a data source.
type CreateDataSourceParams struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
}

// UpdateDataSourceParams are parameters for updating a data source.
type UpdateDataSourceParams struct {
	Name        *string        `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	Config      map[string]any `json:"config,omitempty"`
}

// QueryDataSourceParams are parameters for querying a data source.
type QueryDataSourceParams struct {
	Query string         `json:"query"`
	Params map[string]any `json:"params,omitempty"`
}

// List returns a paginated list of data sources.
func (s *DataSourcesService) List(ctx context.Context, params *ListDataSourcesParams) (*PaginatedResponse[DataSource], error) {
	if params == nil {
		params = &ListDataSourcesParams{}
	}
	p := &ListParams{Page: params.Page, Limit: params.Limit}
	p.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", p.Page),
		"limit": fmt.Sprintf("%d", p.Limit),
	}
	var result PaginatedResponse[DataSource]
	err := s.http.get(ctx, "/api/v1/data-sources", qp, &result)
	return &result, err
}

// Get returns a single data source by ID.
func (s *DataSourcesService) Get(ctx context.Context, id string) (*DataSource, error) {
	var result DataSource
	err := s.http.get(ctx, "/api/v1/data-sources/"+id, nil, &result)
	return &result, err
}

// Create creates a new data source.
func (s *DataSourcesService) Create(ctx context.Context, params *CreateDataSourceParams) (*DataSource, error) {
	var result DataSource
	err := s.http.post(ctx, "/api/v1/data-sources", params, &result)
	return &result, err
}

// Update updates an existing data source.
func (s *DataSourcesService) Update(ctx context.Context, id string, params *UpdateDataSourceParams) (*DataSource, error) {
	var result DataSource
	err := s.http.patch(ctx, "/api/v1/data-sources/"+id, params, &result)
	return &result, err
}

// Delete removes a data source.
func (s *DataSourcesService) Delete(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/data-sources/"+id)
}

// ListVersions returns versions for a data source.
func (s *DataSourcesService) ListVersions(ctx context.Context, id string, params *ListParams) (*PaginatedResponse[DataSourceVersion], error) {
	if params == nil {
		params = &ListParams{}
	}
	params.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", params.Page),
		"limit": fmt.Sprintf("%d", params.Limit),
	}
	var result PaginatedResponse[DataSourceVersion]
	err := s.http.get(ctx, "/api/v1/data-sources/"+id+"/versions", qp, &result)
	return &result, err
}

// CreateVersion creates a new version for a data source.
func (s *DataSourcesService) CreateVersion(ctx context.Context, id string, params *CreateVersionParams) (*DataSourceVersion, error) {
	var result DataSourceVersion
	err := s.http.post(ctx, "/api/v1/data-sources/"+id+"/versions", params, &result)
	return &result, err
}

// TestConnection tests the data source connection.
func (s *DataSourcesService) TestConnection(ctx context.Context, id string) error {
	return s.http.post(ctx, "/api/v1/data-sources/"+id+"/test", nil, nil)
}

// Query executes a query against the data source.
func (s *DataSourcesService) Query(ctx context.Context, id string, params *QueryDataSourceParams) (any, error) {
	var result any
	err := s.http.post(ctx, "/api/v1/data-sources/"+id+"/query", params, &result)
	return result, err
}

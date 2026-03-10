package promptrails

import (
	"context"
	"fmt"
)

// CredentialsService handles credential-related API calls.
type CredentialsService struct {
	http *httpClient
}

// ListCredentialsParams are parameters for listing credentials.
type ListCredentialsParams struct {
	Page  int
	Limit int
}

// CreateCredentialParams are parameters for creating a credential.
type CreateCredentialParams struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
}

// UpdateCredentialParams are parameters for updating a credential.
type UpdateCredentialParams struct {
	Name   *string `json:"name,omitempty"`
	APIKey *string `json:"api_key,omitempty"`
}

// List returns a paginated list of credentials.
func (s *CredentialsService) List(ctx context.Context, params *ListCredentialsParams) (*PaginatedResponse[Credential], error) {
	if params == nil {
		params = &ListCredentialsParams{}
	}
	p := &ListParams{Page: params.Page, Limit: params.Limit}
	p.defaults()
	qp := map[string]string{
		"page":  fmt.Sprintf("%d", p.Page),
		"limit": fmt.Sprintf("%d", p.Limit),
	}
	var result PaginatedResponse[Credential]
	err := s.http.get(ctx, "/api/v1/credentials", qp, &result)
	return &result, err
}

// Get returns a single credential by ID.
func (s *CredentialsService) Get(ctx context.Context, id string) (*Credential, error) {
	var result Credential
	err := s.http.get(ctx, "/api/v1/credentials/"+id, nil, &result)
	return &result, err
}

// Create creates a new credential.
func (s *CredentialsService) Create(ctx context.Context, params *CreateCredentialParams) (*Credential, error) {
	var result Credential
	err := s.http.post(ctx, "/api/v1/credentials", params, &result)
	return &result, err
}

// Update updates an existing credential.
func (s *CredentialsService) Update(ctx context.Context, id string, params *UpdateCredentialParams) (*Credential, error) {
	var result Credential
	err := s.http.patch(ctx, "/api/v1/credentials/"+id, params, &result)
	return &result, err
}

// Delete removes a credential.
func (s *CredentialsService) Delete(ctx context.Context, id string) error {
	return s.http.del(ctx, "/api/v1/credentials/"+id)
}

// SetDefault sets a credential as the default for its provider.
func (s *CredentialsService) SetDefault(ctx context.Context, id string) error {
	return s.http.post(ctx, "/api/v1/credentials/"+id+"/default", nil, nil)
}

// CheckConnection tests the credential connection.
func (s *CredentialsService) CheckConnection(ctx context.Context, id string) (*Credential, error) {
	var result Credential
	err := s.http.post(ctx, "/api/v1/credentials/"+id+"/check", nil, &result)
	return &result, err
}

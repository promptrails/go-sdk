package promptrails

import "time"

const defaultBaseURL = "https://api.promptrails.ai"

type config struct {
	apiKey     string
	baseURL    string
	timeout    time.Duration
	maxRetries int
}

// Option configures the client.
type Option func(*config)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) Option {
	return func(c *config) { c.baseURL = url }
}

// WithTimeout sets the HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *config) { c.timeout = d }
}

// WithMaxRetries sets the maximum number of retries on 5xx/network errors.
func WithMaxRetries(n int) Option {
	return func(c *config) { c.maxRetries = n }
}

func defaultConfig() config {
	return config{
		baseURL:    defaultBaseURL,
		timeout:    30 * time.Second,
		maxRetries: 3,
	}
}

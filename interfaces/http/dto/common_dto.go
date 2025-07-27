package dto

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// TestBBSProviderRequest represents the request to test a BBS provider
type TestBBSProviderRequest struct {
	Provider string `json:"provider" validate:"required"`
}

// TestBBSProviderResponse represents the response from testing a BBS provider
type TestBBSProviderResponse struct {
	Provider    string `json:"provider"`
	Available   bool   `json:"available"`
	Performance string `json:"performance,omitempty"`
	Message     string `json:"message,omitempty"`
}

// BenchmarkBBSProvidersRequest represents the request to benchmark BBS providers
type BenchmarkBBSProvidersRequest struct {
	Providers []string `json:"providers" validate:"required,min=1"`
	Messages  int      `json:"messages,omitempty"` // Default to 5 if not specified
}

// BenchmarkResult represents a single benchmark result
type BenchmarkResult struct {
	Provider        string  `json:"provider"`
	Available       bool    `json:"available"`
	SignTime        float64 `json:"signTime"`        // milliseconds
	VerifyTime      float64 `json:"verifyTime"`      // milliseconds
	ProofCreateTime float64 `json:"proofCreateTime"` // milliseconds
	ProofVerifyTime float64 `json:"proofVerifyTime"` // milliseconds
	Message         string  `json:"message,omitempty"`
}

// BenchmarkBBSProvidersResponse represents the response from benchmarking BBS providers
type BenchmarkBBSProvidersResponse struct {
	Results []BenchmarkResult `json:"results"`
	Summary string            `json:"summary"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

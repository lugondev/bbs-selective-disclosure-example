package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lugondev/bbs-selective-disclosure-example/interfaces/http/dto"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/bbs"
)

// BBSHandler handles BBS provider testing and benchmarking
type BBSHandler struct {
	factory bbs.BBSServiceFactory
}

// NewBBSHandler creates a new BBS handler
func NewBBSHandler(factory bbs.BBSServiceFactory) *BBSHandler {
	return &BBSHandler{
		factory: factory,
	}
}

// TestProvider handles POST /api/bbs/test
func (h *BBSHandler) TestProvider(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	var req dto.TestBBSProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Parse provider
	provider, err := bbs.ParseProvider(req.Provider)
	if err != nil {
		writeErrorResponse(w, "Invalid provider", http.StatusBadRequest, err.Error())
		return
	}

	// Test the provider
	response := h.testSingleProvider(provider)
	writeSuccessResponse(w, response)
}

// BenchmarkProviders handles POST /api/bbs/benchmark
func (h *BBSHandler) BenchmarkProviders(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, "")
		return
	}

	var req dto.BenchmarkBBSProvidersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	// Default to 5 messages if not specified
	messageCount := req.Messages
	if messageCount <= 0 {
		messageCount = 5
	}

	// Benchmark each provider
	results := make([]dto.BenchmarkResult, 0, len(req.Providers))
	for _, providerStr := range req.Providers {
		provider, err := bbs.ParseProvider(providerStr)
		if err != nil {
			results = append(results, dto.BenchmarkResult{
				Provider:  providerStr,
				Available: false,
				Message:   fmt.Sprintf("Invalid provider: %v", err),
			})
			continue
		}

		result := h.benchmarkSingleProvider(provider, messageCount)
		results = append(results, result)
	}

	// Generate summary
	availableCount := 0
	for _, result := range results {
		if result.Available {
			availableCount++
		}
	}

	summary := fmt.Sprintf("Benchmarked %d providers, %d available", len(results), availableCount)

	response := dto.BenchmarkBBSProvidersResponse{
		Results: results,
		Summary: summary,
	}

	writeSuccessResponse(w, response)
}

func (h *BBSHandler) testSingleProvider(provider bbs.Provider) dto.TestBBSProviderResponse {
	config := &bbs.Config{
		Provider: provider,
	}

	service, err := h.factory.CreateService(provider, config)
	if err != nil {
		return dto.TestBBSProviderResponse{
			Provider:  provider.String(),
			Available: false,
			Message:   fmt.Sprintf("Failed to create service: %v", err),
		}
	}

	// Test basic operations
	start := time.Now()

	// Test key generation
	_, err = service.GenerateKeyPair()
	if err != nil {
		return dto.TestBBSProviderResponse{
			Provider:  provider.String(),
			Available: false,
			Message:   fmt.Sprintf("Key generation failed: %v", err),
		}
	}

	elapsed := time.Since(start)

	return dto.TestBBSProviderResponse{
		Provider:    provider.String(),
		Available:   true,
		Performance: fmt.Sprintf("Key generation: %.2fms", float64(elapsed.Nanoseconds())/1e6),
		Message:     "Provider is working correctly",
	}
}

func (h *BBSHandler) benchmarkSingleProvider(provider bbs.Provider, messageCount int) dto.BenchmarkResult {
	config := &bbs.Config{
		Provider: provider,
	}

	service, err := h.factory.CreateService(provider, config)
	if err != nil {
		return dto.BenchmarkResult{
			Provider:  provider.String(),
			Available: false,
			Message:   fmt.Sprintf("Failed to create service: %v", err),
		}
	}

	// Generate test data
	keyPair, err := service.GenerateKeyPair()
	if err != nil {
		return dto.BenchmarkResult{
			Provider:  provider.String(),
			Available: false,
			Message:   fmt.Sprintf("Key generation failed: %v", err),
		}
	}

	// Generate test messages
	messages := make([][]byte, messageCount)
	for i := 0; i < messageCount; i++ {
		messages[i] = []byte(fmt.Sprintf("test message %d for BBS+ benchmarking", i))
	}

	// Benchmark signing
	start := time.Now()
	signature, err := service.Sign(keyPair.PrivateKey, messages)
	if err != nil {
		return dto.BenchmarkResult{
			Provider:  provider.String(),
			Available: false,
			Message:   fmt.Sprintf("Signing failed: %v", err),
		}
	}
	signTime := float64(time.Since(start).Nanoseconds()) / 1e6

	// Benchmark verification
	start = time.Now()
	err = service.Verify(keyPair.PublicKey, signature, messages)
	if err != nil {
		return dto.BenchmarkResult{
			Provider:  provider.String(),
			Available: false,
			Message:   fmt.Sprintf("Verification failed: %v", err),
		}
	}
	verifyTime := float64(time.Since(start).Nanoseconds()) / 1e6

	// Benchmark proof creation (reveal first half of messages)
	revealedIndices := make([]int, messageCount/2)
	for i := 0; i < messageCount/2; i++ {
		revealedIndices[i] = i
	}

	nonce := []byte("test-nonce-for-benchmarking")
	start = time.Now()
	proof, err := service.CreateProof(signature, keyPair.PublicKey, messages, revealedIndices, nonce)
	if err != nil {
		return dto.BenchmarkResult{
			Provider:  provider.String(),
			Available: false,
			Message:   fmt.Sprintf("Proof creation failed: %v", err),
		}
	}
	proofCreateTime := float64(time.Since(start).Nanoseconds()) / 1e6

	// Benchmark proof verification
	revealedMessages := make([][]byte, len(revealedIndices))
	for i, idx := range revealedIndices {
		revealedMessages[i] = messages[idx]
	}

	start = time.Now()
	err = service.VerifyProof(keyPair.PublicKey, proof, revealedMessages, nonce)
	if err != nil {
		return dto.BenchmarkResult{
			Provider:  provider.String(),
			Available: false,
			Message:   fmt.Sprintf("Proof verification failed: %v", err),
		}
	}
	proofVerifyTime := float64(time.Since(start).Nanoseconds()) / 1e6

	return dto.BenchmarkResult{
		Provider:        provider.String(),
		Available:       true,
		SignTime:        signTime,
		VerifyTime:      verifyTime,
		ProofCreateTime: proofCreateTime,
		ProofVerifyTime: proofVerifyTime,
		Message:         fmt.Sprintf("Successfully benchmarked with %d messages", messageCount),
	}
}

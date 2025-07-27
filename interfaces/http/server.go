package http

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/lugondev/bbs-selective-disclosure-example/interfaces/http/handlers"
	"github.com/lugondev/bbs-selective-disclosure-example/internal/holder"
	"github.com/lugondev/bbs-selective-disclosure-example/internal/issuer"
	"github.com/lugondev/bbs-selective-disclosure-example/internal/verifier"
	"github.com/lugondev/bbs-selective-disclosure-example/pkg/bbs"
)

// Server represents the HTTP server
type Server struct {
	issuerHandler   *handlers.IssuerHandler
	holderHandler   *handlers.HolderHandler
	verifierHandler *handlers.VerifierHandler
	healthHandler   *handlers.HealthHandler
	bbsHandler      *handlers.BBSHandler
	port            string
}

// NewServer creates a new HTTP server
func NewServer(
	issuerUC *issuer.UseCase,
	holderUC *holder.UseCase,
	verifierUC *verifier.UseCase,
	bbsFactory bbs.BBSServiceFactory,
	port string,
) *Server {
	return &Server{
		issuerHandler:   handlers.NewIssuerHandler(issuerUC),
		holderHandler:   handlers.NewHolderHandler(holderUC),
		verifierHandler: handlers.NewVerifierHandler(verifierUC),
		healthHandler:   handlers.NewHealthHandler(),
		bbsHandler:      handlers.NewBBSHandler(bbsFactory),
		port:            port,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Health endpoint
	mux.HandleFunc("/health", s.healthHandler.Health)

	// Issuer endpoints
	mux.HandleFunc("/api/issuer/setup", s.issuerHandler.SetupIssuer)
	mux.HandleFunc("/api/issuer/credentials", s.issuerHandler.IssueCredential)
	mux.HandleFunc("/api/issuer/verify", s.issuerHandler.VerifyCredential)

	// Holder endpoints
	mux.HandleFunc("/api/holder/setup", s.holderHandler.SetupHolder)
	mux.HandleFunc("/api/holder/credentials", s.holderHandler.StoreCredential)
	mux.HandleFunc("/api/holder/credentials/list", s.holderHandler.ListCredentials)
	mux.HandleFunc("/api/holder/presentations", s.holderHandler.CreatePresentation)

	// Verifier endpoints
	mux.HandleFunc("/api/verifier/setup", s.verifierHandler.SetupVerifier)
	mux.HandleFunc("/api/verifier/verify", s.verifierHandler.VerifyPresentation)
	mux.HandleFunc("/api/verifier/verification-request", s.verifierHandler.CreateVerificationRequest)
	mux.HandleFunc("/api/verifier/presentations", s.verifierHandler.ListPresentations)

	// BBS endpoints
	mux.HandleFunc("/api/bbs/test", s.bbsHandler.TestProvider)
	mux.HandleFunc("/api/bbs/benchmark", s.bbsHandler.BenchmarkProviders)

	// Serve static files (for the web UI)
	webDir := "./web/"
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	// Add logging middleware
	loggedMux := loggingMiddleware(mux)

	addr := ":" + s.port
	log.Printf("🚀 BBS+ Selective Disclosure API Server starting on http://localhost%s", addr)
	log.Printf("📱 Web UI available at: http://localhost%s", addr)
	log.Printf("🏥 Health check: http://localhost%s/health", addr)
	log.Printf("📖 API Documentation:")
	log.Printf("   Issuer API: http://localhost%s/api/issuer/*", addr)
	log.Printf("   Holder API: http://localhost%s/api/holder/*", addr)
	log.Printf("   Verifier API: http://localhost%s/api/verifier/*", addr)

	return http.ListenAndServe(addr, loggedMux)
}

// loggingMiddleware logs all incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// ServeStaticFile serves a static file from the web directory
func ServeStaticFile(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filepath := filepath.Join("./web", filename)
		http.ServeFile(w, r, filepath)
	}
}

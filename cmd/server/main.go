package main

import (
	"flag"
	"log"
	"os"

	httpServer "github.com/lugon/bbs-selective-disclosure-example/interfaces/http"
	"github.com/lugon/bbs-selective-disclosure-example/internal/holder"
	"github.com/lugon/bbs-selective-disclosure-example/internal/issuer"
	"github.com/lugon/bbs-selective-disclosure-example/internal/verifier"
	"github.com/lugon/bbs-selective-disclosure-example/pkg/bbs"
	"github.com/lugon/bbs-selective-disclosure-example/pkg/did"
	"github.com/lugon/bbs-selective-disclosure-example/pkg/vc"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "8089", "Server port")
	flag.Parse()

	log.Println("🔐 Initializing BBS+ Selective Disclosure API Server")

	// Initialize services (same as in demo)
	didRepo := did.NewInMemoryRepository()
	didService := did.NewService(didRepo)
	bbsService := bbs.NewService()
	credRepo := vc.NewInMemoryCredentialRepository()
	presRepo := vc.NewInMemoryPresentationRepository()
	vcService := vc.NewService(bbsService, credRepo, presRepo)

	// Initialize use cases
	issuerUC := issuer.NewUseCase(didService, vcService, bbsService)
	holderUC := holder.NewUseCase(didService, vcService, credRepo)
	verifierUC := verifier.NewUseCase(didService, vcService, presRepo)

	// Create and start HTTP server
	server := httpServer.NewServer(issuerUC, holderUC, verifierUC, *port)

	log.Printf("✅ All services initialized successfully")

	// Start server
	if err := server.Start(); err != nil {
		log.Printf("❌ Server failed to start: %v", err)
		os.Exit(1)
	}
}

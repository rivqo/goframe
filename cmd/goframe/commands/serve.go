package commands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/example/goframe/config"
	"github.com/example/goframe/routes"
)

// Serve starts the HTTP server
func Serve(cfg *config.Config) {
	fmt.Printf("Starting server on %s:%d...\n", cfg.Server.Host, cfg.Server.Port)

	// Initialize router
	router, err := routes.InitializeRouter(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize router: %v", err)
	}

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(http.ListenAndServe(addr, router))
}


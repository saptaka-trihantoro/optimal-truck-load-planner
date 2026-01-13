package main

import (
	"log"
	"net/http"

	"github.com/saptaka-trihantoro/optimal-truck-load-planner/internal/api"
	"github.com/saptaka-trihantoro/optimal-truck-load-planner/internal/optimizer"
)

func main() {
	// Initialize logic
	solver := &optimizer.BacktrackingOptimizer{}
	handler := &api.Handler{Optimizer: solver}

	// 1. Requirement: Core Optimization Endpoint
	http.HandleFunc("/api/v1/load-optimizer/optimize", handler.OptimizeHandler)

	// 2. Requirement: Health Check Endpoints
	http.HandleFunc("/healthz", handler.HealthHandler)
	http.HandleFunc("/actuator/health", handler.HealthHandler)

	log.Println("SmartLoad Service starting on :8080...")
	log.Println("Endpoint: POST /api/v1/load-optimizer/optimize")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}

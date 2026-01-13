package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saptaka-trihantoro/optimal-truck-load-planner/internal/api"
	"github.com/saptaka-trihantoro/optimal-truck-load-planner/internal/domain"
	"github.com/saptaka-trihantoro/optimal-truck-load-planner/internal/optimizer"
)

func TestOptimizeEndpoint_E2E(t *testing.T) {
	// 1. Setup
	solver := &optimizer.BacktrackingOptimizer{}
	handler := &api.Handler{Optimizer: solver}

	// Create a test server mux to mimic the routing in main.go
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/load-optimizer/optimize", handler.OptimizeHandler)

	// 2. Mock Payload (Similar to sample-request.json)
	requestPayload := domain.LoadRequest{
		Truck: domain.Truck{
			ID: "truck-123", MaxWeightLbs: 44000, MaxVolumeCuft: 3000,
		},
		Orders: []domain.Order{
			{ID: "ord-001", WeightLbs: 10000, VolumeCuft: 500, PayoutCents: 100000, IsHazmat: false},
			{ID: "ord-002", WeightLbs: 10000, VolumeCuft: 500, PayoutCents: 150000, IsHazmat: false},
		},
	}

	body, _ := json.Marshal(requestPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/load-optimizer/optimize", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 3. Execution
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// 4. Assertions
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response domain.LoadResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.TotalPayoutCents != 250000 {
		t.Errorf("Expected payout 250000, got %d", response.TotalPayoutCents)
	}

	if len(response.SelectedOrderIDs) != 2 {
		t.Errorf("Expected 2 selected orders, got %d", len(response.SelectedOrderIDs))
	}
}

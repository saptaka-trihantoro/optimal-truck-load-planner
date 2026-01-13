package api

import (
	"encoding/json"
	"net/http"

	"github.com/saptaka-trihantoro/optimal-truck-load-planner/internal/domain"
)

type Handler struct {
	Optimizer domain.Optimizer
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// A standard response format for health checks
	response := map[string]string{
		"status": "UP",
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) OptimizeHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Validate Method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 2. Validate Payload Size (413 if too large)
	// We assume > 2MB is too large for this task
	if r.ContentLength > 2*1024*1024 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	var req domain.LoadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400 on invalid JSON
		return
	}

	// 3. Logic Validation (400 on invalid input)
	if len(req.Orders) == 0 || req.Truck.MaxWeightLbs <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 4. Solve
	result := h.Optimizer.FindBestLoad(req.Truck, req.Orders)

	// 5. Success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

package optimizer

import (
	"fmt"
	"testing"
	"time"

	"github.com/saptaka-trihantoro/optimal-truck-load-planner/internal/domain"
)

func TestBacktrackingOptimizer_Correctness(t *testing.T) {
	truck := domain.Truck{ID: "T1", MaxWeightLbs: 100, MaxVolumeCuft: 100}

	orders := []domain.Order{
		{ID: "O1", WeightLbs: 50, VolumeCuft: 50, PayoutCents: 1000, IsHazmat: false},
		{ID: "O2", WeightLbs: 60, VolumeCuft: 60, PayoutCents: 2000, IsHazmat: false}, // Fits alone
		{ID: "O3", WeightLbs: 40, VolumeCuft: 40, PayoutCents: 500, IsHazmat: true},   // Hazmat
	}

	optimizer := &BacktrackingOptimizer{}
	res := optimizer.FindBestLoad(truck, orders)

	// Test 1: Payout (Should pick O2 because 2000 > 1000+500)
	if res.TotalPayoutCents != 2000 {
		t.Errorf("Expected payout 2000, got %d", res.TotalPayoutCents)
	}

	// Test 2: Hazmat Isolation
	// If we add another Hazmat order that makes the combo better than O2...
	ordersWithBetterHazmat := append(orders, domain.Order{
		ID: "O4", WeightLbs: 10, VolumeCuft: 10, PayoutCents: 2000, IsHazmat: true,
	})
	resHazmat := optimizer.FindBestLoad(truck, ordersWithBetterHazmat)
	// Should pick O3 and O4 (payout 2500) and NOT mix with O1 or O2
	for _, id := range resHazmat.SelectedOrderIDs {
		if id == "O1" || id == "O2" {
			t.Errorf("Hazmat isolation failed: mixed non-hazmat %s with hazmat", id)
		}
	}
}

func TestPerformance_N22(t *testing.T) {
	truck := domain.Truck{ID: "T-PERF", MaxWeightLbs: 44000, MaxVolumeCuft: 3000}
	var orders []domain.Order
	for i := 0; i < 22; i++ {
		orders = append(orders, domain.Order{
			ID:          fmt.Sprintf("ORD-%d", i),
			WeightLbs:   2000,
			VolumeCuft:  150,
			PayoutCents: int64(1000 + i),
			IsHazmat:    false,
		})
	}

	optimizer := &BacktrackingOptimizer{}
	start := time.Now()
	_ = optimizer.FindBestLoad(truck, orders)
	elapsed := time.Since(start)

	t.Logf("Time for N=22: %s", elapsed)
	if elapsed > 800*time.Millisecond {
		t.Errorf("Performance failed: took %s, limit 800ms", elapsed)
	}
}

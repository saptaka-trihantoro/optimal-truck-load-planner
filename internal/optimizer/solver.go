package optimizer

import (
	"math"

	"sort"

	"github.com/saptaka-trihantoro/optimal-truck-load-planner/internal/domain"
)

type BacktrackingOptimizer struct{}

func (o *BacktrackingOptimizer) FindBestLoad(truck domain.Truck, orders []domain.Order) domain.LoadResponse {
	// 1. Filter: Valid time windows (Pickup <= Delivery)
	validOrders := make([]domain.Order, 0)
	for _, ord := range orders {
		if ord.PickupDate <= ord.DeliveryDate {
			validOrders = append(validOrders, ord)
		}
	}

	// 2. Sort by Value Density (Heuristic to find high payout paths first for pruning)
	sort.Slice(validOrders, func(i, j int) bool {
		valI := float64(validOrders[i].PayoutCents) / float64(validOrders[i].WeightLbs+validOrders[i].VolumeCuft)
		valJ := float64(validOrders[j].PayoutCents) / float64(validOrders[j].WeightLbs+validOrders[j].VolumeCuft)
		return valI > valJ
	})

	var bestPayout int64
	var bestSet []domain.Order

	var backtrack func(idx int, curW, curV int, curP int64, selected []domain.Order)
	backtrack = func(idx int, curW, curV int, curP int64, selected []domain.Order) {
		if curP > bestPayout {
			bestPayout = curP
			bestSet = make([]domain.Order, len(selected))
			copy(bestSet, selected)
		}

		for i := idx; i < len(validOrders); i++ {
			ord := validOrders[i]

			// Constraint: Weight & Volume
			if curW+ord.WeightLbs > truck.MaxWeightLbs || curV+ord.VolumeCuft > truck.MaxVolumeCuft {
				continue
			}

			// Constraint: Hazmat Isolation
			// Rule: All selected orders must have the same Hazmat status
			if len(selected) > 0 && selected[0].IsHazmat != ord.IsHazmat {
				continue
			}

			selected = append(selected, ord)
			backtrack(i+1, curW+ord.WeightLbs, curV+ord.VolumeCuft, curP+ord.PayoutCents, selected)
			selected = selected[:len(selected)-1]
		}
	}

	backtrack(0, 0, 0, 0, []domain.Order{})
	return buildResponse(truck, bestSet, bestPayout)
}

func buildResponse(t domain.Truck, selected []domain.Order, payout int64) domain.LoadResponse {
	var w, v int
	ids := []string{}
	for _, o := range selected {
		ids = append(ids, o.ID)
		w += o.WeightLbs
		v += o.VolumeCuft
	}

	weightUtil := 0.0
	if t.MaxWeightLbs > 0 {
		weightUtil = math.Round((float64(w)/float64(t.MaxWeightLbs)*100)*100) / 100
	}

	volumeUtil := 0.0
	if t.MaxVolumeCuft > 0 {
		volumeUtil = math.Round((float64(v)/float64(t.MaxVolumeCuft)*100)*100) / 100
	}

	return domain.LoadResponse{
		TruckID:                  t.ID,
		SelectedOrderIDs:         ids,
		TotalPayoutCents:         payout,
		TotalWeightLbs:           w,
		TotalVolumeCuft:          v,
		UtilizationWeightPercent: weightUtil,
		UtilizationVolumePercent: volumeUtil,
	}
}

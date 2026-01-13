package domain

type Order struct {
	ID           string `json:"id"`
	WeightLbs    int    `json:"weight_lbs"`
	VolumeCuft   int    `json:"volume_cuft"`
	PayoutCents  int64  `json:"payout_cents"`
	IsHazmat     bool   `json:"is_hazmat"`
	PickupDate   int64  `json:"pickup_date"`
	DeliveryDate int64  `json:"delivery_date"`
}

type Truck struct {
	ID            string `json:"id"`
	MaxWeightLbs  int    `json:"max_weight_lbs"`
	MaxVolumeCuft int    `json:"max_volume_cuft"`
}

type LoadRequest struct {
	Truck  Truck   `json:"truck"`
	Orders []Order `json:"orders"`
}

type LoadResponse struct {
	TruckID                  string   `json:"truck_id"`
	SelectedOrderIDs         []string `json:"selected_order_ids"`
	TotalPayoutCents         int64    `json:"total_payout_cents"`
	TotalWeightLbs           int      `json:"total_weight_lbs"`
	TotalVolumeCuft          int      `json:"total_volume_cuft"`
	UtilizationWeightPercent float64  `json:"utilization_weight_percent"`
	UtilizationVolumePercent float64  `json:"utilization_volume_percent"`
}
type Optimizer interface {
	FindBestLoad(truck Truck, orders []Order) LoadResponse
}

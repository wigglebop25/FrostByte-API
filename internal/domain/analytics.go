package domain

type RevenueAnalyticsResponse struct {
	Period    Period      `json:"period"`
	Summary   Summary     `json:"summary"`
	DailyData []DailyData `json:"daily_data"`
}

type Period struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Summary struct {
	TotalRevenue float64 `json:"total_revenue"`
	TotalOrders  int64   `json:"total_orders"`
}

type DailyData struct {
	Date       string  `json:"date"`
	DayName    string  `json:"day_name"`
	Revenue    float64 `json:"revenue"`
	OrderCount int64   `json:"order_count"`
}

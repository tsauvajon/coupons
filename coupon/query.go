package coupon

import "time"

// '{ id: 1, brand: "Tesco", value: 34, name: "Save £34 at Tesco", expiryUtc: "2019-09-16T11:45:26.371Z" }'

// FilterQuery is a query to list coupons, all fields will be applied as filters
type FilterQuery struct {
	// 1 <= Take <= 20
	Take, Skip int64
	BrandName  string
	// Min <= Max, 0 = disabled
	MinValue, MaxValue float64
	ExpiresAfter       time.Time
}

// SaveQuery is a query to update an existing coupon or save a new one
type SaveQuery struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Brand     string    `json:"brand"`
	Value     float64   `json:"value"`
	ExpiryUTC time.Time `json:"expiryUtc"`
}

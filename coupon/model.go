package coupon

import "time"

/*
Coupon is a discount available in a shop

Example coupon:
{
	“name”: “Save £20 at Tesco”,
	“brand”: "Tesco",
	“value”: 20,
	“createdAt”: “2018-03-01 10:15:53”,
	“expiry”: “2019-03-01 10:15:53”
}
*/
type Coupon struct {
	ID           int64
	Name         string
	Brand        *Brand
	Value        float64
	CreatedAtUTC time.Time
	ExpiryUTC    time.Time
}

// Brand is a shop where you can use coupons
type Brand struct {
	ID   int64 `json:"-"`
	Name string
}

// Query is a query to list coupons, all fields will be applied as filters
type Query struct {
	// 1 <= Take <= 20
	Take, Skip int64
	BrandName  string
	// Min <= Max, 0 = disabled
	MinValue, MaxValue float64
	ExpiresAfter       time.Time
}

package coupon

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tsauvajon/coupons/database"
)

// MaxValueMoney : Maximum value for the money type in Postgres
const MaxValueMoney = float64(92233720368547658.0)

// Client is a Coupons database client
type Client struct {
	db *database.Client
}

// NewClient connects to the db and returns a coupons client
func NewClient() (*Client, error) {
	dbcli, err := database.NewClient()

	return &Client{db: dbcli}, err
}

// ListCoupons will return all coupons that match a filter
func (c *Client) ListCoupons(ctx context.Context, q *Query) ([]*Coupon, error) {
	if q.Take < 0 {
		return nil, newErrBadRequest("take at least 1")
	}
	if q.Take > 20 {
		return nil, newErrBadRequest("take 20 or less")
	}
	if q.Take == 0 {
		q.Take = 20
	}
	if q.Skip < 0 {
		return nil, newErrBadRequest("skip should be positive")
	}
	if q.MaxValue == 0 {
		q.MaxValue = MaxValueMoney
	}
	if q.MaxValue > MaxValueMoney {
		return nil, newErrBadRequest("max value can't be that high")
	}
	if q.MinValue < 0 {
		return nil, newErrBadRequest("min should be positive")
	}
	if q.MaxValue < 0 {
		return nil, newErrBadRequest("max should be positive")
	}
	if q.MinValue > q.MaxValue {
		return nil, newErrBadRequest("min should be less than max")
	}

	filters, _ := json.Marshal(q)
	log.Println("filters: ", string(filters))

	rows, err := c.db.Connection.Query(
		`
SELECT
	coupons.id,
	TRIM(coupons.name),
	value::money::numeric::float8,
	created_at_utc,
	expiry_utc,
	TRIM(brands.name)
FROM coupons
JOIN brands ON brands.id = coupons.brand_id
WHERE
 	value >= $1 AND
 	value <= $2 AND
 	lower(brands.name) like '%' || lower($3) || '%' AND
	expiry_utc >= $4
ORDER BY value desc
LIMIT $5 OFFSET $6
`, q.MinValue, q.MaxValue, q.BrandName, q.ExpiresAfter, q.Take, q.Skip)

	if err != nil {
		log.Println("sql error:", err)
		return nil, newErrInternal(err.Error())
	}

	// Will execute at the end of the scope
	defer rows.Close()

	coupons := []*Coupon{}

	i := 0
	for rows.Next() {
		c := &Coupon{Brand: &Brand{}}

		if err := rows.Scan(&c.ID, &c.Name, &c.Value, &c.CreatedAtUTC, &c.ExpiryUTC, &c.Brand.Name); err != nil {
			return nil, err
		}
		fmt.Println(i, ":", c.Name)
		i++

		coupons = append(coupons, c)
	}

	return coupons, nil
}

package coupon

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

// ListCoupons will return all coupons that match a filter.
// Handles validation
func (c *Client) ListCoupons(ctx context.Context, q *FilterQuery) ([]*Coupon, error) {
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
	if q.MinValue > MaxValueMoney {
		return nil, newErrBadRequest("min value can't be that high")
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

	rows, err := c.db.Connection.Query(`
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
		ORDER BY id desc
		LIMIT $5 OFFSET $6
		`, q.MinValue, q.MaxValue, q.BrandName, q.ExpiresAfter, q.Take, q.Skip)

	if err != nil {
		log.Println("sql error:", err)
		return nil, newErrInternal(err.Error())
	}

	defer rows.Close()

	coupons := []*Coupon{}

	for rows.Next() {
		c := &Coupon{Brand: &Brand{}}

		if err := rows.Scan(&c.ID, &c.Name, &c.Value, &c.CreatedAtUTC, &c.ExpiryUTC, &c.Brand.Name); err != nil {
			return nil, err
		}
		fmt.Printf("%d: %s\n", c.ID, c.Name)

		coupons = append(coupons, c)
	}

	return coupons, nil
}

// GetBrandByName gets a brand from its name, and create it if it doesn't exist
func (c *Client) GetBrandByName(ctx context.Context, brand string) (*Brand, error) {
	b := &Brand{Name: brand}
	err := c.db.Connection.QueryRow(
		`SELECT id FROM brands WHERE lower(name) = lower($1)`,
		brand).Scan(&b.ID)

	// Found the brand, return it
	if err == nil {
		return b, nil
	}

	// Didn't find it, create it
	err = c.db.Connection.
		QueryRow(
			`INSERT INTO brands(name) VALUES ($1) returning id`,
			brand,
		).
		Scan(&b.ID)

	return b, err
}

// CreateCoupon will add a new coupon
func (c *Client) CreateCoupon(ctx context.Context, q *SaveQuery) (*Coupon, error) {
	if err := validatePayload(q); err != nil {
		return nil, err
	}

	b, err := c.GetBrandByName(ctx, q.Brand)

	if err != nil {
		return nil, newErrInternal("couldn't get the brand: " + err.Error())
	}

	coupon := &Coupon{
		Name:      q.Name,
		Value:     q.Value,
		ExpiryUTC: q.ExpiryUTC,
		Brand:     &Brand{Name: q.Brand},
	}

	err = c.db.Connection.
		QueryRow(
			`INSERT INTO coupons(name, value, expiry_utc, brand_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at_utc`,
			coupon.Name, coupon.Value, coupon.ExpiryUTC, b.ID,
		).
		Scan(&coupon.ID, &coupon.CreatedAtUTC)

	if err != nil {
		err = newErrInternal("couldn't insert the coupon: " + err.Error())
		log.Println(err)
		return nil, err
	}

	return coupon, nil
}

// UpdateCoupon will add a new coupon
func (c *Client) UpdateCoupon(ctx context.Context, q *SaveQuery) (*Coupon, error) {
	if err := validatePayload(q); err != nil {
		return nil, err
	}

	b, err := c.GetBrandByName(ctx, q.Brand)

	if err != nil {
		return nil, newErrInternal("couldn't get the brand: " + err.Error())
	}

	coupon := &Coupon{
		ID:        q.ID,
		Name:      q.Name,
		Value:     q.Value,
		ExpiryUTC: q.ExpiryUTC,
		Brand:     &Brand{Name: q.Brand},
	}

	_, err = c.db.Connection.
		Exec(
			`UPDATE coupons
			SET
				name = $1,
				value = $2,
				expiry_utc = $3,
				brand_id = $4
			WHERE ID = $5`,
			coupon.Name, coupon.Value, coupon.ExpiryUTC, b.ID, coupon.ID,
		)

	if err != nil {
		err = newErrInternal("couldn't update the coupon: " + err.Error())
		log.Println(err)
		return nil, err
	}

	return coupon, nil
}

func validatePayload(q *SaveQuery) error {
	if len(q.Brand) < 3 {
		return newErrBadRequest("brand must be at least 3 characters long")
	}

	if len(q.Name) < 3 {
		return newErrBadRequest("name must be at least 3 characters long")
	}

	if q.ExpiryUTC.Before(time.Now().UTC()) {
		return newErrBadRequest("the coupon must expire in the future")
	}

	if q.Value < 0.01 {
		return newErrBadRequest("discount must be at least 1 penny...")
	}

	return nil
}

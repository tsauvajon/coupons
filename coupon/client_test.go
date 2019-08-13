package coupon

import (
	"context"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client, err := NewClient()

	if err != nil {
		t.Errorf("creating client failed: %s", err)
		return
	}

	exp := time.Now().AddDate(1, 0, 0)
	q := &SaveQuery{
		Brand:     "co-op",
		ExpiryUTC: exp,
		Name:      "save some money",
		Value:     12.3,
	}
	coupon, err := client.CreateCoupon(context.Background(), q)

	if err != nil {
		t.Errorf("creating the coupon failed: %s", err)
		return
	}

	if coupon.Brand.Name != q.Brand {
		t.Errorf("brand doesn't match: %s | %s", coupon.Brand.Name, q.Brand)
		return
	}

	if coupon.Name != q.Name {
		t.Errorf("name doesn't match: %s | %s", coupon.Name, q.Name)
		return
	}

	if coupon.Value != q.Value {
		t.Errorf("value doesn't match: %v | %v", coupon.Value, q.Value)
		return
	}

	if coupon.ExpiryUTC != q.ExpiryUTC {
		t.Errorf("expiry date doesn't match: %s | %s", coupon.ExpiryUTC, q.ExpiryUTC)
		return
	}

	if coupon.ID <= 0 {
		t.Errorf("ID doesn't look right: %d", coupon.ID)
		return
	}

	if coupon.CreatedAtUTC.Before(time.Now().AddDate(0, -1, 0)) {
		t.Errorf("the creation date doesn't look right: %s", coupon.CreatedAtUTC)
		return
	}

	q.ID = coupon.ID
	q.Value = 123.45
	q.Brand = "tesco"

	coupon, err = client.UpdateCoupon(context.Background(), q)

	if err != nil {
		t.Errorf("updating the coupon failed: %s", err)
		return
	}

	if coupon.Brand.Name != q.Brand {
		t.Errorf("brand doesn't match: %s | %s", coupon.Brand.Name, q.Brand)
		return
	}

	if coupon.Name != q.Name {
		t.Errorf("name doesn't match: %s | %s", coupon.Name, q.Name)
		return
	}

	if coupon.Value != q.Value {
		t.Errorf("value doesn't match: %v | %v", coupon.Value, q.Value)
		return
	}

	coupons, err := client.ListCoupons(context.Background(), &FilterQuery{})

	if err != nil {
		t.Errorf("getting the coupons failed: %s", err)
		return
	}

	if len(coupons) != 1 {
		t.Errorf("there should be exactly 1 coupon, found %d", len(coupons))
		return
	}

	cp := coupons[0]

	if coupon.Brand.Name != cp.Brand.Name {
		t.Errorf("brand doesn't match: %s | %s", coupon.Brand.Name, cp.Brand.Name)
		return
	}

	if coupon.Name != cp.Name {
		t.Errorf("name doesn't match: %s | %s", coupon.Name, cp.Name)
		return
	}

	if coupon.Value != cp.Value {
		t.Errorf("value doesn't match: %v | %v", coupon.Value, cp.Value)
		return
	}

	// Special case, the date gets rounded to 0.00001 second, so only comparing
	// up to the second... it should be enough for this use case
	expected, actual := coupon.ExpiryUTC.Format("01-02-2006 15:04:05"), cp.ExpiryUTC.Format("01-02-2006 15:04:05")
	if expected != actual {
		t.Errorf("expiry date doesn't match: %s | %s", expected, actual)
		return
	}

	if coupon.ID <= 0 {
		t.Errorf("ID doesn't look right: %d", coupon.ID)
		return
	}

	if cp.CreatedAtUTC.Before(time.Now().AddDate(0, -1, 0)) {
		t.Errorf("the creation date doesn't look right: %s", cp.CreatedAtUTC)
		return
	}
}

func TestValidation(t *testing.T) {
	exp := time.Now().AddDate(1, 0, 0)
	q := SaveQuery{
		Brand:     "co-op",
		ExpiryUTC: exp,
		Name:      "save some money",
		Value:     12.3,
	}

	if err := validatePayload(&q); err != nil {
		t.Errorf("unvalidated good payload: %s", err)
	}

	brand := q
	brand.Brand = ""

	if err := validatePayload(&brand); err == nil {
		t.Error("should have invalided empty brand name")
	}

	date := q
	date.ExpiryUTC = date.ExpiryUTC.AddDate(-1, 0, 0)

	if err := validatePayload(&date); err == nil {
		t.Error("should have invalided empty brand name")
	}

	name := q
	name.Name = ""

	if err := validatePayload(&name); err == nil {
		t.Error("should have invalided empty name")
	}

	value := q
	value.Value = -1

	if err := validatePayload(&value); err == nil {
		t.Error("should have invalided negative value")
	}
}

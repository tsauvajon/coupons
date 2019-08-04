package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	s, err := newServer()

	if err != nil {
		t.Errorf("couldn't create server: %s", err)
		return
	}

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/coupons", nil)
	if err != nil {
		t.Errorf("couldn't create request: %s", err)
		return
	}

	s.GetCoupons(rr, req)

	if st := rr.Code; st != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, st)
		return
	}

	if rr.Body.String() != "[]" {
		t.Errorf("expected an empty array, got %s", rr.Body.String())
	}

	// No need to test the behavior as it is the clients responsibility, and
	// it is tested in the client (coupon) package
}

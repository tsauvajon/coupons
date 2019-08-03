package main

import "net/http"

type server struct {
}

func (s *server) CouponsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.GetCoupons(w, r)
	case http.MethodPost:
		s.CreateCoupon(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *server) GetCoupons(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get coupons"))
	w.WriteHeader(200)
}

func (s *server) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create coupons"))
	w.WriteHeader(200)
}

func (s *server) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update coupon"))
	w.WriteHeader(200)
}

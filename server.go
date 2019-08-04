package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tsauvajon/coupons/coupon"
)

type server struct {
	client *coupon.Client
}

func newServer() (*server, error) {
	c, err := coupon.NewClient()

	return &server{client: c}, err
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		// won't cause an infinite loop because json.Marshal can't fail
		respondWithError(w, code, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respond(w http.ResponseWriter, payload interface{}, err error) {
	if err != nil {
		switch err.(type) {
		case *coupon.ErrBadRequest:
			respondWithError(w, http.StatusBadRequest, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (s *server) CouponsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.GetCoupons(w, r)
	case http.MethodPost:
		s.CreateCoupon(w, r)
	case http.MethodPut:
		s.UpdateCoupon(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// Parses GET queries, formats responses
func (s *server) GetCoupons(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	q := &coupon.FilterQuery{
		BrandName: values.Get("brand"),
	}

	// Parse min as float if provided, else keep default value of 0
	if len(values["min"]) != 0 {
		min, err := strconv.ParseFloat(values.Get("min"), 64)
		if err != nil {
			err = fmt.Errorf("couldn't parse min: %s", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		q.MinValue = min
	}

	if len(values["max"]) != 0 {
		max, err := strconv.ParseFloat(values.Get("max"), 64)
		if err != nil {
			err = fmt.Errorf("couldn't parse max: %s", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		q.MaxValue = max
	}

	if len(values["skip"]) != 0 {
		skip, err := strconv.ParseInt(values.Get("skip"), 10, 64)
		if err != nil {
			err = fmt.Errorf("couldn't parse skip: %s", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		q.Skip = skip
	}

	if len(values["take"]) != 0 {
		take, err := strconv.ParseInt(values.Get("take"), 10, 64)
		if err != nil {
			err = fmt.Errorf("couldn't parse take: %s", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		q.Take = take
	}

	// if no date is provided we'll keep the default date 01/01/0001
	if strExp := values.Get("exp"); strExp != "" {
		exp, err := time.Parse(time.RFC3339, strExp)
		if err != nil {
			err = fmt.Errorf("couldn't parse exp date: %s", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		q.ExpiresAfter = exp
	}

	coupons, err := s.client.ListCoupons(context.Background(), q)

	respond(w, coupons, err)
}

func (s *server) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var q *coupon.SaveQuery
	if err := decoder.Decode(&q); err != nil {
		err = fmt.Errorf("couldn't parse the body: %s", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	coupon, err := s.client.CreateCoupon(context.Background(), q)

	respond(w, coupon, err)
}

func (s *server) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var q *coupon.SaveQuery
	if err := decoder.Decode(&q); err != nil {
		err = fmt.Errorf("couldn't parse the body: %s", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	coupon, err := s.client.UpdateCoupon(context.Background(), q)

	respond(w, coupon, err)
}

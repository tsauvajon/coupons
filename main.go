package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	// default value
	if port == "" {
		port = "4000"
	}

	// TODO: Connect to db

	s := &server{}
	mux := http.NewServeMux()

	mux.Handle("/coupons", http.HandlerFunc(s.CouponsHandler))
	mux.Handle("/coupons/{id}", http.HandlerFunc(s.UpdateCoupon))

	log.Println("Starting server at :" + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Error on creating listener: ", err)
	}
}

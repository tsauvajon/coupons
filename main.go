package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	// default value 4000
	if port == "" {
		port = "4000"
	}

	s, err := newServer()
	if err != nil {
		log.Fatal("couldn't connect to the database: ", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/coupons", http.HandlerFunc(s.CouponsHandler))

	log.Println("starting server at :" + port)
	if err = http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("error on creating listener: ", err)
	}
}

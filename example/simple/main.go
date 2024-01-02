package main

import (
	"log"
	"net/http"
)

const (
	ListenDSN = "localhost:80"
)

func main() {
	s, err := NewServer()
	mustBeNoError(err)

	http.HandleFunc("/", s.rootHandler)
	err = http.ListenAndServe(ListenDSN, nil)
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

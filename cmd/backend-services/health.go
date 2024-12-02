package main

import (
	"fmt"
	"net/http"
)

func setupHealthEndpoint() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Service is up and running now!")
	})
}

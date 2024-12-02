package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ashwingopalsamy/backend-services/pkg/handler"
	"github.com/ashwingopalsamy/backend-services/pkg/store/service"
)

func startServer() {
	storeService := &service.StoreServiceServer{}

	storeHandler := handler.NewStoreHandler(storeService)

	http.HandleFunc("/api/v1/stores", storeHandler.CreateStoreHandler)

	port := "8080"
	fmt.Printf("HTTP server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

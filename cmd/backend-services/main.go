package main

import (
	"log"
)

func main() {
	setupHealthEndpoint()
	setupPersistence()
	
	log.Println("Starting backend services...")
	startServer()
}

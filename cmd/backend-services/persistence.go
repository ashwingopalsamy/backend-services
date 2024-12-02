package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func setupPersistence() {
	db, err := sql.Open("postgres", "postgresql://ourzhop:ourzhop@localhost:5432/ourzhop?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

}

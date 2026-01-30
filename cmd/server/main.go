package main

import (
	"base-skeleton/config"
	"base-skeleton/internal/database"
	"base-skeleton/internal/module/router"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	db, err := database.NewSupabase(cfg)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer db.Close()

	handler := router.New(db)

	log.Println("🚀 HTTP server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

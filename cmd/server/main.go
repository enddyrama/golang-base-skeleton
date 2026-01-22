package main

import (
	"base-skeleton/internal/module/router"
	"log"
	"net/http"
)

func main() {
	handler := router.New()

	log.Println("HTTP server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

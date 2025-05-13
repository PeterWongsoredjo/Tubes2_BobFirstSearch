package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request method - POST required",
		})
		return
	}

	err := runScraping()
	if err != nil {
		log.Printf("Scraping failed: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Scraping failed",
			"details": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Scraping completed successfully",
	})
}

func main() {
	http.HandleFunc("/api/scrape", enableCORS(scrapeHandler))

	fmt.Println("Server starting on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
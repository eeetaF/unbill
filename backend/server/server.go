package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"

	"backend/endpoints"
)

// Expose initializes the HTTP server
func Expose() {
	mux := http.NewServeMux()

	// Handle CORS preflight (OPTIONS requests)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}
	})

	mux.HandleFunc("POST /api/upload_and_analyze", endpoints.Uploader)
	mux.HandleFunc("GET /api/split_equally/{filename}/{num_people}", endpoints.EvenSplitter)
	mux.HandleFunc("PUT /api/products/{filename}", endpoints.ProductsPutter)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8081"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := c.Handler(mux)

	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

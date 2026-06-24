package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("debugprint: entering main")
	port := "3456"

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GET /health from %s", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	fmt.Printf("Mock server starting on port %s...\n", port)
	fmt.Printf("Health check available at http://localhost:%s/health\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

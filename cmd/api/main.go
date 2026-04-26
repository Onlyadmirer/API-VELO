package main

import (
	"VELO-backend/internal/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.ConnectDB()

	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Server is running"}`)
	})

	log.Println("server berjalan di port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("gagal terhubung ke server", err)
	}
}

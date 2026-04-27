package main

import (
	"VELO-backend/internal/config"
	"VELO-backend/internal/handler"
	"VELO-backend/internal/repository"
	"VELO-backend/internal/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Aplikasi berhenti: ", err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/products", productHandler.GetAllProducts)
	mux.HandleFunc("POST /api/products", productHandler.CreateProduct)

	fmt.Println("server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

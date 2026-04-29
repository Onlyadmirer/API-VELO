package main

import (
	"VELO-backend/internal/config"
	"VELO-backend/internal/handler"
	"VELO-backend/internal/middleware"
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

	// user
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// product
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	mux := http.NewServeMux()

	// user
	mux.HandleFunc("POST /api/users/register", userHandler.CreateUser)
	mux.HandleFunc("POST /api/users/login", userHandler.UserLogin)

	// product
	mux.HandleFunc("GET /api/products", productHandler.GetAllProducts)
	mux.HandleFunc("POST /api/products", middleware.JWTMiddleware(productHandler.CreateProduct))
	mux.HandleFunc("DELETE /api/products/{id}", middleware.JWTMiddleware(productHandler.DeleteProduct))
	mux.HandleFunc("PUT /api/products/{id}", middleware.JWTMiddleware(productHandler.UpdateProduct))

	fmt.Println("server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

package main

import (
	"VELO-backend/pkg/config"
	cronhandler "VELO-backend/pkg/cron_Handler"
	"VELO-backend/pkg/handler"
	"VELO-backend/pkg/middleware"
	"VELO-backend/pkg/payment"
	"VELO-backend/pkg/repository"
	"VELO-backend/pkg/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/midtrans/midtrans-go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found")
	}

	midtrans.ServerKey = os.Getenv("SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Aplikasi berhenti: ", err)
	}
	defer db.Close()

	redisClient := config.ConnectRedis()
	_ = redisClient

	midtrans := &payment.MidtransClient{}

	emailService := service.NewEmailService()

	// user
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, emailService)
	userHandler := handler.NewUserHandler(userService)

	// cart
	cartRepo := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepo)
	cartHandler := handler.NewCartHandler(cartService)

	// order
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, cartRepo, midtrans)
	orderHandler := handler.NewOrderHandler(orderService)

	// product
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	mux := http.NewServeMux()

	// user
	mux.HandleFunc("POST /api/users/register", userHandler.CreateUser)
	mux.HandleFunc("GET /api/users/verify", userHandler.HandleVerify)
	mux.HandleFunc("POST /api/users/login", userHandler.UserLogin)
	mux.HandleFunc("POST /api/users/logout", middleware.JWTMiddleware(userHandler.LogOut))

	// cart
	mux.HandleFunc("POST /api/cart", middleware.JWTMiddleware(cartHandler.AddToCart))
	mux.HandleFunc("GET /api/cart", middleware.JWTMiddleware(cartHandler.GetCart))

	// order
	mux.HandleFunc("POST /api/checkout", middleware.JWTMiddleware(orderHandler.CheckOut))
	mux.HandleFunc("POST /api/webhook/midtrans", (orderHandler.MidtransNotifications))
	mux.HandleFunc("GET /api/orders", middleware.JWTMiddleware(orderHandler.GetOrder))

	// product
	mux.HandleFunc("GET /api/products", productHandler.GetAllProducts)
	mux.HandleFunc("POST /api/products", middleware.JWTMiddleware(middleware.RBACMiddleware(productHandler.CreateProduct)))
	mux.HandleFunc("DELETE /api/products/{id}", middleware.JWTMiddleware(middleware.RBACMiddleware(productHandler.DeleteProduct)))
	mux.HandleFunc("PUT /api/products/{id}", middleware.JWTMiddleware(middleware.RBACMiddleware(productHandler.UpdateProduct)))

	// cron
	mux.HandleFunc("GET /api/cron", cronhandler.CronSendResponse)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

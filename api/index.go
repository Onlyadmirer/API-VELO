package api

import (
	"VELO-backend/internal/config"
	"VELO-backend/internal/middleware"
	"VELO-backend/internal/repository"
	"VELO-backend/internal/service"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/midtrans/midtrans-go"

	"VELO-backend/internal/handler"
)

var (
	db             *sql.DB
	userHandler    *handler.UserHandler
	cartHandler    *handler.CartHandler
	orderHandler   *handler.OrderHandler
	productHandler *handler.ProductHandler
)

func init() {
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

	// user
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler = handler.NewUserHandler(userService)

	// cart
	cartRepo := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepo)
	cartHandler = handler.NewCartHandler(cartService)

	// order
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, cartRepo)
	orderHandler = handler.NewOrderHandler(orderService)

	// product
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler = handler.NewProductHandler(productService)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	// user
	mux.Handle("POST /api/users/register", http.HandlerFunc(userHandler.CreateUser))
	mux.Handle("POST /api/users/login", http.HandlerFunc(userHandler.UserLogin))

	// cart
	mux.Handle("POST /api/cart", middleware.JWTMiddleware(http.HandlerFunc(cartHandler.AddToCart)))
	mux.Handle("GET /api/cart", middleware.JWTMiddleware(http.HandlerFunc(cartHandler.GetCart)))

	// order
	mux.Handle("POST /api/checkout", middleware.JWTMiddleware(http.HandlerFunc(orderHandler.CheckOut)))
	mux.Handle("POST /api/webhook/midtrans", http.HandlerFunc(orderHandler.MidtransNotifications))

	// product
	mux.Handle("GET /api/products", http.HandlerFunc(productHandler.GetAllProducts))
	mux.Handle("POST /api/products", middleware.JWTMiddleware(http.HandlerFunc(productHandler.CreateProduct)))
	mux.Handle("DELETE /api/products/{id}", middleware.JWTMiddleware(http.HandlerFunc(productHandler.DeleteProduct)))
	mux.Handle("PUT /api/products/{id}", middleware.JWTMiddleware(http.HandlerFunc(productHandler.UpdateProduct)))

	mux.ServeHTTP(w, r)
}

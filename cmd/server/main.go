package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/your-username/shop01/internal/handlers"
	"github.com/your-username/shop01/internal/store"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Connect to database
	db, err := store.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	// Set up router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Routes
	h := handlers.New(db)
	r.Get("/", h.Home)
	r.Get("/products", h.Products)
	r.Get("/products/{id}", h.ProductDetail)
	r.Get("/cart", h.Cart)
	r.Post("/cart/add", h.AddToCart)
	r.Get("/checkout", h.Checkout)
	r.Post("/checkout", h.PlaceOrder)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

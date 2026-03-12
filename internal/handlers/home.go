package handlers

import (
	"log"
	"net/http"

	"github.com/your-username/shop01/internal/models"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	products, err := h.db.ListFeaturedProducts(r.Context())
	if err != nil {
		log.Printf("error fetching featured products: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title    string
		Products []models.Product
	}{
		Title:    "Welcome to Shop01",
		Products: products,
	}
	h.render(w, "home.html", data)
}

package handlers

import (
	"net/http"

	"github.com/your-username/shop01/internal/models"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// TODO: fetch featured products from DB
	data := struct {
		Title    string
		Products []models.Product
	}{
		Title: "Welcome to Shop01",
		Products: []models.Product{
			{ID: 1, Name: "Sample Product", Description: "A great product", Price: 29.99},
		},
	}
	h.render(w, "home.html", data)
}

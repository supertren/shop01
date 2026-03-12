package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/your-username/shop01/internal/models"
)

func (h *Handler) Products(w http.ResponseWriter, r *http.Request) {
	// TODO: fetch all products from DB
	data := struct {
		Title    string
		Products []models.Product
	}{
		Title:    "Products",
		Products: []models.Product{},
	}
	h.render(w, "products.html", data)
}

func (h *Handler) ProductDetail(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "id")
	// TODO: fetch product by ID from DB
	h.render(w, "product_detail.html", nil)
}

func (h *Handler) Cart(w http.ResponseWriter, r *http.Request) {
	h.render(w, "cart.html", nil)
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	// TODO: add product to session-based cart
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	h.render(w, "checkout.html", nil)
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	// TODO: process payment and create order
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

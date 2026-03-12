package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-username/shop01/internal/models"
	"github.com/your-username/shop01/internal/store"
)

func (h *Handler) Products(w http.ResponseWriter, r *http.Request) {
	products, err := h.db.ListProducts(r.Context())
	if err != nil {
		log.Printf("error listing products: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title    string
		Products []models.Product
	}{
		Title:    "Products",
		Products: products,
	}
	h.render(w, "products.html", data)
}

func (h *Handler) ProductDetail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("invalid product ID in URL: %q", idStr)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.db.GetProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		log.Printf("error getting product %d: %v", id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Product *models.Product
	}{
		Title:   product.Name,
		Product: product,
	}
	h.render(w, "product_detail.html", data)
}

func (h *Handler) Cart(w http.ResponseWriter, r *http.Request) {
	h.render(w, "cart.html", nil)
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("error parsing add to cart form: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	idStr := r.PostFormValue("product_id")
	qtyStr := r.PostFormValue("quantity")

	productID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("invalid product_id in form: %q", idStr)
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(qtyStr)
	if err != nil || quantity < 1 {
		log.Printf("invalid quantity in form: %q", qtyStr)
		http.Error(w, "Invalid Quantity", http.StatusBadRequest)
		return
	}

	// Validate product exists
	product, err := h.db.GetProduct(r.Context(), productID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		log.Printf("error getting product %d: %v", productID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Validate quantity is within stock limits
	if quantity > product.Stock {
		http.Error(w, "Not enough items in stock", http.StatusBadRequest)
		return
	}

	// TODO: Add product to a session-based cart.
	log.Printf("ACTION: Add to cart ProductID=%d, Quantity=%d", productID, quantity)

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	h.render(w, "checkout.html", nil)
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	// TODO: process payment and create order
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

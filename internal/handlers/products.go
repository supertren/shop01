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
	cart := readCart(r)

	var items []models.CartItem
	var total float64
	for productID, qty := range cart {
		product, err := h.db.GetProduct(r.Context(), productID)
		if err != nil {
			continue
		}
		item := models.CartItem{Product: *product, Quantity: qty}
		items = append(items, item)
		total += item.Subtotal()
	}

	data := struct {
		Title string
		Items []models.CartItem
		Total float64
	}{
		Title: "Your Cart",
		Items: items,
		Total: total,
	}
	h.render(w, "cart.html", data)
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

	cart := readCart(r)
	cart[productID] += quantity
	writeCart(w, cart)

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	idStr := r.PostFormValue("product_id")
	productID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	cart := readCart(r)
	delete(cart, productID)
	writeCart(w, cart)

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	h.render(w, "checkout.html", nil)
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	cart := readCart(r)
	if len(cart) == 0 {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	}

	var total float64
	for productID, qty := range cart {
		product, err := h.db.GetProduct(r.Context(), productID)
		if err != nil {
			continue
		}
		total += product.Price * float64(qty)
	}

	payment := r.PostFormValue("payment")
	switch payment {
	case "paypal":
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		base := scheme + "://" + r.Host
		order, err := h.paypal.CreateOrder(r.Context(), total, base+"/paypal/success", base+"/paypal/cancel")
		if err != nil {
			log.Printf("error creating paypal order: %v", err)
			http.Error(w, "Failed to initiate PayPal payment", http.StatusInternalServerError)
			return
		}
		approvalURL := order.ApprovalURL()
		if approvalURL == "" {
			http.Error(w, "No PayPal approval URL returned", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, approvalURL, http.StatusSeeOther)
	default:
		http.Error(w, "Payment method not yet implemented", http.StatusNotImplemented)
	}
}

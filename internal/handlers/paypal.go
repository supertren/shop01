package handlers

import (
	"log"
	"net/http"

	"github.com/your-username/shop01/internal/models"
)

func (h *Handler) PayPalSuccess(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("token")
	if orderID == "" {
		http.Error(w, "Missing PayPal order token", http.StatusBadRequest)
		return
	}

	// Capture the payment
	if err := h.paypal.CaptureOrder(r.Context(), orderID); err != nil {
		log.Printf("error capturing paypal order %s: %v", orderID, err)
		http.Error(w, "Payment capture failed", http.StatusInternalServerError)
		return
	}

	// Build cart items from cookie
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

	// Save order to database
	dbOrderID, err := h.db.CreateOrder(r.Context(), total, items)
	if err != nil {
		log.Printf("error saving order: %v", err)
		http.Error(w, "Failed to save order", http.StatusInternalServerError)
		return
	}

	// Clear cart
	clearCart(w)

	log.Printf("PayPal order captured: paypal_id=%s db_order_id=%d total=%.2f", orderID, dbOrderID, total)

	data := struct {
		Title   string
		OrderID int
		Total   float64
	}{
		Title:   "Order Confirmed",
		OrderID: dbOrderID,
		Total:   total,
	}
	h.render(w, "order_success.html", data)
}

func (h *Handler) PayPalCancel(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/checkout", http.StatusSeeOther)
}

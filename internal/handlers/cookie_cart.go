package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

const cartCookieName = "cart"

// cartData maps product ID to quantity.
type cartData map[int]int

func readCart(r *http.Request) cartData {
	cookie, err := r.Cookie(cartCookieName)
	if err != nil {
		return cartData{}
	}
	raw, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return cartData{}
	}
	var cart cartData
	if err := json.Unmarshal(raw, &cart); err != nil {
		return cartData{}
	}
	return cart
}

func writeCart(w http.ResponseWriter, cart cartData) {
	raw, err := json.Marshal(cart)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     cartCookieName,
		Value:    base64.URLEncoding.EncodeToString(raw),
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7, // 7 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearCart(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   cartCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

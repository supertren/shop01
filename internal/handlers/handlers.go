package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/your-username/shop01/internal/payments/paypal"
	"github.com/your-username/shop01/internal/store"
)

// Handler holds application-wide dependencies.
type Handler struct {
	db     *store.DB
	paypal *paypal.Client
}

// New creates a new Handler with dependencies.
func New(db *store.DB, pp *paypal.Client) *Handler {
	return &Handler{
		db:     db,
		paypal: pp,
	}
}

// render executes a template.
func (h *Handler) render(w http.ResponseWriter, name string, data any) {
	files := []string{
		"web/templates/layout.html",
		fmt.Sprintf("web/templates/%s", name),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("error parsing template %s: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := ts.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("error executing template %s: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

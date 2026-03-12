package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/your-username/shop01/internal/store"
)

// Handler holds application-wide dependencies.
type Handler struct {
	db *store.DB
	// In a real app, templates should be parsed once and cached for performance.
}

// New creates a new Handler with dependencies.
func New(db *store.DB) *Handler {
	return &Handler{
		db: db,
	}
}

// render executes a template.
// NOTE: Parsing templates on every request is inefficient. This should be
// optimized in a production application by caching the parsed templates.
func (h *Handler) render(w http.ResponseWriter, name string, data any) {
	files := []string{
		"web/templates/layouts/base.html",
		fmt.Sprintf("web/templates/pages/%s", name),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("error parsing template %s: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := ts.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("error executing template %s: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}


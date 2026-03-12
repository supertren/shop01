package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/your-username/shop01/internal/store"
)

type Handler struct {
	db        *store.DB
	templates *template.Template
}

func New(db *store.DB) *Handler {
	tmpl := template.Must(template.ParseGlob(filepath.Join("web", "templates", "*.html")))
	return &Handler{db: db, templates: tmpl}
}

func (h *Handler) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
	}
}

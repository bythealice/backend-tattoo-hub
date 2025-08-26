package http

import (
	"backend-tattoo-hub/internal/http/handler"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func NewRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/submit", handler.Submit(db))
	r.Get("/feed", handler.Feed(db))
	return r
}

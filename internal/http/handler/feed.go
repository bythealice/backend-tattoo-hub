package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"backend-tattoo-hub/internal/model"
)

func Feed(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT id, type, payload, created_at
			FROM submissions
			ORDER BY created_at DESC
			LIMIT 50`)
		if err != nil {
			http.Error(w, "erro ao consultar", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var out []model.FeedItem
		for rows.Next() {
			var item model.FeedItem
			if err := rows.Scan(&item.ID, &item.Type, &item.Payload, &item.CreatedAt); err != nil {
				http.Error(w, "erro ao ler linha", http.StatusInternalServerError)
				return
			}
			out = append(out, item)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	}
}

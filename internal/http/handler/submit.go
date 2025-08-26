package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"backend-tattoo-hub/internal/model"
)

func Submit(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body model.SubmitReq
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "json inválido", http.StatusBadRequest)
			return
		}
		if body.Type != "url" && body.Type != "image" {
			http.Error(w, "type deve ser 'url' ou 'image'", http.StatusBadRequest)
			return
		}
		if body.Payload == "" {
			http.Error(w, "payload obrigatório", http.StatusBadRequest)
			return
		}
		if _, err := db.Exec(`INSERT INTO submissions (type, payload) VALUES ($1, $2)`, body.Type, body.Payload); err != nil {
			http.Error(w, "erro ao inserir", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

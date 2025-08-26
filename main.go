package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

type SubmitReq struct {
	Type    string `json:"type"`    // "url" | "image"
	Payload string `json:"payload"` // link ou url da imagem
}

type FeedItem struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

func mustDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL não definido (ex.: postgres://localhost:5432/tattoo?sslmode=disable)")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := mustDB()
	defer db.Close()

	r := chi.NewRouter()

	// POST /submit  -> insere na tabela submissions
	r.Post("/submit", func(w http.ResponseWriter, r *http.Request) {
		var body SubmitReq
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
	})

	// GET /feed -> últimos itens
	r.Get("/feed", func(w http.ResponseWriter, r *http.Request) {
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

		var out []FeedItem
		for rows.Next() {
			var item FeedItem
			if err := rows.Scan(&item.ID, &item.Type, &item.Payload, &item.CreatedAt); err != nil {
				http.Error(w, "erro ao ler linha", http.StatusInternalServerError)
				return
			}
			out = append(out, item)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	})

	log.Println("API on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

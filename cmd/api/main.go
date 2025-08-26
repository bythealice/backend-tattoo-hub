package main

import (
	"log"
	nethttp "net/http"
	"os"

	"backend-tattoo-hub/internal/db"
	apihttp "backend-tattoo-hub/internal/http"
)

func main() {
	conn := db.MustConnect()
	defer conn.Close()

	r := apihttp.NewRouter(conn)

	addr := ":8080"
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		addr = v
	}

	log.Println("API on", addr)
	log.Fatal(nethttp.ListenAndServe(addr, r))
}

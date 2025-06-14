package server

import (
	"database/sql"
	"log"
	"net/http"
)

var Addr = "http://localhost:3001"

func Serve(db *sql.DB) {
	handler := &Handler{DB: db}

	server := &http.Server{
		Addr:    ":3001",
		Handler: routes(handler),
	}

	log.Println("Server running on :3001")
	log.Fatal(server.ListenAndServe())
}

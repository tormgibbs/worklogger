package server
	
import (
	"log"
	"net/http"

	"github.com/tormgibbs/worklogger/data"
)

func Serve() {
	db := data.NewSQLiteDB("db.sqlite")
	defer db.Close()

	handler := &Handler{DB: db}

	server := &http.Server{
		Addr:    ":8080",
		Handler: routes(handler),
	}

	log.Println("Server running on :8080")
	log.Fatal(server.ListenAndServe())
}

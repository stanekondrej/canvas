package main

import (
	"log"
	"net/http"

	"github.com/stanekondrej/canvas/server/internal/app/server"
)

const LISTEN_ADDRESS string = ":9999"

func main() {
	h := server.NewHandler()

	http.HandleFunc("GET /", h.Handle)

	log.Println("Starting server")
	log.Fatalln(http.ListenAndServe(LISTEN_ADDRESS, nil))
}

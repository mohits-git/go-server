package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := "8080"

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Server listening on port", port)
	log.Fatal(server.ListenAndServe())
}

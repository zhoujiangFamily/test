package main

import (
	"git.in.codoon.com/Overseas/runbox/first-test/service"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)

	http.HandleFunc("/", service.Votes)
	http.HandleFunc("/gps", service.Gps)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

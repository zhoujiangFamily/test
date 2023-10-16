package runbox

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("starting server 123...")
	http.HandleFunc("/", Handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

// [END run_helloworld_service]
// [END cloudrun_helloworld_service]

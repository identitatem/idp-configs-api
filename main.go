package main

import (
	"fmt"
	"net/http"
)

// Handler function that responds with Hello World
func helloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello world")
    fmt.Fprintf(w, "See actions in action?")
}

func main() {
    // Register handler function on server route
    http.HandleFunc("/api/hello-world-service/v0/ping", helloWorld)

    fmt.Println("Listening on localhost:8080")
    http.ListenAndServe(":8080", nil)
}

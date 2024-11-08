// main.go
// Application entry point.

package main

import (
    "fmt"
    "net/http"
    "receipt-processor/api"

    "github.com/gorilla/mux"
)

func main() {
    router := mux.NewRouter()

    // Set up routes
    api.SetupRouter(router)

    fmt.Println("Server is running on port 8080...")
    err := http.ListenAndServe(":8080", router)
    if err != nil {
        panic(err)
    }
}
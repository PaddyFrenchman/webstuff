package main

import (
    "log"
    "net/http"
    "myapp/internal/handlers"
)

func main() {
    mux := http.NewServeMux()
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
    mux.HandleFunc("/", handlers.HomeHandler)
    mux.HandleFunc("/api/login", handlers.LoginHandler)

    log.Println("Server running at http://localhost:8080")
    err := http.ListenAndServe(":8080", mux)
    if err != nil {
        log.Fatal(err)
    }
}

package main

import (
    "encoding/json"
    "net/http"
)

func contactsHandler(w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get("Authorization")
    if token != "Bearer securetoken" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    contacts := []map[string]string{
        {"name": "Alice", "email": "alice@example.com"},
        {"name": "Bob", "email": "bob@example.com"},
    }
    json.NewEncoder(w).Encode(contacts)
}

func main() {
    http.HandleFunc("/contacts", contactsHandler)
    http.ListenAndServeTLS(":8082", "cert.pem", "key.pem", nil)
}

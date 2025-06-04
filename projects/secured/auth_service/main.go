
package main

import (
    "encoding/json"
    "net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")
    if username == "admin" && password == "pass" {
        json.NewEncoder(w).Encode(map[string]string{"token": "securetoken"})
    } else {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
    }
}

func main() {
    http.HandleFunc("/login", loginHandler)
    http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil)
}

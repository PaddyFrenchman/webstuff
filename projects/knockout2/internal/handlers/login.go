package handlers

import (
    "encoding/json"
    "net/http"
)

type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var creds Credentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    if creds.Email == "admin@example.com" && creds.Password == "password" {
        json.NewEncoder(w).Encode(map[string]string{"token": "fake-oauth2-token"})
        return
    }

    http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
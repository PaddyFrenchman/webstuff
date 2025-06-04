package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	users = map[string]string{"paddy@x.mail": "pass1"}

	sessions   = map[string]string{}
	sessionsMu sync.Mutex

	contacts   = map[string]ContactData{}
	contactsMu sync.Mutex
)

type ContactData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Msg   string `json:"msg"`
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	http.HandleFunc("/", landingHandler)
	http.HandleFunc("/login_submit", loginSubmitHandler)
	//http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/contact", withAuth(contactHandler))
	http.HandleFunc("/api/contact", withAuth(contactAPIHandler))
	log.Println("Server running on http://localhost:8480")
	log.Fatal(http.ListenAndServe(":8480", nil))
}

func withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || !isValidSession(cookie.Value) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}

func isValidSession(sessionID string) bool {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	_, ok := sessions[sessionID]
	return ok
}

func getUsernameBySession(sessionID string) string {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	return sessions[sessionID]
}

func landingHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/landing.html")
	t.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(w, nil)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/main.html")
	t.Execute(w, nil)

}

func loginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Password != users[req.Email] {
		return
	}

	sessionID := req.Email + "_sess"
	sessionsMu.Lock()
	sessions[sessionID] = req.Email
	sessionsMu.Unlock()

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/main", http.StatusSeeOther)
	//http.Redirect(w, r, "/main.html", http.StatusFound)
	//http.Redirect(w, r, "/contact", http.StatusFound)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token != "Bearer securetoken" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	t, _ := template.ParseFiles("templates/contact.html")
	t.Execute(w, nil)
}

func contactAPIHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	username := getUsernameBySession(cookie.Value)

	switch r.Method {
	case "GET":
		contactsMu.Lock()
		data := contacts[username]
		contactsMu.Unlock()
		json.NewEncoder(w).Encode(data)

	case "POST":
		var data ContactData
		json.NewDecoder(r.Body).Decode(&data)
		contactsMu.Lock()
		contacts[username] = data
		contactsMu.Unlock()
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"

	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Book struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Status     string `json:"status"` // "available" or "borrowed"
	BorrowerID *int   `json:"borrower_id"`
}

type Member struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type BorrowRecord struct {
	ID         int        `json:"id"`
	BookID     int        `json:"book_id"`
	MemberID   int        `json:"member_id"`
	BorrowDate time.Time  `json:"borrow_date"`
	ReturnDate *time.Time `json:"return_date"`
}

type App struct {
	DB     *sql.DB
	Router *mux.Router
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("b81c1382-22ae-47a5-819e-d0034d7f2600")

// fileExists cehcks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true // file or directory exists
	}
	if os.IsNotExist(err) {
		return false // does not exist
	}
	// some other error, e.g., permission denied
	return false
}

// runningInDocker returns true if we are running in a container
func runningInDocker() bool {
	return fileExists("/.dockerenv")
}

func main() {
	dbFile := "./library.db/webknock.sqlite"
	if runningInDocker() {
		dbFile = "/app/library.db"
	}
	if !fileExists(dbFile) {
		os.Exit(1)
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &App{DB: db, Router: mux.NewRouter()}
	app.initializeDB()
	app.setupRoutes()

	// Load self-signed certificate
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal("Error loading certificate:", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	server := &http.Server{
		Addr:      ":8480",
		Handler:   app.Router,
		TLSConfig: config,
	}

	log.Println("Server starting on https://localhost:8480")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func (a *App) initializeDB() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE, password TEXT)`,
		`CREATE TABLE IF NOT EXISTS books (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, author TEXT, status TEXT, borrower_id INTEGER)`,
		`CREATE TABLE IF NOT EXISTS members (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)`,
		`CREATE TABLE IF NOT EXISTS borrow_records (id INTEGER PRIMARY KEY AUTOINCREMENT, book_id INTEGER, member_id INTEGER, borrow_date TEXT, return_date TEXT)`,
	}
	for _, query := range queries {
		_, err := a.DB.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (a *App) setupRoutes() {
	a.Router.HandleFunc("/api/register", a.register).Methods("POST")
	a.Router.HandleFunc("/api/login", a.login).Methods("POST")
	a.Router.HandleFunc("/api/books", a.authMiddleware(a.getBooks)).Methods("GET")
	a.Router.HandleFunc("/api/books", a.authMiddleware(a.createBook)).Methods("POST")
	a.Router.HandleFunc("/api/books/{id}", a.authMiddleware(a.updateBook)).Methods("PUT")
	a.Router.HandleFunc("/api/books/{id}", a.authMiddleware(a.deleteBook)).Methods("DELETE")
	a.Router.HandleFunc("/api/members", a.authMiddleware(a.getMembers)).Methods("GET")
	a.Router.HandleFunc("/api/members", a.authMiddleware(a.createMember)).Methods("POST")
	a.Router.HandleFunc("/api/members/{id}", a.authMiddleware(a.updateMember)).Methods("PUT")
	a.Router.HandleFunc("/api/members/{id}", a.authMiddleware(a.deleteMember)).Methods("DELETE")
	a.Router.HandleFunc("/api/borrow", a.authMiddleware(a.borrowBook)).Methods("POST")
	a.Router.HandleFunc("/api/return", a.authMiddleware(a.returnBook)).Methods("POST")

	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	a.Router.PathPrefix("/js").Handler(http.FileServer(http.Dir("./static/js")))
	a.Router.PathPrefix("/lib").Handler(http.FileServer(http.Dir("./static/lib")))
	a.Router.PathPrefix("/css").Handler(http.FileServer(http.Dir("./static/css")))
}

func (a *App) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	_, err = a.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, hashedPassword)
	if err != nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	var storedPassword string
	err := a.DB.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (a *App) getBooks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	query := "SELECT id, title, author, status, borrower_id FROM books"
	if status != "" {
		query += " WHERE status = ?"
	}
	rows, err := a.DB.Query(query, status)
	if err != nil {
		http.Error(w, "Error fetching books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var books []Book
	for rows.Next() {
		var book Book
		var borrowerID sql.NullInt64
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Status, &borrowerID); err != nil {
			http.Error(w, "Error scanning books", http.StatusInternalServerError)
			return
		}
		if borrowerID.Valid {
			bid := int(borrowerID.Int64)
			book.BorrowerID = &bid
		}
		books = append(books, book)
	}
	json.NewEncoder(w).Encode(books)
}

func (a *App) createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	book.Status = "available"
	_, err := a.DB.Exec("INSERT INTO books (title, author, status) VALUES (?, ?, ?)", book.Title, book.Author, book.Status)
	if err != nil {
		http.Error(w, "Error creating book", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *App) updateBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	_, err := a.DB.Exec("UPDATE books SET title = ?, author = ?, status = ? WHERE id = ?", book.Title, book.Author, book.Status, id)
	if err != nil {
		http.Error(w, "Error updating book", http.StatusInternalServerError)
		return
	}
}

func (a *App) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := a.DB.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error deleting book", http.StatusInternalServerError)
		return
	}
}

func (a *App) getMembers(w http.ResponseWriter, r *http.Request) {
	rows, err := a.DB.Query("SELECT id, name FROM members")
	if err != nil {
		http.Error(w, "Error fetching members", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var members []Member
	for rows.Next() {
		var member Member
		if err := rows.Scan(&member.ID, &member.Name); err != nil {
			http.Error(w, "Error scanning members", http.StatusInternalServerError)
			return
		}
		members = append(members, member)
	}
	json.NewEncoder(w).Encode(members)
}

func (a *App) createMember(w http.ResponseWriter, r *http.Request) {
	var member Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	_, err := a.DB.Exec("INSERT INTO members (name) VALUES (?)", member.Name)
	if err != nil {
		http.Error(w, "Error creating member", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *App) updateMember(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var member Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	_, err := a.DB.Exec("UPDATE members SET name = ? WHERE id = ?", member.Name, id)
	if err != nil {
		http.Error(w, "Error updating member", http.StatusInternalServerError)
		return
	}
}

func (a *App) deleteMember(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := a.DB.Exec("DELETE FROM members WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error deleting member", http.StatusInternalServerError)
		return
	}
}

func (a *App) borrowBook(w http.ResponseWriter, r *http.Request) {
	var record BorrowRecord
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	_, err := a.DB.Exec("UPDATE books SET status = 'borrowed', borrower_id = ? WHERE id = ?", record.MemberID, record.BookID)
	if err != nil {
		http.Error(w, "Error borrowing book", http.StatusInternalServerError)
		return
	}
	_, err = a.DB.Exec("INSERT INTO borrow_records (book_id, member_id, borrow_date) VALUES (?, ?, ?)", record.BookID, record.MemberID, time.Now())
	if err != nil {
		http.Error(w, "Error recording borrow", http.StatusInternalServerError)
		return
	}
}

func (a *App) returnBook(w http.ResponseWriter, r *http.Request) {
	var record BorrowRecord
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	_, err := a.DB.Exec("UPDATE books SET status = 'available', borrower_id = NULL WHERE id = ?", record.BookID)
	if err != nil {
		http.Error(w, "Error returning book", http.StatusInternalServerError)
		return
	}
	_, err = a.DB.Exec("UPDATE borrow_records SET return_date = ? WHERE book_id = ? AND member_id = ? AND return_date IS NULL", time.Now(), record.BookID, record.MemberID)
	if err != nil {
		http.Error(w, "Error recording return", http.StatusInternalServerError)
		return
	}
}

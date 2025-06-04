/*
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'Paddy0724827140#';
FLUSH PRIVILEGES;
EXIT;
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"onepageapp/data"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	if data.ValidateUser(db, username, password) {
		json.NewEncoder(w).Encode(map[string]string{"token": "securetoken"})
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	err := data.RegisterUser(db, username, password)
	if err != nil {
		http.Error(w, "Registration error", http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

func contactsHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !strings.Contains(token, "securetoken") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	contacts := data.GetContacts()
	json.NewEncoder(w).Encode(contacts)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// token := r.Header.Get("Authorization")
	// if !strings.Contains(token, "securetoken") {
	// 	http.Error(w, "Unauthorized to main page", http.StatusUnauthorized)
	// 	return
	// }
	//http.ServeFile(w, r, "static/index.html")

	vars := mux.Vars(r)
	fmt.Printf("%v\n", vars)

	//title := vars["title"]
	//page := vars["page"]
	if r.RequestURI == "/" {
		http.ServeFile(w, r, "static/index.html")
	}
	if r.RequestURI == "/about" {
		http.ServeFile(w, r, "static/views/about.html")
	}
	if r.RequestURI == "/home" {
		http.ServeFile(w, r, "static/views/home.html")
	}
}
func main() {
	dsn := "root:Paddy0724827140#@tcp(localhost:3306)/truely?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	db.AutoMigrate(&data.User{})
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexHandler)

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "static/index.html")
	// })
	//http.Handle("/", fs)
	// http.HandleFunc("/login", loginHandler)
	// http.HandleFunc("/register", registerHandler)
	// http.HandleFunc("/contacts", contactsHandler)

	fmt.Println("Server running at https://localhost:8443")
	log.Fatal(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil))
}

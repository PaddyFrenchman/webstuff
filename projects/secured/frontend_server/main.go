package main

import (
	"html/template"
	"log"
	"net/http"
)

func contactHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/contact.html")
	t.Execute(w, nil)
}

func main() {
	fs1 := http.FileServer(http.Dir("./assets"))
	fs2 := http.FileServer(http.Dir("./templates"))
	http.Handle("/", fs1)
	http.Handle("/assets", fs2)
	log.Println("Frontend server running on https://localhost:8080")
	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil))
}

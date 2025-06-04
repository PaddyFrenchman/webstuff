package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handleBookPage(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	title := vars["title"]
	page := vars["page"]
	fmt.Fprintf(response, "You've requested the book: %s on page %s\n", title, page)
}

// http://localhost:8800/books/1984/page/10
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/books/{title}/page/{page}", handleBookPage)
	http.ListenAndServe(":8800", r)
}

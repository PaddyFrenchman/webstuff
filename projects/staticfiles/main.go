package main

import (
	"fmt"
	"net/http"
)

// http://localhost:8800/static/form1.html
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})
	fs1 := http.FileServer(http.Dir("static/"))
	fs2 := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs1))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs2))
	http.ListenAndServe(":8800", nil)
}

package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles(filepath.Join("web", "templates", "index.html")))
    tmpl.Execute(w, nil)
}
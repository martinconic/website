package main

import (
	"html/template"
	"net/http"
)

var TPL *template.Template

func init() {
	TPL = template.Must(template.ParseGlob("*.html"))
}


func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	reqID := r.URL.Query().Get("id")


	TPL.ExecuteTemplate(w, "index.html", reqID)
}

func main() {
	http.HandleFunc("/", home)
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.ListenAndServe(":80", http.DefaultServeMux)
}
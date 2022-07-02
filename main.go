package main

import (
	"golang.org/x/crypto/acme/autocert"
	"html/template"
	"log"
	"net/http"
)

var TPL *template.Template

func init() {
	TPL = template.Must(template.ParseGlob("*.html"))
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		http.StatusTemporaryRedirect)
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

	m := &autocert.Manager{
		Cache:      autocert.DirCache("secret-dir"),
		Prompt:     autocert.AcceptTOS,
		Email:      "martinconic@gmail.com",
		HostPolicy: autocert.HostWhitelist("martinconic.ro", "www.martinconic.ro"),
	}
	s := &http.Server{
		Addr:      ":https",
		TLSConfig: m.TLSConfig(),
	}

	go http.ListenAndServe(":80", http.HandlerFunc(redirect))

	log.Fatalln(s.ListenAndServeTLS("", ""))
}

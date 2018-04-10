package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"html/template"

	"github.com/apex/log"
	jsonlog "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/gorilla/pat"
)

func init() {
	if os.Getenv("UP_STAGE") == "" {
		log.SetHandler(text.Default)
	} else {
		log.SetHandler(jsonlog.Default)
	}
}

func main() {
	addr := ":" + os.Getenv("PORT")
	app := pat.New()

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	app.PathPrefix("/static/").Handler(s)

	app.Get("/", get)

	if err := http.ListenAndServe(addr, app); err != nil {
		log.WithError(err).Fatal("error listening")
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("").ParseGlob("*.html"))
	content, err := ioutil.ReadFile("bins")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	bins := strings.Split(string(content), "\n")
	t.ExecuteTemplate(w, "polyfill.html", bins)
}
